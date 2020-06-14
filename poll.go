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
	"context"
	"fmt"
	"math/big"
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

		var rpcUrl string
		for _, envVar := range currentEnv {
			if strings.HasPrefix(envVar, "HTTP_RPC_URL") {
				rpcUrl = strings.Split(envVar, "=")[1]
			}
		}

		// try connecting to http rpc url
		client, err := ethclient.Dial(rpcUrl)
		if err != nil {
			fmt.Printf("Failed to connect to ethereum endpoint: %s\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		fmt.Println("Successfully connected to eth1 rpc endpoint")

    // check block number to see if eth client is synced
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			fmt.Printf("Failed to retrieve latest block header %s\n", err)
      time.Sleep(10 * time.Second)
			continue
		}
		zero := big.NewInt(int64(0))
		blockNumber := header.Number
		if blockNumber.Cmp(zero) == 1 {
			fmt.Println("Eth endpoint is synchronized, starting keep client...")
			break
		}

		time.Sleep(10 * time.Second)
	}
}
