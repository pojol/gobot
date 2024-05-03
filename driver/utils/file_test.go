package utils

import (
	"fmt"
	"testing"
)

func TestGetDirectoryFiels(t *testing.T) {
	scripts := GetDirectoryFiels("/Users/pojol/github/gobot/script", ".lua")
	fmt.Println(scripts)
}
