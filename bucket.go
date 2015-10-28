package aliyunoss

import (
	"encoding/xml"
	"fmt"
)

type Bucket struct {
	Location     string `xml:"Location"`
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}

type BucketList struct {
	// XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type BucketConfiguration struct {
	XMLName            xml.Name `xml:"CreateBucketConfiguration"`
	LocationConstraint string   `xml:LocationConstraint`
}

func (c *AliOSSClient) ListBucket(prefix string, marker string, max_size int) (*BucketList, error) {
	uri := "/"
	query := make(map[string]string)

	if prefix != "" {
		query["prefix"] = prefix
	}

	if marker != "" {
		query["maker"] = marker
	}

	if max_size > 0 {
		query["max-keys"] = fmt.Sprintf("%d", max_size)
	}

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  "http://" + c.EndPoint,
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
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		xml.Unmarshal(xml_result, v)
		return v, nil
	} else {
		xml.Unmarshal(xml_result, e)
		return nil, e
	}
}

func (c *AliOSSClient) CreateBucket(name string, location string) error {
	bucket_config := &BucketConfiguration{LocationConstraint: location}
	xml_content, _ := xml.MarshalIndent(bucket_config, "", "  ")
	uri := fmt.Sprintf("/%s/", name)
	query := make(map[string]string)
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "PUT",
		Url:                  fmt.Sprintf("http://%s.%s", name, c.EndPoint),
		CanonicalizedHeaders: make(map[string]string),
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              []byte(xml.Header + string(xml_content)),
		ContentType:          "application/xml",
		Debug:                c.Debug,
		logger:               c.logger,
	}

	v := &BucketList{}
	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		xml.Unmarshal(xml_result, v)
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

func (c *AliOSSClient) GetLocationOfBucket(bucket string) (string, error) {
	uri := fmt.Sprintf("/%s/?location", bucket)
	query := make(map[string]string)
	query["location"] = ""
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  fmt.Sprintf("http://%s.%s", bucket, c.EndPoint),
		CanonicalizedHeaders: make(map[string]string),
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              []byte(""),
		ContentType:          "",
		Debug:                c.Debug,
		logger:               c.logger,
	}

	v := ""
	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {
		xml.Unmarshal(xml_result, &v)
		return v, nil
	} else {
		xml.Unmarshal(xml_result, e)
		return "", e
	}
}
