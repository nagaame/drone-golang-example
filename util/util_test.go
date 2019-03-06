package util


import (
	"testing"
)

func Test_GetNameLen_1(t *testing.T) {
	if l := GetNameLen("nick"); l != 4 {
		t.Error("test failed, the length of nick is not correct.")
	} else {
		t.Log("test passed.")
	}
}

func Test_GetNameLen_2(t *testing.T) {
	if l := GetNameLen(""); l != 0 {
		t.Error("test failed, the length of empty string is not correct.")
	} else {
		t.Log("test passed.")
	}
}