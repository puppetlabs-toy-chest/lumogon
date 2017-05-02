package utils

import (
	"strings"
	"testing"
)

func Test_GetRandomName_hasPrefix(t *testing.T) {
	if !strings.HasPrefix(GetRandomName("myprefix_"), "myprefix_") {
		t.Error("Expected generated random name to have supplied prefix")
	}
}
func Test_GetRandomName_unique(t *testing.T) {
	prefix := "myprefix_"
	name1 := GetRandomName(prefix)
	name2 := GetRandomName(prefix)
	if name1 == name2 {
		t.Errorf("Generated random names are not random, %s == %s", name1, name2)
	}
}
