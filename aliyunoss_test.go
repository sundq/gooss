package aliyunoss

import (
	// "fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func getAk() (string, string) {
	data, _ := ioutil.ReadFile("./data.txt")
	data_parse := strings.Split(string(data), ":")
	return data_parse[0], data_parse[1]
}

func TestListBucket(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, false)
	result, err := c.ListBucket("", "", 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}

// func TestCreateBucket(t *testing.T) {
// 	accesskey, access_key_secret := getAk()
// 	c := New(accesskey, access_key_secret, nil, false)
// 	err := c.CreateBucket("fff", "oss-cn-hangzhou", "")
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

func TestCreateObjectForBuff(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, true)
	err := c.CreateObjectForBuff("xiagnjun129866", "test44444444444", []byte("sssdddddddddyyy"), "")
	if err != nil {
		t.Error(err)
	} else {
		t.Log()
	}
}
