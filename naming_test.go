package main

import (
	"regexp"
	"testing"
)

func TestNameWithPrefix_Name(t *testing.T) {
	name := NameWithPrefix("full_name", "")
	if name != "full_name" {
		t.Fatalf("NameWithPrefix should return the name when name is present, name was %s", name)
	}
}

func TestNameWithPrefix_NameAndPrefix(t *testing.T) {
	name := NameWithPrefix("full_name", "prefix")
	if name != "full_name" {
		t.Fatalf("NameWithPrefix should return the name when name and prefix are present, name was %s", name)
	}
}

func TestNameWithPrefix_Prefix(t *testing.T) {
	r := regexp.MustCompile("prefix_[0-9]{26}")
	name := NameWithPrefix("", "prefix_")
	if !r.Match([]byte(name)) {
		t.Fatalf("NameWithPrefix should return a generated name with prefix when prefix is present, name was %s", name)
	}
}

func TestNameWithPrefix_NoNameAndNoPrefix(t *testing.T) {
	r := regexp.MustCompile("terraform-[0-9]{26}")
	name := NameWithPrefix("", "")
	if !r.Match([]byte(name)) {
		t.Fatalf("NameWithPrefix should return a generated name when prefix is present, name was %s", name)
	}
}

func TestNamePrefixFromName_ValidPrefixedName(t *testing.T) {
	prefix := NamePrefixFromName("test_role_20211103145836382400000001")
	if *prefix != "test_role_" {
		t.Fatalf("unexpected name prefix %s", *prefix)
	}
}

func TestNamePrefixFromName_InvalidPrefixedName(t *testing.T) {
	prefix := NamePrefixFromName("test_role")
	if prefix != nil {
		t.Fatal("expected prefix to be nil")
	}
}
