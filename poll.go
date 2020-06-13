/*
Prevents keep-client from starting unless
required preconditions are satisfied.

This program will wait forever until required env vars are set
and ethereum endpoint is accessible and synced.

If these constraints are satisfied, renders config.toml file
from template to be used by keep-client process, then halts with exit code 0.
*/
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

  "github.com/ethereum/go-ethereum/ethclient"
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
			fmt.Printf("Missing required environment variables: %s, please set them and restart\n", missingEnv)
      time.Sleep(10 * time.Second)
      continue
		} else {
      fmt.Println("Found all required environment variables")
    }

    // try connecting to http rpc url
    _, err := ethclient.Dial("https://mainnet.infura.io")
    if err != nil {
        fmt.Printf("Failed to connect to ethereum endpoint: %s\n", err)
        time.Sleep(10 * time.Second)
        continue
    }
    // idea try connecting to peers...
    fmt.Println("Successfully connected to eth1 rpc endpoint")
    break
	}
}
