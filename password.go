// +build !darwin

package main

import "github.com/hashicorp/terraform/helper/schema"

func passwordRetrievalFunc(env_var string, dv interface{}) schema.SchemaDefaultFunc {
	return schema.EnvDefaultFunc(env_var, dv)
}
