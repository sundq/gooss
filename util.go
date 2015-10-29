package aliyunoss

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"
	"time"
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
	Content              *bytes.Reader //io.Reader
	ContentType          string
	ContentMd5           string
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

func (s *oss_agent) calc_signature(date string) string {
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

	content_md5 := ""
	md5h := md5.New()
	io.Copy(md5h, s.Content)
	sum := md5h.Sum(nil)
	content_md5 = hex.EncodeToString(sum[:])
	s.ContentMd5 = content_md5
	s.Content.Seek(0, 0)

	canonicalized_resource_str := s.CanonicalizedUri

	signature_ele := []string{s.Verb, "", s.ContentType, date, sorted_canonicalized_headers_str + canonicalized_resource_str}
	signature_str := strings.Join(signature_ele, "\n")
	if s.Debug {
		s.logger.Println("signature string:", signature_str)
	}
	mac := hmac.New(sha1.New, []byte(s.AccessKeySecret))
	mac.Write([]byte(signature_str))
	result := mac.Sum(nil)
	b4str := base64.StdEncoding.EncodeToString(result)
	return b4str
}

func (s *oss_agent) send_request(is_stream bool) (*http.Response, []byte, error) {
	t := time.Now().UTC()
	date := t.Format("Mon, 02 Jan 2006 15:04:05 GMT")
	client := &http.Client{}
	sig := s.calc_signature(date)
	req, err := http.NewRequest(s.Verb, s.Url, s.Content)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Date", date)
	req.Header.Add("Authorization", fmt.Sprintf("OSS %s:%s", s.AccessKey, sig))
	// req.Header.Add("Content-Md5", s.ContentMd5)
	if s.ContentType != "" {
		req.Header.Add("Content-Type", s.ContentType)
	}
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
	defer resp.Body.Close()
	if s.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if nil != err {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Response: %s", string(dump))
		}
	}

	if err != nil {
		return nil, nil, err
	}
	if !is_stream {
		body, _ := ioutil.ReadAll(resp.Body)
		return resp, body, nil
	} else {
		return resp, nil, nil
	}
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
