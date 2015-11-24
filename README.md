# aliyunoss
--
    import "github.com/sundq/aliyunoss"


## Usage

#### type AccessControl

```go
type AccessControl struct {
	Grant string `xml:"Grant"`
}
```


#### type AliOSSClient

```go
type AliOSSClient struct {
	AccessKey       string
	AccessKeySecret string
	EndPoint        string
	Debug           bool
}
```

#### type AliOssError

```go
type AliOssError struct {
	Code           string `xml:"Code"`
	Message        string `xml:"Message"`
	RequestId      string `xml:"RequestId"`
	HostId         string `xml:"HostId"`
	OSSAccessKeyId string `xml:"OSSAccessKeyId"`
}
```


#### func (AliOssError) Error

```go
func (e AliOssError) Error() string
```

#### type Bucket

```go
type Bucket struct {
	Location     string `xml:"Location"`
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}
```


#### type BucketACL

```go
type BucketACL struct {
	XMLName           xml.Name      `xml:"AccessControlPolicy"`
	AccessControlList AccessControl `xml:AccessControlList`
}
```


#### type BucketConfiguration

```go
type BucketConfiguration struct {
	XMLName            xml.Name `xml:"CreateBucketConfiguration"`
	LocationConstraint string   `xml:LocationConstraint`
}
```


#### type BucketList

```go
type BucketList struct {
	// XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}
```


#### type BucketLog

```go
type BucketLog struct {
	XMLName        xml.Name       `xml:"BucketLoggingStatus"`
	LoggingEnabled LoggingEnabled `xml:LoggingEnabled`
}
```


#### type BucketLogging

```go
type BucketLogging struct {
	XMLName            xml.Name `xml:"BucketLoggingStatus"`
	LocationConstraint string   `xml:LocationConstraint`
}
```


#### type BucketWebsite

```go
type BucketWebsite struct {
	XMLName       xml.Name                `xml:"WebsiteConfiguration"`
	IndexDocument BucketWebsiteIndex      `xml:"IndexDocument"`
	ErrorDocument BucketWebsiteErrorIndex `xml:"ErrorDocument"`
}
```


#### type BucketWebsiteErrorIndex

```go
type BucketWebsiteErrorIndex struct {
	Key string `xml:"Key"`
}
```


#### type BucketWebsiteIndex

```go
type BucketWebsiteIndex struct {
	Suffix string `xml:"Suffix"`
}
```


#### type CompleteUpload

```go
type CompleteUpload struct {
	XMLName xml.Name     `xml:"CompleteMultipartUpload"`
	Part    []PartUpload `xml:"Part"`
}
```


#### type DeleteObject

```go
type DeleteObject struct {
	Object string `xml:"Key"`
}
```


#### type DeleteObjectList

```go
type DeleteObjectList struct {
	XMLName xml.Name       `xml:"Delete"`
	Quiet   string         `xml:"Quiet"`
	Objects []DeleteObject `xml:"Object"`
}
```


#### type Lifecycle

```go
type Lifecycle struct {
	XMLName xml.Name        `xml:"LifecycleConfiguration"`
	Rule    []LifecycleRule `xml:Rule`
}
```


#### type LifecycleRule

```go
type LifecycleRule struct {
	RuleID     string      `xml:Rule`
	Prefix     string      `xml:Prefix`
	Status     string      `xml:Status` //Disabled Enabled
	Expiration interface{} `xml:Expiration`
}
```


#### type LifecycleRuleExpireDay

```go
type LifecycleRuleExpireDay struct {
	Day string `xml:Day`
}
```


#### type LifecycleRuleExpireDays

```go
type LifecycleRuleExpireDays struct {
	Days int `xml:Days`
}
```


#### type LoggingEnabled

```go
type LoggingEnabled struct {
	TargetBucket string `xml:"TargetBucket"`
	TargetPrefix string `xml:"TargetPrefix"`
}
```


#### type MultiUpload

```go
type MultiUpload struct {
	Key       string `xml:"Key"`
	UploadId  string `xml:"UploadId"`
	Initiated string `xml:"Initiated"`
}
```


#### type MultiUploadInit

```go
type MultiUploadInit struct {
	XMLName  xml.Name `xml:"InitiateMultipartUploadResult"`
	Bucket   string   `xml:"Bucket"`
	Key      string   `xml:"Key"`
	UploadId string   `xml:"UploadId"`
}
```


#### type MultiUploadList

```go
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
```


#### type Object

```go
type Object struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Type         string `xml:"Type"`
	Size         string `xml:"Size"`
}
```


#### type ObjectList

```go
type ObjectList struct {
	BucketName string   `xml:"Name"`
	Prefix     string   `xml:"Prefix"`
	Marker     string   `xml:"Marker"`
	Delimiter  string   `xml:"Delimiter"`
	MaxKeys    string   `xml:"MaxKeys"`
	Objects    []Object `xml:"Contents"`
}
```


#### type PartUpload

```go
type PartUpload struct {
	PartNumber int    `xml:"PartNumber"`
	ETag       string `xml:"ETag"`
}
```


#### func  New

```go
func New(access_key string, access_key_secret string, endpoint interface{}, debug interface{}) *AliOSSClient
```

#### func (*AliOSSClient) AddBucketRefer

```go
func (c *AliOSSClient) AddBucketRefer(name string, allow_empty_referer bool, refer_list []string) error
```

#### func (*AliOSSClient) AppendObjectForBuff

```go
func (c *AliOSSClient) AppendObjectForBuff(bucket string, key string, position int, data []byte) (int, string, error)
```
AppendObjectForBuff like CreateObjectForBuff but it will append the exist key

