package aliyunoss

import (
	"bytes"
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

type BucketLogging struct {
	XMLName            xml.Name `xml:"BucketLoggingStatus"`
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
		Content:              &bytes.Reader{},
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

func (c *AliOSSClient) CreateBucket(name string, location string, permission string) error {
	bucket_config := &BucketConfiguration{LocationConstraint: location}
	xml_content, _ := xml.MarshalIndent(bucket_config, "", "  ")
	uri := fmt.Sprintf("/%s/", name)
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
		Url:                  fmt.Sprintf("http://%s.%s", name, c.EndPoint),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              bytes.NewReader([]byte(xml.Header + string(xml_content))),
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
	if resp.StatusCode/100 == 2 {
		xml.Unmarshal(xml_result, v)
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

func (c *AliOSSClient) ModifyBucketAcl(name string, permission string) error {
	uri := fmt.Sprintf("/%s/?acl", name)
	query := make(map[string]string)
	header := make(map[string]string)
	query["acl"] = ""
	if permission == "" {
		header["x-oss-acl"] = "private"
	} else {
		header["x-oss-acl"] = permission
	}

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "PUT",
		Url:                  fmt.Sprintf("http://%s.%s", name, c.EndPoint),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
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
	if resp.StatusCode/100 == 2 {
		xml.Unmarshal(xml_result, v)
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

func (c *AliOSSClient) OpenBucketLogging(name string, target_bucket string, obj_prefix string) error {
	xml_template := `<?xml version="1.0" encoding="UTF-8"?><BucketLoggingStatus><LoggingEnabled><TargetBucket>%s</TargetBucket>%s</LoggingEnabled></BucketLoggingStatus>`
	prefix_str := ""
	if obj_prefix != "" {
		prefix_str = fmt.Sprintf("<TargetPrefix>%s</TargetPrefix>", obj_prefix)
	}
	xml_content := fmt.Sprintf(xml_template, target_bucket, prefix_str)
	uri := fmt.Sprintf("/%s/?logging", name)
	query := make(map[string]string)
	header := make(map[string]string)
	query["logging"] = ""

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "PUT",
		Url:                  fmt.Sprintf("http://%s.%s", name, c.EndPoint),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              bytes.NewReader([]byte(xml_content)),
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
	if resp.StatusCode/100 == 2 {
		xml.Unmarshal(xml_result, v)
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

func (c *AliOSSClient) CloseBucketLogging(name string) error {
	xml_content := `<?xml version="1.0" encoding="UTF-8"?><BucketLoggingStatus></BucketLoggingStatus>`
	uri := fmt.Sprintf("/%s/?logging", name)
	query := make(map[string]string)
	header := make(map[string]string)
	query["logging"] = ""

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "PUT",
		Url:                  fmt.Sprintf("http://%s.%s", name, c.EndPoint),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              bytes.NewReader([]byte(xml_content)),
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
	if resp.StatusCode/100 == 2 {
		xml.Unmarshal(xml_result, v)
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

//Delete bucket
func (c *AliOSSClient) DeleteBucket(name string) error {
	uri := fmt.Sprintf("/%s/", name)
	query := make(map[string]string)
	header := make(map[string]string)

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "DELETE",
		Url:                  fmt.Sprintf("http://%s.%s", name, c.EndPoint),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	v := &BucketList{}
	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 2 {
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
		Content:              &bytes.Reader{},
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

func (c *AliOSSClient) CreateBucketWebsite(name string, index string, error_file string) error {
	xml_template :=
		`<?xml version="1.0" encoding="UTF-8"?><WebsiteConfiguration><IndexDocument><Suffix>%s</Suffix></IndexDocument>%s</WebsiteConfiguration>`
	xml_error_file :=
		`<ErrorDocument><Key>%s</Key></ErrorDocument>`
	error_file_content := ""
	if error_file_content != "" {
		error_file_content = fmt.Sprintf(xml_error_file, error_file)
	}
	xml_content := fmt.Sprintf(xml_template, index, error_file_content)
	uri := fmt.Sprintf("/%s/?website", name)
	query := make(map[string]string)
	header := make(map[string]string)
	query["website"] = ""

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "PUT",
		Url:                  fmt.Sprintf("http://%s.%s", name, c.EndPoint),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              bytes.NewReader([]byte(xml_content)),
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
	if resp.StatusCode/100 == 2 {
		xml.Unmarshal(xml_result, v)
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

func (c *AliOSSClient) AddBucketRefer(name string, allow_empty_referer bool, refer_list []string) error {
	xml_template :=
		`<?xml version="1.0" encoding="UTF-8"?><RefererConfiguration><AllowEmptyReferer>%t</AllowEmptyReferer><RefererList>%s</RefererList></RefererConfiguration>`
	xml_refer_list :=
		`<Referer>%s</Referer>`
	xml_refer_content := ""
	for _, refer := range refer_list {
		xml_refer_content += fmt.Sprintf(xml_refer_list, refer)
	}

	xml_content := fmt.Sprintf(xml_template, allow_empty_referer, xml_refer_content)
	uri := fmt.Sprintf("/%s/?referer", name)
	query := make(map[string]string)
	header := make(map[string]string)
	query["referer"] = ""

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "PUT",
		Url:                  fmt.Sprintf("http://%s.%s", name, c.EndPoint),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              bytes.NewReader([]byte(xml_content)),
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
	if resp.StatusCode/100 == 2 {
		xml.Unmarshal(xml_result, v)
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}
