package main

import (
	"reflect"
	"testing"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/go-cty/cty"
)

func TestRemoveDefaultTags(t *testing.T) {
	cases := []struct {
		allTagsMap    TagMap
		defaultTagMap TagMap
		expected      TagMap
	}{
		{
			allTagsMap: TagMap{
				"resourceKey1": "resourceValue1",
				"defaultKey1":  "defaultValue1",
			},
			defaultTagMap: TagMap{"defaultKey1": "defaultValue1"},
			expected:      TagMap{"resourceKey1": "resourceValue1"},
		},
		{
			allTagsMap: TagMap{
				"defaultKey2": "defaultValue2",
				"defaultKey1": "resourceValue2",
			},
			defaultTagMap: TagMap{
				"defaultKey2": "defaultValue2",
				"defaultKey1": "defaultValue2",
			},
			expected: TagMap{
				"defaultKey1": "resourceValue2", //Should not remove this key.  We are assuming that if the key matches one in default but not the value, that the default key was overwritten on purpose in the role definition and shouldnt be removed
			},
		},
	}

	for _, c := range cases {
		resourceTagsSlice := removeDefaultTags(c.allTagsMap, c.defaultTagMap)
		if !reflect.DeepEqual(resourceTagsSlice, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", resourceTagsSlice, c.expected)
		}
	}
}

func TestRemoveDefaultTags_NilDefault(t *testing.T) {
	allTags := TagMap{"key1": "val1"}
	result := removeDefaultTags(allTags, nil)
	if !reflect.DeepEqual(result, allTags) {
		t.Fatalf("expected allTags unchanged when defaultTags is nil, got %#v", result)
	}
}

func TestRemoveIgnoredTags(t *testing.T) {
	cases := []struct {
		allTags     TagMap
		ignoredTags IgnoreTags
		expected    TagMap
	}{
		{
			allTags: TagMap{
				"Key1":           "Value1",
				"Key2":           "Value2",
				"KeyPrefix:Key3": "Value3",
			},
			ignoredTags: IgnoreTags{
				Keys:        TagMap{"Key1": ""},
				KeyPrefixes: TagMap{"KeyPrefix:": ""},
			},
			expected: TagMap{"Key2": "Value2"},
		},
	}

	for _, c := range cases {
		tagMap := removeIgnoredTags(c.allTags, c.ignoredTags)
		if !reflect.DeepEqual(tagMap, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", tagMap, c.expected)
		}
	}
}

func TestTagMapToSlice(t *testing.T) {
	cases := []struct {
		tagMap   TagMap
		expected []alks.Tag
	}{
		{
			tagMap: TagMap{
				"key1": "value1",
			},
			expected: []alks.Tag{{Key: "key1", Value: "value1"}},
		},
	}

	for _, c := range cases {
		tagSlice := tagMapToSlice(c.tagMap)
		if !reflect.DeepEqual(tagSlice, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", tagSlice, c.expected)
		}
	}
}

func TestTagSliceToMap(t *testing.T) {
	cases := []struct {
		tagSlice []alks.Tag
		expected TagMap
	}{
		{
			tagSlice: []alks.Tag{{Key: "defaultKey1", Value: "defaultValue1"}},
			expected: TagMap{"defaultKey1": "defaultValue1"},
		},
	}

	for _, c := range cases {
		tagMap := tagSliceToMap(c.tagSlice)
		if !reflect.DeepEqual(tagMap, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", tagMap, c.expected)
		}
	}
}

func TestGetExternalyManagedTags(t *testing.T) {
	cases := []struct {
		roleTags    TagMap
		ignoredTags IgnoreTags
		expected    TagMap
	}{
		{
			roleTags: TagMap{
				"Key1":           "Value1",
				"Key2":           "Value2",
				"KeyPrefix:Key3": "Value3",
			},
			ignoredTags: IgnoreTags{
				Keys:        TagMap{"Key1": ""},
				KeyPrefixes: TagMap{"KeyPrefix:": ""},
			},
			expected: TagMap{
				"Key1":           "Value1",
				"KeyPrefix:Key3": "Value3",
			},
		},
	}

	for _, c := range cases {
		tagMap := getExternalyManagedTags(c.roleTags, c.ignoredTags)
		if !reflect.DeepEqual(tagMap, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", tagMap, c.expected)
		}
	}
}

func TestCombineMaps(t *testing.T) {
	cases := []struct {
		defaultTagMap  TagMap
		resourceTagMap TagMap
		expected       TagMap
	}{
		{
			defaultTagMap:  TagMap{"defaultKey1": "defaultValue1"},
			resourceTagMap: TagMap{"defaultKey1": "resourceValue1"},
			expected:       TagMap{"defaultKey1": "resourceValue1"},
		},
	}

	for _, c := range cases {
		tagMap := combineTagMaps(c.defaultTagMap, c.resourceTagMap)
		if !reflect.DeepEqual(tagMap, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", tagMap, c.expected)
		}
	}
}

