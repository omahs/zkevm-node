package e2e

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

type Account struct {
	Auth *bind.TransactOpts

	Name           string         `json:"name"`
	PrivateKey     string         `json:"private_key"`
	Address        common.Address `json:"address"`
	InitialBalance *big.Int       `json:"balance"`
	AmountToSend   *big.Int       `json:"amount_to_send"`
}

func loadAccountsAsTestCases(path string) ([]Account, error) {
	var testCases []Account

	jsonFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return testCases, err
	}
	defer func() { _ = jsonFile.Close() }()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return testCases, err
	}

	err = json.Unmarshal(bytes, &testCases)
	if err != nil {
		return testCases, err
	}

	return testCases, nil
}
func Test_FullBalanceInTransit(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// Setup

	defer func() {
		//	require.NoError(t, operations.Teardown())
	}()

	log.Info(">>> Setting up...")
	opsCfg := operations.GetDefaultOperationsConfig()
	opsCfg.State.MaxCumulativeGasUsed = 80000000000
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsman.Setup()
	require.NoError(t, err)

	// Load eth client
	client, err := ethclient.Dial("http://localhost:8123")
	require.NoError(t, err)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(1000000)
	chainID := big.NewInt(1000)
	log.Info(">>> Connected to ZKEVM-Node ETH client...")

	// Hydrate accounts

	accounts, err := loadAccountsAsTestCases("./../vectors/src/e2e/balance_in_transit.json")
	require.NoError(t, err)

	for index, account := range accounts {
		auth, err := operations.GetAuth(account.PrivateKey, chainID)
		require.NoError(t, err)
		accounts[index].Auth = auth
	}
	log.Info(">>> Accounts hydrated with their respective Auth objects...")

	// Account 0 => Account 1
	// Account 2 => Account 0
	for index, account := range accounts {
		log.Infof(">>> Starting with: %s", account.Name)
		balance, err := client.BalanceAt(ctx, account.Address, nil)
		require.NoError(t, err)
		require.Equal(t, account.InitialBalance, balance)
		nonce := uint64(index)

		if index == len(accounts)-1 {
			log.Infof("\t>>> Last TX: Account %s has %d - will send %d to %s", account.Name, balance, account.AmountToSend, accounts[0].Address)

			tx := types.NewTransaction(nonce, accounts[0].Address, account.AmountToSend, gasLimit, gasPrice, nil)
			signedTx, err := account.Auth.Signer(account.Auth.From, tx)
			require.NoError(t, err)
			err = client.SendTransaction(context.Background(), signedTx)
			require.NoError(t, err)
			break
		}

		log.Infof("\t>>> Account %s has %d - will send %d to %s", account.Name, balance, account.AmountToSend, accounts[index+1].Name)

		tx := types.NewTransaction(nonce, accounts[index+1].Address, account.AmountToSend, gasLimit, gasPrice, nil)
		signedTx, err := account.Auth.Signer(account.Auth.From, tx)
		require.NoError(t, err)
		err = client.SendTransaction(context.Background(), signedTx)
		require.NoError(t, err)
		log.Info(">>> Sent")
	}
}
