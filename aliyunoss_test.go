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

func TestAppendObject(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, false)
	next_position, crc64, err := c.AppendObjectForBuff(test_bucket, test_append_object, 0, []byte(test_content))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(next_position, crc64)
	}
}
