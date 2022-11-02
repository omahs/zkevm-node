package aggregator2

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// Consumer interfaces required by the package.

type proverInterface interface {
	ID() string
	IsIdle() bool
	AggregateProofs(ctx context.Context) error
	VerifyBatch(ctx context.Context) error
}

// ethTxManager contains the methods required to send txs to
// ethereum.
type ethTxManager interface {
	VerifyBatch(ctx context.Context, batchNum uint64, proof *pb.GetProofResponse) error
}

// etherman contains the methods required to interact with ethereum
type etherman interface {
	GetLatestVerifiedBatchNum() (uint64, error)
	GetPublicAddress() common.Address
}

// aggregatorTxProfitabilityChecker interface for different profitability
// checking algorithms.
type aggregatorTxProfitabilityChecker interface {
	IsProfitable(context.Context, *big.Int) (bool, error)
}

// proverClient is a wrapper to the prover service
type proverClientInterface interface {
	GetURI() string
	IsIdle(ctx context.Context) bool
	GetGenProofID(ctx context.Context, inputProver *pb.InputProver) (string, error)
	GetResGetProof(ctx context.Context, genProofID string, batchNumber uint64) (*pb.GetProofResponse, error)
}

// stateInterface gathers the methods to interact with the state.
type stateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetVirtualBatchToProve(ctx context.Context, lastVerfiedBatchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	GetProofsToAggregate(ctx context.Context, dbTx pgx.Tx) (*state.Proof2, *state.Proof2, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	AddGeneratedProof2(ctx context.Context, proof *state.Proof2, dbTx pgx.Tx) error
	UpdateGeneratedProof2(ctx context.Context, proof *state.Proof2, dbTx pgx.Tx) error
	DeleteGeneratedProof2(ctx context.Context, batchNumber uint64, batchNumberFinal uint64, dbTx pgx.Tx) error
	DeleteUngeneratedProofs2(ctx context.Context, dbTx pgx.Tx) error
	GetWIPProofByProver2(ctx context.Context, prover string, dbTx pgx.Tx) (*state.Proof2, error)
}
