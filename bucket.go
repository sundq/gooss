package aliyunoss

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"
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

func (c *AliOSSClient) ListBucket(prefix string, maker string, max_size int) (*BucketList, error) {
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
		return &BucketList{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	xml_result := string(body)
	if resp.StatusCode == 200 {
		xml.Unmarshal([]byte(xml_result), v)
		return v, nil
	} else {
		xml.Unmarshal([]byte(xml_result), e)
		return &BucketList{}, e
	}
}
