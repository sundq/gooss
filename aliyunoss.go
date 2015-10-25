package aliyunoss

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"
	"time"
)

type oss_agent struct {
	AccessKey            string
	AccessKeySecret      string
	Verb                 string
	CanonicalizedHeaders map[string]string
	CanonicalizedUri     string
	CanonicalizedQuery   map[string]string
	Content              []byte
	ContentType          string
	Date                 string
	Url                  string
	Debug                bool
	logger               *log.Logger
}

type AliOSSClient struct {
	AccessKey       string
	AccessKeySecret string
	EndPoint        string
	Debug           bool
	logger          *log.Logger
}

type Bucket struct {
	Location     string `xml:"Location"`
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}

type BucketList struct {
	// XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type AliOssError struct {
	Code           string `xml:"Code"`
	Message        string `xml:"Message"`
	RequestId      string `xml:"RequestId"`
	HostId         string `xml:"HostId"`
	OSSAccessKeyId string `xml:"OSSAccessKeyId"`
}

func (e AliOssError) Error() string {
	return fmt.Sprintf("%v: %v RequestId:%s", e.Code, e.Message, e.RequestId)
}

// type Request *gorequest.Request
// type Response *gorequest.Response

func (s *oss_agent) calc_signature() string {
	//sort the canonicalized headers
	sorted_canonicalized_headers_str := ""
	var header_keys []string
	for k := range s.CanonicalizedHeaders {
		header_keys = append(header_keys, k)
	}
	sort.Strings(header_keys)
	for _, k := range header_keys {
		sorted_canonicalized_headers_str += (strings.Trim(k, " ") + ":" + strings.Trim(s.CanonicalizedHeaders[k], " ") + "\n")
	}

	date := s.Date
	content_md5 := ""
	if len(s.Content) > 0 {
		sum := md5.Sum(s.Content)
		content_md5 = hex.EncodeToString(sum[:])
	}
	canonicalized_resource_str := s.CanonicalizedUri

	signature_ele := []string{s.Verb, content_md5, s.ContentType, date, sorted_canonicalized_headers_str + canonicalized_resource_str}
	signature_str := strings.Join(signature_ele, "\n")
	mac := hmac.New(sha1.New, []byte(s.AccessKeySecret))
	mac.Write([]byte(signature_str))
	result := mac.Sum(nil)
	b4str := base64.StdEncoding.EncodeToString(result)
	return b4str
}

func (s *oss_agent) send_request() (*http.Response, error) {
	client := &http.Client{}
	sig := s.calc_signature()
	req, err := http.NewRequest(s.Verb, s.Url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Date", s.Date)
	req.Header.Add("Authorization", "OSS "+s.AccessKey+":"+sig)
	for k, v := range s.CanonicalizedHeaders {
		req.Header.Add(k, v)
	}

	q := req.URL.Query()
	for k, v := range s.CanonicalizedQuery {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	if s.Debug {
		dump, err := httputil.DumpRequest(req, true)
		if nil != err {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Request: %s", string(dump))
		}
	}
	return client.Do(req)

}

const (
	default_endpoint = "oss.aliyuncs.com"
)

func New(access_key string, access_key_secret string, endpoint interface{}, debug interface{}) *AliOSSClient {
	cur_endpoint := default_endpoint
	debug_mode := false
	if endpoint != nil {
		if v, ok := endpoint.(string); ok {
			cur_endpoint = v
		}
	}

	if debug != nil {
		if v, ok := debug.(bool); ok {
			debug_mode = v
		}
	}
	s := &AliOSSClient{
		AccessKey:       access_key,
		AccessKeySecret: access_key_secret,
		EndPoint:        cur_endpoint,
		Debug:           debug_mode,
		logger:          log.New(os.Stderr, "[aliyunoss]", log.LstdFlags),
	}

	return s
}

func (c *AliOSSClient) ListBucket(prefix string, maker string, max_size int) (*simplejson.Json, error) {
	t := time.Now().UTC()
	date := t.Format("Mon, 02 Jan 2006 15:04:05 GMT")
	uri := "/"
	query := make(map[string]string)

	if prefix != "" {
		query["prefix"] = prefix
	}

	if maker != "" {
		query["maker"] = maker
	}

	if max_size > 0 {
		query["max-keys"] = fmt.Sprintf("%d", max_size)
	}

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  "http://" + c.EndPoint,
		Date:                 date,
		CanonicalizedHeaders: make(map[string]string),
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              []byte(""),
		ContentType:          "",
		Debug:                c.Debug,
		logger:               c.logger,
	}

	v := &BucketList{}
	e := &AliOssError{}
	resp, err := s.send_request()
	defer resp.Body.Close()
	if err != nil {
		return &simplejson.Json{}, err
	}
	if s.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if nil != err {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Response: %s", string(dump))
		}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	xml_result := string(body)
	if resp.StatusCode == 200 {
		xml.Unmarshal([]byte(xml_result), v)
		result, _ := json.Marshal(v)
		js, _ := simplejson.NewJson(result)
		return js, nil
	} else {
		xml.Unmarshal([]byte(xml_result), e)
		return &simplejson.Json{}, e
	}
}
