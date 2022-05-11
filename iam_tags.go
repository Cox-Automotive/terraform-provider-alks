package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TagMap map[string]interface{}
type IgnoreTags struct {
	Keys        TagMap
	KeyPrefixes TagMap
}

func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func TagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func SetTagsDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	defaultTags := meta.(*AlksClient).defaultTags
	ignoredTags := meta.(*AlksClient).ignoreTags
	resourceTags := (diff.Get("tags")).(map[string]interface{})
	//default tag values will be overwritten by resource values if key exists in both maps
	allTags := combineTagMaps(defaultTags, resourceTags)
	localTags := removeIgnoredTags(allTags, *ignoredTags)

	// To ensure "tags_all" is correctly computed, we explicitly set the attribute diff
	// when the merger of resource-level tags onto provider-level tags results in n > 0 tags,
	// otherwise we mark the attribute as "Computed" only when their is a known diff (excluding an empty map)
	// or a change for "tags_all".

	if len(localTags) > 0 {
		if err := diff.SetNew("tags_all", localTags); err != nil {
			return fmt.Errorf("error setting new tags_all diff: %w", err)
		}
	} else if len(diff.Get("tags_all").(map[string]interface{})) > 0 {
		if err := diff.SetNewComputed("tags_all"); err != nil {
			return fmt.Errorf("error setting tags_all to computed: %w", err)
		}
	} else if diff.HasChange("tags_all") {
		if err := diff.SetNewComputed("tags_all"); err != nil {
			return fmt.Errorf("error setting tags_all to computed: %w", err)
		}
	}

	return nil
}

//Removes default tags from a map of role specific + default tags
func removeDefaultTags(allTags TagMap, defaultTags TagMap) TagMap {
	for k, v := range defaultTags {
		//If the key and value of a tag returned from the role exists in the defaultTags list
		//We will assume it was set as a default tag and remove it from role specific tag list
		if val, ok := allTags[k]; ok {
			if val == v {
				delete(allTags, k)
			}
		}

	}
	return allTags
}

func removeIgnoredTags(allTags TagMap, ignoredTags IgnoreTags) TagMap {
	localMap := TagMap{}
	for k, v := range allTags {
		localMap[k] = v.(string)
	}

	for k := range allTags {
		if _, ok := ignoredTags.Keys[k]; ok {
			delete(localMap, k)
		} else {
			for kp := range ignoredTags.KeyPrefixes {
				if strings.HasPrefix(k, kp) {
					delete(localMap, k)
				}
			}
		}

	}
	return localMap
}

func tagMapToSlice(tagMap TagMap) []alks.Tag {
	tags := []alks.Tag{}
	for k, v := range tagMap {
		tag := alks.Tag{Key: k, Value: v.(string)}
		tags = append(tags, tag)
	}
	return tags
}

func tagSliceToMap(tagSlice []alks.Tag) TagMap {
	tagMap := make(TagMap)
	for _, t := range tagSlice {
		tagMap[t.Key] = t.Value
	}
	return tagMap
}

func getExternalyManagedTags(roleTags TagMap, ignoredTags IgnoreTags) TagMap {
	externalTags := TagMap{}
	//Loop Through ignored keys and ignored key prefixes, checking if a tag exists that should be ignored
	for k := range ignoredTags.Keys {
		if val, ok := roleTags[k]; ok {
			externalTags[k] = val.(string)
		}
	}

	for p := range ignoredTags.KeyPrefixes {
		for k, v := range roleTags {
			if strings.HasPrefix(k, p) {
				externalTags[k] = v.(string)
			}
		}
	}
	return externalTags
}

//Combine two tag maps.  Values in map2 will overwrite values in map1 if they exist in both maps
func combineTagMaps(map1 TagMap, map2 TagMap) TagMap {
	LocalMap := TagMap{}

	for k, v := range map1 {
		LocalMap[k] = v
	}
	for k, v := range map2 {
		LocalMap[k] = v
	}

	return LocalMap
}
