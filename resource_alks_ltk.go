package main

import (
	"context"
	"log"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlksLtk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlksLtkCreate,
		ReadContext:   resourceAlksLtkRead,
		UpdateContext: resourceAlksLtkUpdate,
		DeleteContext: resourceAlksLtkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"iam_username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"iam_user_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key": {
				Sensitive: true,
				Type:      schema.TypeString,
				Computed:  true,
			},
			"secret_key": {
				Sensitive: true,
				Type:      schema.TypeString,
				Computed:  true,
			},
			"tags":     TagsSchema(),
			"tags_all": TagsSchemaComputed(),
		},
		CustomizeDiff: customdiff.All(
			SetTagsDiff,
		),
	}
}

func resourceAlksLtkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS LTK User Create")

	var iamUsername = d.Get("iam_username").(string)
	var tags = d.Get("tags").(map[string]interface{})

	providerStruct := meta.(*AlksClient)
	client := providerStruct.client

	allTags := tagMapToSlice(combineTagMaps(providerStruct.defaultTags, tags))

	options := &alks.IamUserOptions{
		IamUserName: &iamUsername,
		Tags:        &allTags,
	}
	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateIamUser(options)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(iamUsername)
	_ = d.Set("iam_user_arn", resp.IAMUserArn)
	_ = d.Set("access_key", resp.AccessKey)
	_ = d.Set("secret_key", resp.SecretKey)

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	return resourceAlksLtkRead(ctx, d, meta)
}

func resourceAlksLtkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS LTK User Read")

	providerStruct := meta.(*AlksClient)
	client := providerStruct.client

	defaultTags := providerStruct.defaultTags
	ignoreTags := providerStruct.ignoreTags

	// Check if role exists.
	if d.Id() == "" || d.Id() == "none" {
		return nil
	}

	resp, err := client.GetIamUser(d.Id())

	if err != nil {
		//If error is 404, UserNotFound, we log it and let terraform decide how to handle it.
		//All other errors cause a failure
		if err.StatusCode == 404 {
			log.Printf("[Error] %s", err)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	_ = d.Set("iam_username", resp.User.UserName)
	_ = d.Set("access_key", resp.User.AccessKey)

	allTags := tagSliceToMap(resp.User.Tags)
	localTags := removeIgnoredTags(allTags, *ignoreTags)

	if err := d.Set("tags_all", localTags); err != nil {
		return diag.FromErr(err)
	}

	userSpecificTags := removeDefaultTags(localTags, defaultTags)

	if err := d.Set("tags", userSpecificTags); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAlksLtkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS LTK Update")

	// enable partial state mode
	d.Partial(true)

	if d.HasChange("tags_all") {
		// try updating enable_alks_access
		if err := updateUserTags(d, meta); err != nil {
			return diag.FromErr(err)
		}
	}

	d.Partial(false)

	return resourceAlksLtkRead(ctx, d, meta)
}

func resourceAlksLtkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS LTK User Delete")

	providerStruct := meta.(*AlksClient)
	client := providerStruct.client
	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	if _, err := client.DeleteIamUser(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updateUserTags(d *schema.ResourceData, meta interface{}) error {
	providerStruct := meta.(*AlksClient)
	client := providerStruct.client

	if err := validateIAMEnabled(client); err != nil {
		return err
	}

	//Do a read to get existing tags.  If any of those are in ignore_tags, then they are externally managed
	//and they should be included in the update so they don't get removed.
	resp, err := client.GetIamUser(d.Id())

	if err != nil {
		return err
	}

	existingTags := tagSliceToMap(resp.User.Tags)
	externalTags := getExternalyManagedTags(existingTags, *providerStruct.ignoreTags)
	internalTags := d.Get("tags_all").(map[string]interface{})

	//Tags includes default tags, role specific tags, and tags that exist externally on the role itself and are specified in ignored_tags
	tags := tagMapToSlice(combineTagMaps(internalTags, externalTags))

	options := alks.IamUserOptions{
		IamUserName: &resp.User.UserName,
		Tags:        &tags,
	}

	if _, err := client.UpdateIamUser(&options); err != nil {
		return err
	}
	return nil
}
