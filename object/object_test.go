package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello"}
	hello2 := &String{Value: "Hello"}
	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have differnt hash keys")
	}

	diff1 := &String{Value: "My name is something"}
	diff2 := &String{Value: "My name is something"}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have differnt hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
