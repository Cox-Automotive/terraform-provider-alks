package main

import (
	"context"
	"fmt"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

//Removes default tags from a map of role specific + default tags
func removeDefaultTags(allTags map[string]interface{}, defalutTags []alks.Tag) []alks.Tag {
	for _, t := range defalutTags {
		delete(allTags, t.Key)
	}

	return tagMapToSlice(allTags)
}

func tagMapToSlice(tagMap map[string]interface{}) []alks.Tag {
	tags := []alks.Tag{}
	for k, v := range tagMap {
		tag := alks.Tag{Key: k, Value: v.(string)}
		tags = append(tags, tag)
	}
	return tags
}

func tagSliceToMap(tagSlice []alks.Tag) map[string]interface{} {
	tagMap := make(map[string]interface{})
	for _, t := range tagSlice {
		tagMap[t.Key] = t.Value
	}
	return tagMap
}

//Combines tags defined on an individual resource with the default tags listed on the provider block
//Resource specific tags will overwrite default tags
func combineTagsWithDefault(tags []alks.Tag, defaultTags []alks.Tag) []alks.Tag {
	defaultTagsMap := tagSliceToMap(defaultTags)

	for _, t := range tags {
		defaultTagsMap[t.Key] = t.Value
	}
	allTags := tagMapToSlice(defaultTagsMap)

	return allTags
}

func SetTagsDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	defaultTags := meta.(*AlksClient).defaultTags

	resourceTags := tagMapToSlice(diff.Get("tags").(map[string]interface{}))

	allTags := combineTagsWithDefault(resourceTags, defaultTags)

	// To ensure "tags_all" is correctly computed, we explicitly set the attribute diff
	// when the merger of resource-level tags onto provider-level tags results in n > 0 tags,
	// otherwise we mark the attribute as "Computed" only when their is a known diff (excluding an empty map)
	// or a change for "tags_all".

	if len(allTags) > 0 {
		if err := diff.SetNew("tags_all", tagSliceToMap(allTags)); err != nil {
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