package aliyunoss

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"
)

const (
	default_endpoint = "oss.aliyuncs.com"
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
	resp, err := client.Do(req)
	if s.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if nil != err {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Response: %s", string(dump))
		}
	}
	return resp, err

}

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
