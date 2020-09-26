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
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
)

var RequiredEnvVars = [...]string{
	"ANNOUNCED_ADDRESSES",
	"HTTP_RPC_URL",
	"KEEP_ETHEREUM_PASSWORD",
	"LOG_LEVEL",
	"PEERS",
	"WS_RPC_URL",
}

const TEMPLATE_PATH = "./config.toml.tmpl"
const CONFIG_PATH = "./config.toml"
const KEYSTORE = "/mnt/keystore"
var ACCOUNT_PATH = fmt.Sprintf("%s/keep_wallet.json", KEYSTORE)


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
	/*
	Use your dappnode's static ip address or dns name

	eg. "/ip4/80.20.40.233/tcp/3919"
	*/
	AnnouncedAddresses []string

	// 0x...
	EthAddress string

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

/*
createKeyStore creates a keystore with a single account
and then returns that account
*/
func createKeyStore(password string, accountPath string) (accounts.Account, error) {

	ks := keystore.NewKeyStore(KEYSTORE, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		return accounts.Account{}, err
	}

	err = os.Rename(account.URL.Path, accountPath)
	if err != nil {
		return accounts.Account{}, err
	}

	account.URL.Path = accountPath
	return account, nil
}

func loadKeyStore(password string, filePath string) (accounts.Account, error) {
	ks := keystore.NewKeyStore(KEYSTORE, keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return accounts.Account{}, err
	}

	account, err := ks.Import(jsonBytes, password, password)

	// NOTE: ErrAccountAlreadyExists hit only on linux/container (not on OS X)
	if err != nil && err != keystore.ErrAccountAlreadyExists {
		return accounts.Account{}, err
	}

	return account, nil
}

func main() {
	var announcedAddresses []string
	var httpRpcUrl string
	var peers []string
	var wsRpcUrl string
	var password string
	var account accounts.Account

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
			case "ANNOUNCED_ADDRESSES":
				announcedAddresses = strings.Split(value, ",")
			case "HTTP_RPC_URL":
				httpRpcUrl = value
			case "WS_RPC_URL":
				wsRpcUrl = value
			case "PEERS":
				peers = strings.Split(value, ",")
			case "KEEP_ETHEREUM_PASSWORD":
				password = value
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
		if blockNumber.Cmp(zero) != 1 {
			fmt.Println("Eth endpoint still synchronizing, waiting...")
			time.Sleep(10 * time.Second)
			continue
		}

		// load operator account (or create if necessary)
		if _, err := os.Stat(ACCOUNT_PATH); os.IsNotExist(err) {
			account, err = createKeyStore(password, ACCOUNT_PATH)
			if err != nil {
				fmt.Printf("Failed to create operator's ethereum account: %s\n", err)
				time.Sleep(10 * time.Second)
				continue
			}
			fmt.Printf("No existing account found. Created new operator account with address %s\n", account.Address.Hex())
		} else {
			account, err = loadKeyStore(password, ACCOUNT_PATH)
			if err != nil {
				fmt.Printf("Failed to load operator's ethereum account: %s\n", err)
				time.Sleep(10 * time.Second)
				continue
			}
			fmt.Printf("Loaded existing operator account with address: %s\n", account.Address.Hex())
		}

		break
	}

	config := Config{
		AnnouncedAddresses: announcedAddresses,
		EthAddress:    account.Address.Hex(),
		HttpRpcUrl: httpRpcUrl,
		Peers:      peers,
		WsRpcUrl:   wsRpcUrl,
	}
	renderTemplate(config)
}
