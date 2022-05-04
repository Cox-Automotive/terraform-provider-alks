package main

import (
	"reflect"
	"testing"

	"github.com/Cox-Automotive/alks-go"
)

func TestRemoveDefaultTags(t *testing.T) {
	cases := []struct {
		allTagsMap       map[string]interface{}
		defaultTagsSlice []alks.Tag
		expected         []alks.Tag
	}{
		{
			allTagsMap: map[string]interface{}{
				"resourceKey1": "resourceValue1",
				"defaultKey1":  "defaultValue1",
			},
			defaultTagsSlice: []alks.Tag{{Key: "defaultKey1", Value: "defaultValue1"}},
			expected:         []alks.Tag{{Key: "resourceKey1", Value: "resourceValue1"}},
		},
		{
			allTagsMap: map[string]interface{}{
				"defaultKey2": "defaultValue2",
				"defaultKey1": "resourceValue2",
			},
			defaultTagsSlice: []alks.Tag{
				{Key: "defaultKey2", Value: "defaultValue2"},
				{Key: "defaultKey1", Value: "defaultValue2"},
			},
			expected: []alks.Tag{
				{Key: "defaultKey1", Value: "resourceValue2"}, //Should not remove this key.  We are assuming that if the key matches one in default but not the value, that the default key was overwritten on purpose in the role definition and shouldnt be removed
			},
		},
	}

	for _, c := range cases {
		resourceTagsSlice := removeDefaultTags(c.allTagsMap, c.defaultTagsSlice)
		if !reflect.DeepEqual(resourceTagsSlice, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", resourceTagsSlice, c.expected)
		}
	}
}

func TestTagMapToSlice(t *testing.T) {
	cases := []struct {
		tagMap   map[string]interface{}
		expected []alks.Tag
	}{
		{
			tagMap: map[string]interface{}{
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
		expected map[string]interface{}
	}{
		{
			tagSlice: []alks.Tag{{Key: "defaultKey1", Value: "defaultValue1"}},
			expected: map[string]interface{}{"defaultKey1": "defaultValue1"},
		},
	}

	for _, c := range cases {
		tagMap := tagSliceToMap(c.tagSlice)
		if !reflect.DeepEqual(tagMap, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", tagMap, c.expected)
		}
	}
}

func TestCombineTagsWithDefault(t *testing.T) {
	cases := []struct {
		defaultTagSlice  []alks.Tag
		resourceTagSlice []alks.Tag
		expected         []alks.Tag
	}{
		{
			defaultTagSlice:  []alks.Tag{{Key: "defaultKey1", Value: "defaultValue1"}},
			resourceTagSlice: []alks.Tag{{Key: "defaultKey1", Value: "resourceValue1"}},
			expected:         []alks.Tag{{Key: "defaultKey1", Value: "resourceValue1"}},
		},
	}

	for _, c := range cases {
		tagSlice := combineTagsWithDefault(c.resourceTagSlice, c.defaultTagSlice)
		if !reflect.DeepEqual(tagSlice, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", tagSlice, c.expected)
		}
	}
}
