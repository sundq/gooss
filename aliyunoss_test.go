package aliyunoss

import (
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
