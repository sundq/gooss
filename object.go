package aliyunoss

import (
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
	Objects []Object `xml:"Contents"`
}

func (c *AliOSSClient) ListObject(bucket string) (*ObjectList, error) {
	uri := fmt.Sprintf("/%s/", bucket)
	query := make(map[string]string)

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
