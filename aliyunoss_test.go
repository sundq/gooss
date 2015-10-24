package aliyunoss

import (
	"fmt"
	"testing"
)

func TestListBucket(t *testing.T) {
	c := New("C59MdXDbgVItTf2O", "z3EAIwHJ4FHGBpDW84j7iW2eSWilZQ", nil, nil)
	cc, result, err := c.ListBucket("xiao", "", 0)
	l := result.Get("Buckets").GetIndex(0).MustMap()
	fmt.Println(l["Name"])
	cc = cc
	err = err
	t.Log(result)
}
