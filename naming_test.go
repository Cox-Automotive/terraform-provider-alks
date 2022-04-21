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

func TestValidRoleName_NameTooShort(t *testing.T) {
	name := ""
	warn, err := validRoleName(name, "")
	if len(warn) != 0 || len(err) != 2 {
		t.Fatalf("expected an error")
	}
}

func TestValidRoleName_NameMaxLength(t *testing.T) {
	name := "0123456789012345678901234567890123456789012345678901234567890123"
	warn, err := validRoleName(name, "")
	if len(warn) != 0 || len(err) != 0 {
		t.Fatalf("expected an error")
	}
}

func TestValidRoleName_NameTooLong(t *testing.T) {
	name := "01234567890123456789012345678901234567890123456789012345678901234"
	warn, err := validRoleName(name, "")
	if len(warn) != 0 || len(err) != 1 {
		t.Fatalf("expected an error")
	}
}

func TestValidRoleName_NameInvalidChar(t *testing.T) {
	name := "name!!!!@#$%^&*()-="
	warn, err := validRoleName(name, "")
	if len(warn) != 0 || len(err) != 1 {
		t.Fatalf("expected an error")
	}
}

func TestValidRolePrefix_NamePrefixTooShort(t *testing.T) {
	prefix := ""
	warn, err := validRolePrefix(prefix, "")
	if len(warn) != 0 || len(err) != 2 {
		t.Fatalf("expected an error")
	}
}

func TestValidRolePrefix_NamePrefixMaxLength(t *testing.T) {
	prefix := "01234567890123456789012345678901234567"
	warn, err := validRolePrefix(prefix, "")
	if len(warn) != 0 || len(err) != 0 {
		t.Fatalf("expected an error")
	}
}

func TestValidRolePrefix_NamePrefixTooLong(t *testing.T) {
	prefix := "012345678901234567890123456789012345678"
	warn, err := validRolePrefix(prefix, "")
	if len(warn) != 0 || len(err) != 1 {
		t.Fatalf("expected an error")
	}
}

func TestValidRolePrefix_NamePrefixInvalidChar(t *testing.T) {
	prefix := "name!!!!@#$%^&*()-="
	warn, err := validRolePrefix(prefix, "")
	if len(warn) != 0 || len(err) != 1 {
		t.Fatalf("expected an error")
	}
}