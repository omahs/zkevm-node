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
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

type genericCheckTestCase struct {
	TestName     string `json:"test_name"`
	Transactions []struct {
		Nonce int `json:"nonce"`
	} `json:"transactions"`
}

func loadGenericCheckTestCases(path string) ([]genericCheckTestCase, error) {
	var testCases []genericCheckTestCase

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

func TestNonceOrder(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	defer func() {
		require.NoError(t, operations.Teardown())
	}()
	opsCfg := operations.GetDefaultOperationsConfig()
	opsCfg.State.MaxCumulativeGasUsed = 80000000000
	opsman, err := operations.NewManager(ctx, opsCfg)
	require.NoError(t, err)
	err = opsman.Setup()
	require.NoError(t, err)

	testCases, err := loadGenericCheckTestCases("./../vectors/src/e2e/nonces.json")
	require.NoError(t, err)

	// use genesis address
	auth, err := operations.GetAuth("0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d", big.NewInt(1000))
	require.NoError(t, err)

	zkCounters := pool.ZkCounters{
		CumulativeGasUsed:    1000000,
		UsedKeccakHashes:     1,
		UsedPoseidonHashes:   1,
		UsedPoseidonPaddings: 1,
		UsedMemAligns:        1,
		UsedArithmetics:      1,
		UsedBinaries:         1,
		UsedSteps:            1,
	}
	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {

			log.Infof("Running test: %s", testCase.TestName)
			ctx := context.Background()

			s, err := pgpoolstorage.NewPostgresPoolStorage(dbutils.NewConfigFromEnv())
			require.NoError(t, err)

			p := pool.NewPool(s, opsman.State(), common.Address{})

			if testing.Short() {
				t.Skip()
			}

			defer func() {
				require.NoError(t, operations.Teardown())
			}()
			opsCfg := operations.GetDefaultOperationsConfig()
			opsCfg.State.MaxCumulativeGasUsed = 80000000000
			opsman, err := operations.NewManager(ctx, opsCfg)
			require.NoError(t, err)
			err = opsman.Setup()
			require.NoError(t, err)

			for _, tx := range testCase.Transactions {
				tx := types.NewTransaction(uint64(tx.Nonce), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10+int64(uint64(tx.Nonce))), []byte{})
				log.Infof("\nInserting: %d", tx.Nonce())
				signedTx, err := auth.Signer(auth.From, tx)
				require.NoError(t, err)
				if err := p.AddTx(ctx, *signedTx); err != nil {
					t.Error(err)
				}
				tx_got, err := p.GetTopPendingTxByProfitabilityAndZkCounters(ctx, zkCounters)
				require.NoError(t, err)
				//assert.Equal(t, uint64(5), tx.Nonce())
				p.DeleteTxsByHashes(ctx, []common.Hash{
					tx_got.Hash(),
				})
			}
			tx, err := p.GetTopPendingTxByProfitabilityAndZkCounters(ctx, zkCounters)
			require.Error(t, err)
			require.Nil(t, tx)
		})
	}

}
