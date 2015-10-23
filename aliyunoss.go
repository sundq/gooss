package aliyunoss

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"github.com/bitly/go-simplejson"
	"github.com/parnurzeal/gorequest"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type oss_agent struct {
	AccessKey             string
	AccessKeySecret       string
	Verb                  string
	CanonicalizedHeaders  map[string]string
	CanonicalizedResource string
	Content               []byte
	ContentType           string
	Date                  string
	Url                   string
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

type Request *gorequest.Request
type Response *gorequest.Response

func (s *oss_agent) calc_signature() string {
	//sort the canonicalized headers
	var keys []string
	for k := range s.CanonicalizedHeaders {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sorted_canonicalized_headers_str := ""
	for _, k := range keys {
		sorted_canonicalized_headers_str += (strings.Trim(k, " ") + ":" + strings.Trim(s.CanonicalizedHeaders[k], "") + "\n")
	}

	date := s.Date
	content_md5 := ""
	if len(s.Content) > 0 {
		sum := md5.Sum(s.Content)
		content_md5 = hex.EncodeToString(sum[:])
	}
	signature_str := s.Verb + "\n" + content_md5 + "\n" + s.ContentType + "\n" + date + "\n" + sorted_canonicalized_headers_str + s.CanonicalizedResource
	mac := hmac.New(sha1.New, []byte(s.AccessKeySecret))
	mac.Write([]byte(signature_str))
	result := mac.Sum(nil)
	b4str := base64.StdEncoding.EncodeToString(result)
	return b4str
}

func (s *oss_agent) send_request() (gorequest.Response, string, []error) {
	request := gorequest.New()
	request.Get(s.Url)
	request.Set("Date", s.Date)
	for k, v := range s.CanonicalizedHeaders {
		request.Set(k, v)
	}

	request.Set("Authorization", "OSS "+s.AccessKey+":"+s.calc_signature())
	return request.End()
}

type AliOSSClient struct {
	AccessKey       string
	AccessKeySecret string
	EndPoint        string
	Debug           bool
	logger          *log.Logger
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
		if v, ok := endpoint.(bool); ok {
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

func (c *AliOSSClient) ListBucket() (gorequest.Response, *simplejson.Json, []error) {
	t := time.Now().UTC()
	date := t.Format("Mon, 02 Jan 2006 15:04:05 GMT")
	s := &oss_agent{
		AccessKey:             c.AccessKey,
		AccessKeySecret:       c.AccessKeySecret,
		Verb:                  "GET",
		Url:                   "http://" + c.EndPoint,
		Date:                  date,
		CanonicalizedHeaders:  make(map[string]string),
		CanonicalizedResource: "/",
		Content:               []byte(""),
		ContentType:           "",
	}

	v := &BucketList{}
	resp, xml_result, errs := s.send_request()
	if errs != nil {
		return resp, &simplejson.Json{}, errs
	}
	xml.Unmarshal([]byte(xml_result), v)
	result, _ := json.Marshal(v)
	js, _ := simplejson.NewJson(result)
	return resp, js, errs
}