func TestCombineMaps_BothEmpty(t *testing.T) {
	result := combineTagMaps(TagMap{}, TagMap{})
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %#v", result)
	}
}

func TestCombineMaps_MergeMap2IntoMap1(t *testing.T) {
	base := TagMap{"k1": "v1", "k2": "v2"}
	merge := TagMap{"k2": "override", "k3": "v3"}
	result := combineTagMaps(base, merge)
	if result["k1"] != "v1" {
		t.Fatalf("expected k1=v1, got %s", result["k1"])
	}
	if result["k2"] != "override" {
		t.Fatalf("expected k2=override, got %s", result["k2"])
	}
	if result["k3"] != "v3" {
		t.Fatalf("expected k3=v3, got %s", result["k3"])
	}
}

func TestResolveDuplicates_RemovesDefaultTagsBasicPath(t *testing.T) {
	allTags := TagMap{
		"resourceKey": "resourceValue",
		"defaultKey":  "defaultValue",
	}
	defaultTags := TagMap{"defaultKey": "defaultValue"}

	d := resourceAlksIamRole().Data(nil)
	result := resolveDuplicates(allTags, defaultTags, d)

	if _, ok := result["defaultKey"]; ok {
		t.Fatal("defaultKey should have been removed by removeDefaultTags")
	}
	if result["resourceKey"] != "resourceValue" {
		t.Fatalf("resourceKey should be preserved, got %v", result["resourceKey"])
	}
}

func TestResolveDuplicates_NilDefaultTags(t *testing.T) {
	allTags := TagMap{"key1": "value1", "key2": "value2"}
	d := resourceAlksIamRole().Data(nil)

	result := resolveDuplicates(allTags, nil, d)

	if result["key1"] != "value1" {
		t.Fatalf("expected key1=value1, got %v", result["key1"])
	}
	if result["key2"] != "value2" {
		t.Fatalf("expected key2=value2, got %v", result["key2"])
	}
}

func TestResolveDuplicates_EmptyAllTags(t *testing.T) {
	d := resourceAlksIamRole().Data(nil)
	result := resolveDuplicates(TagMap{}, TagMap{"k": "v"}, d)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %#v", result)
	}
}

func TestNormalizeTagsFromRaw_AddsTags(t *testing.T) {
	m := map[string]cty.Value{
		"key1": cty.StringVal("value1"),
		"key2": cty.StringVal("value2"),
	}
	configTags := &TagMap{}

	normalizeTagsFromRaw(m, configTags)

	if (*configTags)["key1"] != "value1" {
		t.Fatalf("expected key1=value1, got %v", (*configTags)["key1"])
	}
	if (*configTags)["key2"] != "value2" {
		t.Fatalf("expected key2=value2, got %v", (*configTags)["key2"])
	}
}

func TestNormalizeTagsFromRaw_SkipsNullValues(t *testing.T) {
	m := map[string]cty.Value{
		"nullKey":   cty.NullVal(cty.String),
		"normalKey": cty.StringVal("normalValue"),
	}
	configTags := &TagMap{}

	normalizeTagsFromRaw(m, configTags)

	if _, ok := (*configTags)["nullKey"]; ok {
		t.Fatal("nullKey should have been skipped")
	}
	if (*configTags)["normalKey"] != "normalValue" {
		t.Fatalf("expected normalKey=normalValue, got %v", (*configTags)["normalKey"])
	}
}

func TestNormalizeTagsFromRaw_DoesNotOverwriteExisting(t *testing.T) {
	m := map[string]cty.Value{
		"existingKey": cty.StringVal("newValue"),
		"newKey":      cty.StringVal("newKeyValue"),
	}
	configTags := &TagMap{"existingKey": "originalValue"}

	normalizeTagsFromRaw(m, configTags)

	if (*configTags)["existingKey"] != "originalValue" {
		t.Fatalf("expected existingKey to retain originalValue, got %v", (*configTags)["existingKey"])
	}
	if (*configTags)["newKey"] != "newKeyValue" {
		t.Fatalf("expected newKey=newKeyValue, got %v", (*configTags)["newKey"])
	}
}

func TestNormalizeTagsFromRaw_EmptyMap(t *testing.T) {
	configTags := &TagMap{}
	normalizeTagsFromRaw(map[string]cty.Value{}, configTags)
	if len(*configTags) != 0 {
		t.Fatalf("expected empty configTags, got %#v", *configTags)
	}
}
