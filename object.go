package aliyunoss

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type Object struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Type         string `xml:"Type"`
	Size         string `xml:"Size"`
}

type ObjectList struct {
	// XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	BucketName string   `xml:"Name"`
	Prefix     string   `xml:"Prefix"`
	Marker     string   `xml:"Marker"`
	Delimiter  string   `xml:"Delimiter"`
	MaxKeys    string   `xml:"MaxKeys"`
	Objects    []Object `xml:"Contents"`
}

func (c *AliOSSClient) ListObject(bucket string, delimiter string, marker string, max_size int, prefix string) (*ObjectList, error) {
	uri := fmt.Sprintf("/%s/", bucket)
	query := make(map[string]string)

	if delimiter != "" {
		query["delimiter"] = delimiter
	}

	if prefix != "" {
		query["prefix"] = prefix
	}

	if marker != "" {
		query["marker"] = marker
	}

	if max_size > 0 {
		query["max-keys"] = fmt.Sprintf("%d", max_size)
	}

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  fmt.Sprintf("http://%s.%s", bucket, c.EndPoint),
		CanonicalizedHeaders: make(map[string]string),
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		ContentType:          "",
		Debug:                c.Debug,
		logger:               c.logger,
	}

	v := &ObjectList{}
	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		xml.Unmarshal(xml_result, v)
		return v, nil
	} else {
		xml.Unmarshal(xml_result, e)
		return &ObjectList{}, e
	}
}

func (c *AliOSSClient) CreateObjectForBuff(bucket string, key string, data []byte, permission string) error {
	uri := fmt.Sprintf("/%s/%s", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)

	if permission == "" {
		header["x-oss-acl"] = "private"
	} else {
		header["x-oss-acl"] = permission
	}

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "PUT",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              bytes.NewReader(data),
		ContentType:          "application/octet-stream",
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}
