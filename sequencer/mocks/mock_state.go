// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v4"

	state "github.com/0xPolygonHermez/zkevm-node/state"

	time "time"

	types "github.com/ethereum/go-ethereum/core/types"
)

// StateMock is an autogenerated mock type for the stateInterface type
type StateMock struct {
	mock.Mock
}

// BeginStateTransaction provides a mock function with given fields: ctx
func (_m *StateMock) BeginStateTransaction(ctx context.Context) (pgx.Tx, error) {
	ret := _m.Called(ctx)

	var r0 pgx.Tx
	if rf, ok := ret.Get(0).(func(context.Context) pgx.Tx); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CloseBatch provides a mock function with given fields: ctx, receipt, dbTx
func (_m *StateMock) CloseBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error {
	ret := _m.Called(ctx, receipt, dbTx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, state.ProcessingReceipt, pgx.Tx) error); ok {
		r0 = rf(ctx, receipt, dbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBatchByNumber provides a mock function with given fields: ctx, batchNumber, dbTx
func (_m *StateMock) GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error) {
	ret := _m.Called(ctx, batchNumber, dbTx)

	var r0 *state.Batch
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) *state.Batch); ok {
		r0 = rf(ctx, batchNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*state.Batch)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, batchNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockNumAndMainnetExitRootByGER provides a mock function with given fields: ctx, ger, dbTx
func (_m *StateMock) GetBlockNumAndMainnetExitRootByGER(ctx context.Context, ger common.Hash, dbTx pgx.Tx) (uint64, common.Hash, error) {
	ret := _m.Called(ctx, ger, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, pgx.Tx) uint64); ok {
		r0 = rf(ctx, ger, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 common.Hash
	if rf, ok := ret.Get(1).(func(context.Context, common.Hash, pgx.Tx) common.Hash); ok {
		r1 = rf(ctx, ger, dbTx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(common.Hash)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, common.Hash, pgx.Tx) error); ok {
		r2 = rf(ctx, ger, dbTx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetLastBatch provides a mock function with given fields: ctx, dbTx
func (_m *StateMock) GetLastBatch(ctx context.Context, dbTx pgx.Tx) (*state.Batch, error) {
	ret := _m.Called(ctx, dbTx)

	var r0 *state.Batch
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) *state.Batch); ok {
		r0 = rf(ctx, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*state.Batch)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastBatchNumber provides a mock function with given fields: ctx, dbTx
func (_m *StateMock) GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	ret := _m.Called(ctx, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) uint64); ok {
		r0 = rf(ctx, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastL2BlockNumber provides a mock function with given fields: ctx, dbTx
func (_m *StateMock) GetLastL2BlockNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	ret := _m.Called(ctx, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) uint64); ok {
		r0 = rf(ctx, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastVirtualBatchNum provides a mock function with given fields: ctx, dbTx
func (_m *StateMock) GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	ret := _m.Called(ctx, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) uint64); ok {
		r0 = rf(ctx, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestGlobalExitRoot provides a mock function with given fields: ctx, maxBlockNumber, dbTx
func (_m *StateMock) GetLatestGlobalExitRoot(ctx context.Context, maxBlockNumber uint64, dbTx pgx.Tx) (state.GlobalExitRoot, time.Time, error) {
	ret := _m.Called(ctx, maxBlockNumber, dbTx)

	var r0 state.GlobalExitRoot
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) state.GlobalExitRoot); ok {
		r0 = rf(ctx, maxBlockNumber, dbTx)
	} else {
		r0 = ret.Get(0).(state.GlobalExitRoot)
	}

	var r1 time.Time
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) time.Time); ok {
		r1 = rf(ctx, maxBlockNumber, dbTx)
	} else {
		r1 = ret.Get(1).(time.Time)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, uint64, pgx.Tx) error); ok {
		r2 = rf(ctx, maxBlockNumber, dbTx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetNonce provides a mock function with given fields: ctx, address, blockNumber, dbTx
func (_m *StateMock) GetNonce(ctx context.Context, address common.Address, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	ret := _m.Called(ctx, address, blockNumber, dbTx)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, uint64, pgx.Tx) uint64); ok {
		r0 = rf(ctx, address, blockNumber, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, address, blockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStateRootByBatchNumber provides a mock function with given fields: ctx, batchNum, dbTx
func (_m *StateMock) GetStateRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error) {
	ret := _m.Called(ctx, batchNum, dbTx)

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) common.Hash); ok {
		r0 = rf(ctx, batchNum, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, batchNum, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTimeForLatestBatchVirtualization provides a mock function with given fields: ctx, dbTx
func (_m *StateMock) GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
	ret := _m.Called(ctx, dbTx)

	var r0 time.Time
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Tx) time.Time); ok {
		r0 = rf(ctx, dbTx)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pgx.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionsByBatchNumber provides a mock function with given fields: ctx, batchNumber, dbTx
func (_m *StateMock) GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]types.Transaction, error) {
	ret := _m.Called(ctx, batchNumber, dbTx)

	var r0 []types.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) []types.Transaction); ok {
		r0 = rf(ctx, batchNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, batchNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxsOlderThanNL1Blocks provides a mock function with given fields: ctx, nL1Blocks, dbTx
func (_m *StateMock) GetTxsOlderThanNL1Blocks(ctx context.Context, nL1Blocks uint64, dbTx pgx.Tx) ([]common.Hash, error) {
	ret := _m.Called(ctx, nL1Blocks, dbTx)

	var r0 []common.Hash
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) []common.Hash); ok {
		r0 = rf(ctx, nL1Blocks, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Hash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, nL1Blocks, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsBatchClosed provides a mock function with given fields: ctx, batchNum, dbTx
func (_m *StateMock) IsBatchClosed(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error) {
	ret := _m.Called(ctx, batchNum, dbTx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) bool); ok {
		r0 = rf(ctx, batchNum, dbTx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, batchNum, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsBatchVirtualized provides a mock function with given fields: ctx, batchNumber, dbTx
func (_m *StateMock) IsBatchVirtualized(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (bool, error) {
	ret := _m.Called(ctx, batchNumber, dbTx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, uint64, pgx.Tx) bool); ok {
		r0 = rf(ctx, batchNumber, dbTx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, pgx.Tx) error); ok {
		r1 = rf(ctx, batchNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OpenBatch provides a mock function with given fields: ctx, processingContext, dbTx
func (_m *StateMock) OpenBatch(ctx context.Context, processingContext state.ProcessingContext, dbTx pgx.Tx) error {
	ret := _m.Called(ctx, processingContext, dbTx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, state.ProcessingContext, pgx.Tx) error); ok {
		r0 = rf(ctx, processingContext, dbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessSequencerBatch provides a mock function with given fields: ctx, batchNumber, txs, dbTx
func (_m *StateMock) ProcessSequencerBatch(ctx context.Context, batchNumber uint64, txs []types.Transaction, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	ret := _m.Called(ctx, batchNumber, txs, dbTx)

	var r0 *state.ProcessBatchResponse
	if rf, ok := ret.Get(0).(func(context.Context, uint64, []types.Transaction, pgx.Tx) *state.ProcessBatchResponse); ok {
		r0 = rf(ctx, batchNumber, txs, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*state.ProcessBatchResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64, []types.Transaction, pgx.Tx) error); ok {
		r1 = rf(ctx, batchNumber, txs, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreTransactions provides a mock function with given fields: ctx, batchNum, processedTxs, dbTx
func (_m *StateMock) StoreTransactions(ctx context.Context, batchNum uint64, processedTxs []*state.ProcessTransactionResponse, dbTx pgx.Tx) error {
	ret := _m.Called(ctx, batchNum, processedTxs, dbTx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, []*state.ProcessTransactionResponse, pgx.Tx) error); ok {
		r0 = rf(ctx, batchNum, processedTxs, dbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateGERInOpenBatch provides a mock function with given fields: ctx, ger, dbTx
func (_m *StateMock) UpdateGERInOpenBatch(ctx context.Context, ger common.Hash, dbTx pgx.Tx) error {
	ret := _m.Called(ctx, ger, dbTx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash, pgx.Tx) error); ok {
		r0 = rf(ctx, ger, dbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStateMock interface {
	mock.TestingT
	Cleanup(func())
}

// NewStateMock creates a new instance of StateMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStateMock(t mockConstructorTestingTNewStateMock) *StateMock {
	mock := &StateMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
