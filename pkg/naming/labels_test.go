package naming

import (
	"testing"
)

func TestHelloName(t *testing.T) {
	t1 := make(map[string]string)
	t2 := make(map[string]string)
	t1["t1"] = "t1"
	t2["t2"] = "t2"
	if len(mergeLabels(nil, nil)) != 0 {
		t.Fatalf("merge of nil maps should be empty and not fail")
	}

	res := mergeLabels(t1, nil)
	if res["t1"] != "t1" {
		t.Fatalf("merge with nil map incorrect")
	}
	res = mergeLabels(nil, t2)
	if res["t2"] != "t2" {
		t.Fatalf("merge with nil map incorrect")
	}
	res = mergeLabels(t1, t2)
	if res["t1"] != "t1" || res["t2"] != "t2" {
		t.Fatalf("merge of maps incorrect")
	}
}