#### func (*AliOSSClient) CloseBucketLogging

```go
func (c *AliOSSClient) CloseBucketLogging(name string) error
```

#### func (*AliOSSClient) CompleteUploadPart

```go
func (c *AliOSSClient) CompleteUploadPart(bucket string, key string, upload_id string, part []PartUpload) error
```
tell oss a multipart upload have been completed

#### func (*AliOSSClient) CreateBucket

```go
func (c *AliOSSClient) CreateBucket(name string, location string, permission string) error
```

#### func (*AliOSSClient) CreateBucketLifecycleRule

```go
func (c *AliOSSClient) CreateBucketLifecycleRule(name string, rule_list []LifecycleRule) error
```

#### func (*AliOSSClient) CreateBucketWebsite

```go
func (c *AliOSSClient) CreateBucketWebsite(name string, index string, error_file string) error
```

#### func (*AliOSSClient) CreateObjectAcl

```go
func (c *AliOSSClient) CreateObjectAcl(bucket string, key string, permission string) error
```
set object access control

#### func (*AliOSSClient) CreateObjectForBuff

```go
func (c *AliOSSClient) CreateObjectForBuff(bucket string, key string, data []byte, permission string) error
```
CreateObjectForBuff create a oss key for buffer, the permission can be private,
public-read or public_write

#### func (*AliOSSClient) CreateObjectForFile

```go
func (c *AliOSSClient) CreateObjectForFile(bucket string, key string, filepath string, permission string) error
```
CreateObjectForBuff create a oss key for local file, the permission can be
private, public-read or public_write

#### func (*AliOSSClient) DeleteBucket

```go
func (c *AliOSSClient) DeleteBucket(name string) error
```
Delete bucket

#### func (*AliOSSClient) DeleteBucketLifecycleRule

```go
func (c *AliOSSClient) DeleteBucketLifecycleRule(name string) error
```

#### func (*AliOSSClient) DeleteBucketWebsite

```go
func (c *AliOSSClient) DeleteBucketWebsite(name string) error
```

#### func (*AliOSSClient) DeleteMultiObject

```go
func (c *AliOSSClient) DeleteMultiObject(bucket string, keys []string) error
```
DeleteMultiObject delete multi-object

#### func (*AliOSSClient) DeleteObject

```go
func (c *AliOSSClient) DeleteObject(bucket string, key string) error
```
DeleteObject delete single key

#### func (*AliOSSClient) DeleteUploadPart

```go
func (c *AliOSSClient) DeleteUploadPart(bucket string, key string, upload_id string) error
```
delete a multipart upload

#### func (*AliOSSClient) GetBucketAcl

```go
func (c *AliOSSClient) GetBucketAcl(name string) (string, error)
```

#### func (*AliOSSClient) GetBucketLifecycleRule

```go
func (c *AliOSSClient) GetBucketLifecycleRule(name string) (*Lifecycle, error)
```

#### func (*AliOSSClient) GetBucketLogging

```go
func (c *AliOSSClient) GetBucketLogging(name string) (string, string, error)
```

#### func (*AliOSSClient) GetBucketWebsite

```go
func (c *AliOSSClient) GetBucketWebsite(name string) (*BucketWebsite, error)
```

#### func (*AliOSSClient) GetInitMultipartUpload

```go
func (c *AliOSSClient) GetInitMultipartUpload(bucket string, key string) (*MultiUploadInit, error)
```
Init multipart upload

#### func (*AliOSSClient) GetLocationOfBucket

```go
func (c *AliOSSClient) GetLocationOfBucket(bucket string) (string, error)
```

#### func (*AliOSSClient) GetObjectAcl

```go
func (c *AliOSSClient) GetObjectAcl(bucket string, key string) (*BucketACL, error)
```
get object access control

#### func (*AliOSSClient) GetObjectAsBuffer

```go
func (c *AliOSSClient) GetObjectAsBuffer(bucket string, key string) ([]byte, error)
```
GetObjectAsBuffer get object as buffer

#### func (*AliOSSClient) GetObjectAsFile

```go
func (c *AliOSSClient) GetObjectAsFile(bucket string, key string, filepath string) error
```
GetObjectAsFile get a object as local file

#### func (*AliOSSClient) GetObjectInfo

```go
func (c *AliOSSClient) GetObjectInfo(bucket string, key string) (http.Header, error)
```
GetObjectInfo information of object

#### func (*AliOSSClient) GetObjectMetaData

```go
func (c *AliOSSClient) GetObjectMetaData(bucket string, key string) (http.Header, error)
```
GetObjectMetaData information of object

#### func (*AliOSSClient) ListBucket

```go
func (c *AliOSSClient) ListBucket(prefix string, marker string, max_size int) (*BucketList, error)
```

#### func (*AliOSSClient) ListMultiUploadPart

```go
func (c *AliOSSClient) ListMultiUploadPart(bucket string) (*MultiUploadList, error)
```
list multiupload

#### func (*AliOSSClient) ListObject

```go
func (c *AliOSSClient) ListObject(bucket string, delimiter string, marker string, max_size int, prefix string) (*ObjectList, error)
```
ListObject get the list of key for the specified bucket

#### func (*AliOSSClient) ModifyBucketAcl

```go
func (c *AliOSSClient) ModifyBucketAcl(name string, permission string) error
```

#### func (*AliOSSClient) OpenBucketLogging

```go
func (c *AliOSSClient) OpenBucketLogging(name string, target_bucket string, obj_prefix string) error
```

#### func (*AliOSSClient) UploadPart

```go
func (c *AliOSSClient) UploadPart(bucket string, key string, part_number int, upload_id string, data []byte) error
```
upload a part for multiupload


