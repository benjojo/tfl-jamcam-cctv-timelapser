package main

import (
	"os"
	"testing"
)

func TestXMLParse(t *testing.T) {
	f, err := os.ReadFile("./camera-list-example.xml")
	if err != nil {
		t.FailNow()
	}
	test := getCameraFromXML(f)

	if !test.Available || test.ID == "" {
		t.FailNow()
	}
}
