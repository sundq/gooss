package aliyunoss

import (
	// "fmt"
	"testing"
)

func TestListBucket(t *testing.T) {
	result, err := c.ListBucket("", "", 1)
	// l := result.Get("Buckets").GetIndex(0).MustMap()
	// fmt.Println(l["Name"])
	// err = err
	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}
