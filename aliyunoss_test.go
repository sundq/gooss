package aliyunoss

import (
	// "fmt"
	"io/ioutil"
	"strings"
	"testing"
)

const (
	test_bucket        = "xiagnjun129866"
	test_object        = "test44444444444"
	test_object1       = "test_1"
	test_append_object = "test_2"
	test_content       = "aliyun oss test"
	test_create_bucket = "test-create"
)

func getAk() (string, string) {
	data, _ := ioutil.ReadFile("./data.txt")
	data_parse := strings.Split(string(data), ":")
	return data_parse[0], data_parse[1]
}

// func TestListBucket(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	result, err := c.ListBucket("", "", 1)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(result)
// 	}
// }

// func TestCreateBucket(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.CreateBucket(test_create_bucket, "oss-cn-hangzhou", "")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log()
// 	}
// }

// func TestDeleteBucket(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.DeleteBucket(test_create_bucket)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log()
// 	}
// }

// func TestListObject(t *testing.T) {

// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	result, err := c.ListObject("jiagouyun-cn", "", "", 0, "jakin-jiagouyun-sever/")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(result)
// 	}
// }

// func TestListBucketOfLocation(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	result, err := c.GetLocationOfBucket("xiagnjun129866")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(result)
// 	}
// }

// func TestCreateObjectForBuff(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.CreateObjectForBuff(test_bucket, test_object, []byte(test_content), "")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log()
// 	}
// }

// func TestCreateObjectForFile(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.CreateObjectForFile(test_bucket, test_object1, "test.txt", "")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log()
// 	}
// }

// func TestGetObjectAsBuff(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	result, err := c.GetObjectAsBuffer(test_bucket, test_object)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(string(result))
// 	}
// }

// func TestGetObjectAsFile(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.GetObjectAsFile(test_bucket, test_object, "test.txt")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log()
// 	}

// 	// defer os.Remove("test.txt")
// }

// func TestDeleteObject(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.DeleteObject(test_bucket, test_object)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log()
// 	}
// }

// func TestAppendObject(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	next_position, crc64, err := c.AppendObjectForBuff(test_bucket, test_append_object, 0, []byte(test_content))
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(next_position, crc64)
// 	}
// }

// func TestDeleteMultiObject(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, true)
// 	err := c.DeleteMultiObject(test_bucket, []string{"a", "v"})
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(err)
// 	}
// }

// func TestGetObjectInfo(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, true)
// 	result, err := c.GetObjectInfo(test_bucket, test_object)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(result)
// 	}
// }

// func TestModifyBucketAcl(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.ModifyBucketAcl(test_bucket, "public-read")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(err)
// 	}
// }

// func TestModifyBucketAcl1(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.ModifyBucketAcl(test_bucket, "")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(err)
// 	}
// }

// func TestOpenBucketLogging(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.OpenBucketLogging(test_bucket, test_bucket, "log-")
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(err)
// 	}
// }

// func TestGetBucketLogging(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	p, r, err := c.GetBucketLogging(test_bucket)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(p, r)
// 	}
// }

// func TestCloseBucketLogging(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.CloseBucketLogging(test_bucket)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(err)
// 	}
// }

func TestCreateBucketWebsite(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, true)
	err := c.CreateBucketWebsite(test_bucket, "index.html", "")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
}

func TestGetBucketWebsite(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, false)
	r, err := c.GetBucketWebsite(test_bucket)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(r)
	}
}

func TestDeleteBucketWebsite(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, false)
	err := c.DeleteBucketWebsite(test_bucket)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
}

func TestCreateBucketLifecycleRule(t *testing.T) {
	accesskey, access_key_secret := getAk()
	rule := []LifecycleRule{{Prefix: "test-", Status: "Disabled", Expiration: LifecycleRuleExpireDays{Days: 6}}}
	c := New(accesskey, access_key_secret, nil, true)
	err := c.CreateBucketLifecycleRule(test_bucket, rule)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
}

func TestGetBucketLifecycleRule(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, false)
	r, err := c.GetBucketLifecycleRule(test_bucket)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(r)
	}
}

func TestDeleteBucketLifecycleRule(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, false)
	err := c.DeleteBucketLifecycleRule(test_bucket)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
}

// func TestAddBucketRefer(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.AddBucketRefer(test_bucket, false, []string{"ss", "rrr"})
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(err)
// 	}
// }

// func TestGetBucketAcl(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	r, err := c.GetBucketAcl(test_bucket)
// 	if err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(r)
// 	}
// }
