package main

import (
	"log"
	"runtime"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/keybase/go-keychain"
)

func passwordRetrievalFunc(envVar string, dv interface{}) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if runtime.GOOS == "darwin" {
			log.Println("[INFO] On macOS so trying the keychain")
			query := keychain.NewItem()
			query.SetSecClass(keychain.SecClassGenericPassword)
			query.SetService("alkscli")
			query.SetAccount("alksuid")
			query.SetMatchLimit(keychain.MatchLimitOne)
			query.SetReturnData(true)
			results, err := keychain.QueryItem(query)
			if err != nil {
				log.Println("[WARN] Error accessing the macOS keychain. Falling back to environment variables")
				log.Println(err)
			} else {
				return string(results[0].Data), nil
			}
		}

		return schema.EnvDefaultFunc(envVar, dv)()
	}
}
