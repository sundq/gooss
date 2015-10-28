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

func TestListObject(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, false)
	result, err := c.ListObject("xiagnjun129866")
	// for index, object := range result.Objects {
	// 	fmt.Println(fmt.Sprintf("%d:%s %s %s", index, object.Key, object.Size, object.Type))
	// }
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result.Objects)
	}
}

func TestListBucketOfLocation(t *testing.T) {
	accesskey, access_key_secret := getAk()
	c := New(accesskey, access_key_secret, nil, false)
	result, err := c.GetLocationOfBucket("xiagnjun129866")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}
