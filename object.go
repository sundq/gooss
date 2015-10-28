package aliyunoss

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"
)

type Object struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Type         string `xml:"Type"`
	Size         string `xml:"Size"`
}

// <Contents>
//     <Key>fun/movie/001.avi</Key>
//     <LastModified>2012-02-24T08:43:07.000Z</LastModified>
//     <ETag>&quot;5B3C1A2E053D763E1B002CC607C5A0FE&quot;</ETag>
//     <Type>Normal</Type>
//     <Size>344606</Size>
//     <StorageClass>Standard</StorageClass>
//     <Owner>
//         <ID>00220120222</ID>
//         <DisplayName>user-example</DisplayName>
//     </Owner>
// </Contents>

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
	resp, err := s.send_request()
	defer resp.Body.Close()
	if err != nil {
		return &ObjectList{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	xml_result := string(body)
	if resp.StatusCode == 200 {
		xml.Unmarshal([]byte(xml_result), v)
		return v, nil
	} else {
		xml.Unmarshal([]byte(xml_result), e)
		return &ObjectList{}, e
	}
}
