package aliyunoss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Object struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Type         string `xml:"Type"`
	Size         string `xml:"Size"`
}

type ObjectList struct {
	BucketName string   `xml:"Name"`
	Prefix     string   `xml:"Prefix"`
	Marker     string   `xml:"Marker"`
	Delimiter  string   `xml:"Delimiter"`
	MaxKeys    string   `xml:"MaxKeys"`
	Objects    []Object `xml:"Contents"`
}

type DeleteObject struct {
	Object string `xml:"Key"`
}

type DeleteObjectList struct {
	XMLName xml.Name       `xml:"Delete"`
	Quiet   string         `xml:"Quiet"`
	Objects []DeleteObject `xml:"Object"`
}

type MultiUploadInit struct {
	XMLName  xml.Name `xml:"InitiateMultipartUploadResult"`
	Bucket   string   `xml:"Bucket"`
	Key      string   `xml:"Key"`
	UploadId string   `xml:"UploadId"`
}

type PartUpload struct {
	PartNumber int    `xml:"PartNumber"`
	ETag       string `xml:"ETag"`
}
type CompleteUpload struct {
	XMLName xml.Name     `xml:"CompleteMultipartUpload"`
	Part    []PartUpload `xml:"Part"`
}

type MultiUpload struct {
	Key       string `xml:"Key"`
	UploadId  string `xml:"UploadId"`
	Initiated string `xml:"Initiated"`
}

type MultiUploadList struct {
	XMLName            xml.Name      `xml:"ListMultipartUploadsResult"`
	Bucket             string        `xml:"Bucket"`
	KeyMarker          string        `xml:"KeyMarker"`
	UploadIdMarker     string        `xml:"UploadIdMarker"`
	NextKeyMarker      string        `xml:"NextKeyMarker"`
	NextUploadIdMarker string        `xml:"NextUploadIdMarker"`
	Delimiter          string        `xml:"Delimiter"`
	Prefix             string        `xml:"Prefix"`
	MaxUploads         string        `xml:"MaxUploads"`
	IsTruncated        string        `xml:"IsTruncated"`
	Uploads            []MultiUpload `xml:"Upload"`
}

//ListObject get the list of key for the specified bucket
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

//CreateObjectForBuff create a oss key for buffer, the permission can be private, public-read or public_write
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

