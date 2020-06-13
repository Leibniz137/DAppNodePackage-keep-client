package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var RequiredEnvVars = [...]string{
	"ETHEREUM_ADDRESS",
	"HTTP_RPC_URL",
	"KEEP_ETHEREUM_PASSWORD",
	"LOG_LEVEL",
	"PEERS",
	"WS_RPC_URL",
}

func checkEnvVars(currentEnv []string) []string {

	errors := make([]string, 0, len(RequiredEnvVars))

  var found bool
  var empty bool
  var envValue []string
	for _, requiredEnvVar := range RequiredEnvVars {
    found = false
    empty = false
		for _, envVar := range currentEnv {
			if strings.HasPrefix(envVar, requiredEnvVar) {
				found = true
        envValue = strings.Split(envVar, "=")
        if strings.TrimSpace(envValue[1]) == "" {
          fmt.Println("empty value")
          empty = true
        }
			}
		}

		if !found || empty {
			errors = append(errors, requiredEnvVar)
		}
	}
	return errors
}

func main() {
	currentEnv := os.Environ()

  var missingEnv []string
	for true {
    missingEnv = checkEnvVars(currentEnv)
		if len(missingEnv) != 0 {
			// TODO: log which vars are missing
			fmt.Printf("Missing required environment variables: %s, please set them and restart\n", missingEnv)
		} else {
      fmt.Println("Found all required environment variables")
    }
		time.Sleep(10 * time.Second)
	}
}
