package aliyunoss

import (
	"encoding/xml"
	"fmt"
	"time"
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
	t := time.Now().UTC()
	date := t.Format("Mon, 02 Jan 2006 15:04:05 GMT")
	uri := fmt.Sprintf("/%s/", bucket)
	query := make(map[string]string)

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  fmt.Sprintf("http://%s.%s", bucket, c.EndPoint),
		Date:                 date,
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
	resp, xml_result, err := s.send_request(true)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		xml.Unmarshal([]byte(xml_result), v)
		return v, nil
	} else {
		xml.Unmarshal([]byte(xml_result), e)
		return &ObjectList{}, e
	}
}