//CreateObjectForBuff create a oss key for local file, the permission can be private, public-read or public_write
func (c *AliOSSClient) CreateObjectForFile(bucket string, key string, filepath string, permission string) error {
	uri := fmt.Sprintf("/%s/%s", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

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
		Content:              file,
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

//AppendObjectForBuff like CreateObjectForBuff but it will append the exist key
func (c *AliOSSClient) AppendObjectForBuff(bucket string, key string, position int, data []byte) (int, string, error) {
	uri := fmt.Sprintf("/%s/%s?append&position=%d", bucket, key, position)
	query := make(map[string]string)
	header := make(map[string]string)

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "POST",
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
	resp, _, err := s.send_request(true)
	if err != nil {
		return 0, "", err
	} else {
		defer resp.Body.Close()
	}
	if resp.StatusCode/100 == 2 {
		position, _ := strconv.Atoi(resp.Header.Get("x-oss-next-append-position"))
		return position, resp.Header.Get("x-oss-hash-crc64ecma"), nil
	} else {
		xml_result, _ := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(xml_result, e)
		return 0, "", e
	}
}

//DeleteObject delete a key
func (c *AliOSSClient) DeleteObject(bucket string, key string) error {
	uri := fmt.Sprintf("/%s/%s", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "DELETE",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 2 {
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

//DeleteMultiObject delete a multi-object
func (c *AliOSSClient) DeleteMultiObject(bucket string, keys []string) error {
	uri := fmt.Sprintf("/%s/?delete", bucket)
	query := make(map[string]string)
	header := make(map[string]string)

	objs := []DeleteObject{}
	for _, value := range keys {
		objs = append(objs, DeleteObject{value})
	}

	tmp := DeleteObjectList{Quiet: "true", Objects: objs}
	xml_content, _ := xml.MarshalIndent(tmp, "", "  ")

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "POST",
		Url:                  fmt.Sprintf("http://%s.%s/?delete", bucket, c.EndPoint),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              bytes.NewReader([]byte(xml.Header + string(xml_content))),
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 2 {
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

//GetObjectAsBuffer a object as buffer
func (c *AliOSSClient) GetObjectAsBuffer(bucket string, key string) ([]byte, error) {
	uri := fmt.Sprintf("/%s/%s", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return xml_result, err
	}
	if resp.StatusCode == 200 {
		return xml_result, nil
	} else {
		xml.Unmarshal(xml_result, e)
		return xml_result, e
	}
}

//GetObjectAsFile a object as local file
func (c *AliOSSClient) GetObjectAsFile(bucket string, key string, filepath string) error {
	uri := fmt.Sprintf("/%s/%s", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, _, err := s.send_request(true)
	if err != nil {
		return err
	} else {
		defer resp.Body.Close()
	}
	if resp.StatusCode/100 == 2 {
		_, err := io.Copy(file, resp.Body)
		return err
	} else {
		xml_result, _ := ioutil.ReadAll(resp.Body)
		xml.Unmarshal(xml_result, e)
		return e
	}
}

//GetObjectInfo information of object
func (c *AliOSSClient) GetObjectInfo(bucket string, key string) (http.Header, error) {
	uri := fmt.Sprintf("/%s/%s", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)

	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "HEAD",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	resp, _, err := s.send_request(true)
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
	}
	if resp.StatusCode/100 == 2 {
		return resp.Header, nil
	} else {
		e := &AliOssError{Code: "NotFound", Message: "the object does not exist."}
		return nil, e
	}
}

//GetObjectMetaData information of object
func (c *AliOSSClient) GetObjectMetaData(bucket string, key string) (http.Header, error) {
	uri := fmt.Sprintf("/%s/%s", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)
	query["objectMeta"] = ""
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 == 2 {
		return resp.Header, nil
	} else {
		xml.Unmarshal(xml_result, e)
		return nil, e
	}
}

func (c *AliOSSClient) CreateObjectAcl(bucket string, key string, permission string) error {
	uri := fmt.Sprintf("/%s/%s?acl", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)
	query["acl"] = ""
	header["x-oss-object-acl"] = permission
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "PUT",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 2 {
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

func (c *AliOSSClient) GetObjectAcl(bucket string, key string) (*BucketACL, error) {
	uri := fmt.Sprintf("/%s/%s?acl", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)
	query["acl"] = ""
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	v := &BucketACL{}
	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 == 2 {
		xml.Unmarshal(xml_result, v)
		return v, nil
	} else {
		xml.Unmarshal(xml_result, e)
		return nil, e
	}
}

func (c *AliOSSClient) GetInitMultipartUpload(bucket string, key string) (*MultiUploadInit, error) {
	uri := fmt.Sprintf("/%s/%s?uploads", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)
	query["uploads"] = ""
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "POST",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	v := &MultiUploadInit{}
	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 == 2 {
		xml.Unmarshal(xml_result, v)
		return v, nil
	} else {
		xml.Unmarshal(xml_result, e)
		return nil, e
	}
}

func (c *AliOSSClient) UploadPart(bucket string, key string, part_number int, upload_id string, data []byte) error {
	uri := fmt.Sprintf("/%s", key)
	query := make(map[string]string)
	header := make(map[string]string)
	query["partNumber"] = fmt.Sprintf("%d", part_number)
	query["uploadId"] = upload_id
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "POST",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              bytes.NewReader(data),
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 2 {
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

func (c *AliOSSClient) CompleteUploadPart(bucket string, key string, upload_id string, part []PartUpload) error {
	data := CompleteUpload{Part: part}
	xml_content, _ := xml.MarshalIndent(data, "", "  ")
	uri := fmt.Sprintf("/%s/%s?uploads", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)
	query["uploadId"] = upload_id
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "POST",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              bytes.NewReader([]byte(xml.Header + string(xml_content))),
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 2 {
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

func (c *AliOSSClient) DeleteUploadPart(bucket string, key string, upload_id string) error {
	uri := fmt.Sprintf("/%s/%s", bucket, key)
	query := make(map[string]string)
	header := make(map[string]string)
	query["uploadId"] = upload_id
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "DELETE",
		Url:                  fmt.Sprintf("http://%s.%s/%s", bucket, c.EndPoint, key),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 2 {
		return nil
	} else {
		xml.Unmarshal(xml_result, e)
		return e
	}
}

func (c *AliOSSClient) ListMultiUploadPart(bucket string) (*MultiUploadList, error) {
	uri := fmt.Sprintf("/%s/?uploads", bucket)
	query := make(map[string]string)
	header := make(map[string]string)
	query["uploads"] = ""
	s := &oss_agent{
		AccessKey:            c.AccessKey,
		AccessKeySecret:      c.AccessKeySecret,
		Verb:                 "GET",
		Url:                  fmt.Sprintf("http://%s.%s", bucket, c.EndPoint),
		CanonicalizedHeaders: header,
		CanonicalizedUri:     uri,
		CanonicalizedQuery:   query,
		Content:              &bytes.Reader{},
		Debug:                c.Debug,
		logger:               c.logger,
	}

	v := &MultiUploadList{}
	e := &AliOssError{}
	resp, xml_result, err := s.send_request(false)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 == 2 {
		xml.Unmarshal(xml_result, v)
		return v, nil
	} else {
		xml.Unmarshal(xml_result, e)
		return nil, e
	}
}
