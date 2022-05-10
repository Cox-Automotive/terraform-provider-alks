package main

import (
	"reflect"
	"testing"

	"github.com/Cox-Automotive/alks-go"
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

func TestCombineMaps(t *testing.T) {
	cases := []struct {
		defaultTagSlice  TagMap
		resourceTagSlice TagMap
		expected         TagMap
	}{
		{
			defaultTagSlice:  TagMap{"defaultKey1": "defaultValue1"},
			resourceTagSlice: TagMap{"defaultKey1": "resourceValue1"},
			expected:         TagMap{"defaultKey1": "resourceValue1"},
		},
	}

	for _, c := range cases {
		tagMap := combineTagMaps(c.defaultTagSlice, c.resourceTagSlice)
		if !reflect.DeepEqual(tagMap, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", tagMap, c.expected)
		}
	}
}
