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
	"bufio"
	"context"
	"fmt"
	"html/template"
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

const TEMPLATE_PATH = "./config.toml.tmpl"
const CONFIG_PATH = "./config.toml"

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

type Config struct {
	// 0x...
	Address string

	// http://ropsten.dnp.dappnode.eth:8545/
	HttpRpcUrl string

	// ws://ropsten.dnp.dappnode.eth:8546/
	WsRpcUrl string

	/*
	[
	  "/dns4/bootstrap-1.core.keep.test.boar.network/tcp/3001/ipfs/16Uiu2HAkuTUKNh6HkfvWBEkftZbqZHPHi3Kak5ZUygAxvsdQ2UgG",
	  "/dns4/bootstrap-2.core.keep.test.boar.network/tcp/3001/ipfs/16Uiu2HAmQirGruZBvtbLHr5SDebsYGcq6Djw7ijF3gnkqsdQs3wK"
	]
	*/
	Peers []string
}

func renderTemplate(config Config) {
	t, err := template.ParseFiles(TEMPLATE_PATH)
	if err != nil {
    panic(err)
	}

	f, err := os.Create(CONFIG_PATH)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	t.Execute(w, config)
	w.Flush()
}


func main() {

	var address string
	var httpRpcUrl string
	var peers []string
	var wsRpcUrl string

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

		for _, envVarStr := range currentEnv {
			envVar := strings.Split(envVarStr, "=")
			key := envVar[0]
			value := envVar[1]
			switch key {
			case "ETHEREUM_ADDRESS":
				address = value
			case "HTTP_RPC_URL":
				httpRpcUrl = value
			case "WS_RPC_URL":
				wsRpcUrl = value
			case "PEERS":
				peers = strings.Split(value, ",")
			}
		}

		_, err := ethclient.Dial(wsRpcUrl)
		if err != nil {
			fmt.Printf("Failed to connect to ethereum endpoint using websockets: %s\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		client, err := ethclient.Dial(httpRpcUrl)
		if err != nil {
			fmt.Printf("Failed to connect to ethereum endpoint using http: %s\n", err)
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

	config := Config{
		Address: address,
		HttpRpcUrl: httpRpcUrl,
		Peers: peers,
		WsRpcUrl: wsRpcUrl,
	}
	renderTemplate(config)
}
