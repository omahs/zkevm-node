package aggregator2

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/proverclient/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/keccak256"
	"google.golang.org/grpc"
)

// Aggregator represents an aggregator
type Aggregator2 struct {
	cfg Config

	State                stateInterface
	EthTxManager         ethTxManager
	Ethman               etherman
	ProverClients        []proverClientInterface
	ProfitabilityChecker aggregatorTxProfitabilityChecker
	TimeSendFinalProof   time.Time
	StateDBMutex         *sync.Mutex
}

// NewAggregator creates a new aggregator
func NewAggregator2(
	cfg Config,
	stateInterface stateInterface,
	ethTxManager ethTxManager,
	etherman etherman,
	grpcClientConns []*grpc.ClientConn,
) (Aggregator2, error) {
	var profitabilityChecker aggregatorTxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}

	proverClients := make([]proverClientInterface, 0, len(cfg.ProverURIs))
	//ctx := context.Background()

	a := Aggregator2{
		cfg: cfg,

		State:                stateInterface,
		EthTxManager:         ethTxManager,
		Ethman:               etherman,
		ProverClients:        proverClients,
		ProfitabilityChecker: profitabilityChecker,
		StateDBMutex:         &sync.Mutex{},
	}

	rand.Seed(time.Now().UnixNano()) //·

	for _, proverURI := range cfg.ProverURIs {
		proverClient := prover.NewClient(proverURI, cfg.IntervalFrequencyToGetProofGenerationState)
		proverClients = append(proverClients, proverClient)
		grpcClientConns = append(grpcClientConns, proverClient.Prover.Conn)
		log.Infof("Connected to prover %v", proverURI)

		/*		// Check if prover is already working in a proof generation
				proof, err := stateInterface.GetWIPProofByProver2(ctx, proverURI, nil)
				if err != nil && err != state.ErrNotFound {
					log.Errorf("Error while getting WIP proof for prover %v", proverURI)
					continue
				}

				if proof != nil {
					log.Infof("Resuming WIP proof generation for batchNumber %v in prover %v", proof.BatchNumber, *proof.Prover)
					go func() {
						a.resumeWIPProofGeneration(ctx, proof, proverClient)
					}()
				}*/
	}

	a.ProverClients = proverClients

	return a, nil
}

func (a *Aggregator2) resumeWIPProofGeneration(ctx context.Context, proof *state.Proof2, prover proverClientInterface) {
	/*	err := a.getAndStoreProof(ctx, proof, prover)
		if err != nil {
			log.Warnf("Could not resume WIP Proof Generation for prover %v and batchNumber %v", *proof.Prover, proof.BatchNumber)
		}*/
}

// Start starts the aggregator
func (a *Aggregator2) Start(ctx context.Context) {
	// define those vars here, bcs it can be used in case <-a.ctx.Done()
	tickerVerifyBatch := time.NewTicker(a.cfg.IntervalToConsolidateState.Duration)
	tickerSendVerifiedBatch := time.NewTicker(a.cfg.IntervalToSendFinalProof.Duration)
	defer tickerVerifyBatch.Stop()
	defer tickerSendVerifiedBatch.Stop()

	a.TimeSendFinalProof = time.Now().Add(a.cfg.IntervalToSendFinalProof.Duration)

	for i := 0; i < len(a.ProverClients); i++ {
		go func(prover proverClientInterface) {
			for {
				if prover.IsIdle(ctx) {
					proofGenerated, _ := a.tryAggregateProofs(ctx, prover, tickerVerifyBatch)
					if !proofGenerated {
						proofGenerated, _ = a.tryGenerateProofs(ctx, prover, tickerVerifyBatch)
					}
					if !proofGenerated {
						// if no proof was generated (aggregated or batch) wait some time waiting before retry
						time.Sleep(a.cfg.IntervalToConsolidateState.Duration)
					} // if proof was generated we retry inmediatly as probably we have more proofs to process
				} else {
					log.Warn("Prover %s is not idle", prover.GetURI())
					time.Sleep(a.cfg.IntervalToConsolidateState.Duration)
				}
			}
		}(a.ProverClients[i])
		time.Sleep(time.Second)
	}

	go func() {
		for {
			//			a.tryToSendVerifiedBatch(ctx, tickerSendVerifiedBatch)
		}
	}()
	// Wait until context is done
	<-ctx.Done()
}

/*func (a *Aggregator2) tryToSendVerifiedBatch(ctx context.Context, ticker *time.Ticker) {
	log.Debug("Checking if network is synced")
	for !a.isSynced(ctx) {
		log.Infof("Waiting for synchronizer to sync...")
		waitTick(ctx, ticker)
		continue
	}
	log.Debug("Checking if there is any consolidated batch to be verified")
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("Failed to get last consolidated batch, err: %v", err)
		waitTick(ctx, ticker)
		return
	} else if err == state.ErrNotFound {
		log.Debug("No consolidated batch found")
		waitTick(ctx, ticker)
		return
	}

	batchNumberToVerify := lastVerifiedBatch.BatchNumber + 1

	proof, err := a.State.GetGeneratedProofByBatchNumber(ctx, batchNumberToVerify, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("Failed to get last proof for batch %v, err: %v", batchNumberToVerify, err)
		waitTick(ctx, ticker)
		return
	}

	if proof != nil && proof.Proof != nil {
		log.Infof("Sending verified proof to the ethereum smart contract, batchNumber %d", batchNumberToVerify)
		err := a.EthTxManager.VerifyBatch(ctx, batchNumberToVerify, proof.Proof)
		if err != nil {
			log.Errorf("Error verifying batch %d, err: %w", batchNumberToVerify, err)
		} else {
			log.Infof("Proof for the batch was sent, batchNumber: %v", batchNumberToVerify)
			err := a.State.DeleteGeneratedProof2(ctx, batchNumberToVerify, batchNumberToVerify, nil)
			if err != nil {
				log.Warnf("Failed to delete generated proof for batchNumber %v, err: %v", batchNumberToVerify, err)
			}
		}
	} else {
		log.Debugf("No generated proof for batchNumber %v has been found", batchNumberToVerify)
		waitTick(ctx, ticker)
	}
}*/

func (a *Aggregator2) trySendFinalProof(ctx context.Context, proof *state.Proof2, ticker *time.Ticker) (bool, error) {
	if a.TimeSendFinalProof.Before(time.Now()) {
		log.Debug("Send final proof time reached")

		log.Debug("Checking if network is synced")
		for !a.isSynced(ctx) {
			log.Info("Waiting for synchronizer to sync...")
			waitTick(ctx, ticker)
			continue
		}

		lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
		if err != nil && err != state.ErrNotFound {
			log.Warnf("Failed to get last verified batch, err: %v", err)
			return false, err
		} else if err == state.ErrNotFound {
			log.Debug("Last verified batch not found")
			return false, err
		}

		batchNumberToVerify := lastVerifiedBatch.BatchNumber + 1

		if proof.BatchNumber == batchNumberToVerify {
			//· Calcular la final proof antes de enviarla
			log.Infof("Generating final proof for batches %d-%d", proof.BatchNumber, proof.BatchNumberFinal)

			log.Infof("Verfiying final proof with ethereum smart contract, batches %d-%d", proof.BatchNumber, proof.BatchNumberFinal)
			err := a.EthTxManager.VerifyBatch(ctx, proof.BatchNumberFinal, proof.Proof)
			if err != nil {
				log.Errorf("Error verifiying final proof for batches %d-%d, err: %w", proof.BatchNumber, proof.BatchNumberFinal, err)
				return false, err
			} else {
				log.Infof("Final proof for batches %d-%d verified", proof.BatchNumber, proof.BatchNumberFinal)
				a.TimeSendFinalProof = time.Now().Add(a.cfg.IntervalToSendFinalProof.Duration)
				return true, nil
			}
		} else {
			log.Infof("Proof batch number %d is not the following to last verfied batch number %d", proof.BatchNumber, batchNumberToVerify)
			return false, nil
		}
	} else {
		return false, nil
	}
}
func (a *Aggregator2) unlockProofsToAggregate(ctx context.Context, proof1 *state.Proof2, proof2 *state.Proof2, ticker *time.Ticker) error {
	// Release proofs from aggregating state in a single transaction
	dbTx, err := a.State.BeginStateTransaction(ctx)
	if err != nil {
		log.Warnf("Failed to begin transaction to release proof aggregation state, err: %v", err)
		return err
	}

	proof1.Aggregating = false
	err = a.State.UpdateGeneratedProof2(ctx, proof1, dbTx)
	if err == nil {
		proof2.Aggregating = false
		err = a.State.UpdateGeneratedProof2(ctx, proof2, dbTx)
	}

	if err != nil {
		log.Warnf("Failed to release proof aggregation state, err: %v", err)
		dbTx.Rollback(ctx)
		return err
	}

	dbTx.Commit(ctx)

	return nil
}

func (a *Aggregator2) getAndLockProofsToAggregate(ctx context.Context, prover proverClientInterface, ticker *time.Ticker) (*state.Proof2, *state.Proof2, error) {
	a.StateDBMutex.Lock()
	defer a.StateDBMutex.Unlock()

	proof1, proof2, err := a.State.GetProofsToAggregate(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Set proofs in aggregating state in a single transaction
	dbTx, err := a.State.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("Failed to begin transaction to set proof aggregation state, err: %v", err)
		return nil, nil, err
	}

	proof1.Aggregating = true
	err = a.State.UpdateGeneratedProof2(ctx, proof1, dbTx)
	if err == nil {
		proof2.Aggregating = true
		err = a.State.UpdateGeneratedProof2(ctx, proof2, dbTx)
	}

	if err != nil {
		log.Errorf("Failed to set proof aggregation state, err: %v", err)
		dbTx.Rollback(ctx)
		return nil, nil, err
	}

	dbTx.Commit(ctx)

	return proof1, proof2, nil
}

func (a *Aggregator2) tryAggregateProofs(ctx context.Context, prover proverClientInterface, ticker *time.Ticker) (bool, error) {
	log.Debugf("tryAggregateProols start %s", prover.GetURI())

	proof1, proof2, err0 := a.getAndLockProofsToAggregate(ctx, prover, ticker)
	if err0 != nil {
		//··waitTick(ctx, ticker)
		return false, err0
	}

	var err error

	defer func() {
		if err != nil {
			err2 := a.unlockProofsToAggregate(ctx, proof1, proof2, ticker)
			if err2 != nil {
				log.Errorf("Failed to release aggregated proofs, err: %v", err2)
			}
		}
	}()

	/*log.Infof("sending zki + batch to the prover, batchNumber: %d", batchToVerify.BatchNumber)
	inputProver, err := a.buildInputProver(ctx, batchToVerify)
	if err != nil {
		log.Warnf("failed to build input prover, err: %v", err)
		return false, err
	}*/

	log.Infof("Prover %s is going to be used to aggregate proofs: %d-%d and %d-%d", prover.GetURI(), proof1.BatchNumber, proof1.BatchNumberFinal, proof2.BatchNumber, proof2.BatchNumberFinal)

	/*log.Infof("sending a batch to the prover, OLDSTATEROOT: %s, NEWSTATEROOT: %s, BATCHNUM: %d",
	inputProver.PublicInputs.OldStateRoot, inputProver.PublicInputs.NewStateRoot, inputProver.PublicInputs.BatchNum)
	*/
	proverURI := prover.GetURI()
	//· Definir inputProver
	proof := &state.Proof2{BatchNumber: proof1.BatchNumber, BatchNumberFinal: proof2.BatchNumberFinal, Prover: &proverURI, InputProver: proof1.InputProver, Aggregating: false}

	var genProofID string     //· Eliminar
	genProofID = "1234567890" //· Eliminar
	/*genProofID, err := prover.GetGenProofID(ctx, inputProver)
	if err != nil {
		log.Warnf("failed to get gen proof id, err: %v", err)
		err2 := a.State.DeleteGeneratedProof(ctx, proof.BatchNumber, nil)
		if err2 != nil {
			log.Errorf("failed to delete proof generation mark after error, err: %v", err2)
		}
		//··waitTick(ctx, ticker)
		return false, err
	}*/

	proof.ProofID = &genProofID

	log.Infof("Proof ID for aggregated proof %d-%d: %v", proof.BatchNumber, proof.BatchNumberFinal, *proof.ProofID)

	/*err = a.getAndStoreProof(ctx, proof, prover)*/
	proof.Proof = proof1.Proof //· Quitar
	time.Sleep(time.Duration(rand.Intn(20)+10) * time.Second)

	proofSent, _ := a.trySendFinalProof(ctx, proof, ticker)

	dbTx, err := a.State.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("Failed to begin transaction to store proof aggregation result, err: %v", err)
		return false, err
	}

	// If the new aggregated proof has not been sent to L1 we store it
	if !proofSent {
		err = a.State.AddGeneratedProof2(ctx, proof, dbTx)
	}

	// Delete aggregated proofs
	if err == nil {
		err = a.State.DeleteGeneratedProof2(ctx, proof1.BatchNumber, proof1.BatchNumberFinal, nil)
	}
	if err == nil {
		err = a.State.DeleteGeneratedProof2(ctx, proof2.BatchNumber, proof2.BatchNumberFinal, nil)
	}

	if err != nil {
		dbTx.Rollback(ctx)
		log.Errorf("Failed to store proof aggregation result, err: %v", err)
		//··waitTick(ctx, ticker)
		return false, err
	}

	dbTx.Commit(ctx)

	log.Debug("tryAggregateProols end")

	return true, nil
}

func (a *Aggregator2) getAndLockBatchToProve(ctx context.Context, prover proverClientInterface, ticker *time.Ticker) (*state.Batch, *state.Proof2, error) {
	a.StateDBMutex.Lock()
	defer a.StateDBMutex.Unlock()

	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Get virtual batch pending to generate proof
	batchToVerify, err := a.State.GetVirtualBatchToProve(ctx, lastVerifiedBatch.BatchNumber, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("Found virtual batch %d pending to generate proof", batchToVerify.BatchNumber)

	log.Infof("Checking profitability to aggregate batch, batchNumber: %d", batchToVerify.BatchNumber)
	// pass matic collateral as zero here, bcs in smart contract fee for aggregator is not defined yet
	isProfitable, err := a.ProfitabilityChecker.IsProfitable(ctx, big.NewInt(0))
	if err != nil {
		log.Errorf("Failed to check aggregator profitability, err: %v", err)
		//··waitTick(ctx, ticker)
		return nil, nil, err
	}

	if !isProfitable {
		log.Infof("Batch %d is not profitable, matic collateral %d", batchToVerify.BatchNumber, big.NewInt(0))
		//··waitTick(ctx, ticker)
		return nil, nil, err
	}

	proverURI := prover.GetURI()
	proof := &state.Proof2{BatchNumber: batchToVerify.BatchNumber, BatchNumberFinal: batchToVerify.BatchNumber, Prover: &proverURI, Aggregating: false}

	// Avoid other thread to process the same batch
	err = a.State.AddGeneratedProof2(ctx, proof, nil)
	if err != nil {
		log.Errorf("Failed to add batch proof, err: %v", err)
		//··waitTick(ctx, ticker)
		return nil, nil, err
	}

	return batchToVerify, proof, nil
}

func (a *Aggregator2) tryGenerateProofs(ctx context.Context, prover proverClientInterface, ticker *time.Ticker) (bool, error) {
	batchToProve, proof, err0 := a.getAndLockBatchToProve(ctx, prover, ticker)
	if err0 != nil {
		//··waitTick(ctx, ticker)
		return false, err0
	}

	var err error

	defer func() {
		if err != nil {
			err2 := a.State.DeleteGeneratedProof2(ctx, proof.BatchNumber, proof.BatchNumberFinal, nil)
			if err2 != nil {
				log.Errorf("Failed to delete proof in progress, err: %v", err2)
			}
		}
	}()

	log.Infof("Prover %s is going to be used for batchNumber: %d", prover.GetURI(), batchToProve.BatchNumber)

	log.Infof("Sending zki + batch to the prover, batchNumber: %d", batchToProve.BatchNumber)
	proof.InputProver, err = a.buildInputProver(ctx, batchToProve)
	if err != nil {
		log.Errorf("Failed to build input prover, err: %v", err)
		//··waitTick(ctx, ticker)
		return false, err
	}

	log.Infof("Sending a batch to the prover, OLDSTATEROOT: %s, NEWSTATEROOT: %s, BATCHNUM: %d",
		proof.InputProver.PublicInputs.OldStateRoot, proof.InputProver.PublicInputs.NewStateRoot, proof.InputProver.PublicInputs.BatchNum)

	genProofID, err := prover.GetGenProofID(ctx, proof.InputProver)
	if err != nil {
		log.Errorf("Failed to get gen proof id, err: %v", err)
		//··waitTick(ctx, ticker)
		return false, err
	}

	proof.ProofID = &genProofID

	log.Infof("Proof ID for batchNumber %d: %v", proof.BatchNumber, *proof.ProofID)

	resGetProof, err := prover.GetResGetProof(ctx, *proof.ProofID, proof.BatchNumber)
	if err != nil {
		log.Errorf("Failed to get proof from prover, err: %v", err)
		//··waitTick(ctx, ticker)
		return false, err
	}

	proof.Proof = resGetProof

	a.compareInputHashes(proof.InputProver, proof.Proof)

	// Handle local exit root in the case of the mock prover
	if proof.Proof.Public.PublicInputs.NewLocalExitRoot == "0x17c04c3760510b48c6012742c540a81aba4bca2f78b9d14bfd2f123e2e53ea3e" {
		// This local exit root comes from the mock, use the one captured by the executor instead
		log.Warnf(
			"NewLocalExitRoot looks like a mock value, using value from executor instead: %v",
			proof.InputProver.PublicInputs.NewLocalExitRoot,
		)
		resGetProof.Public.PublicInputs.NewLocalExitRoot = proof.InputProver.PublicInputs.NewLocalExitRoot
	}

	proofSent, _ := a.trySendFinalProof(ctx, proof, ticker)

	if !proofSent {
		// Store proof
		err = a.State.UpdateGeneratedProof2(ctx, proof, nil)
		if err != nil {
			log.Errorf("Failed to store batch proof result, err: %v", err)
			//··waitTick(ctx, ticker)
			return false, err
		}
	}

	return true, nil
}

func (a *Aggregator2) isSynced(ctx context.Context) bool {
	lastVerifiedBatch, err := a.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Warnf("Failed to get last consolidated batch, err: %v", err)
		return false
	}
	if lastVerifiedBatch == nil {
		return false
	}
	lastVerifiedEthBatchNum, err := a.Ethman.GetLatestVerifiedBatchNum()
	if err != nil {
		log.Warnf("Failed to get last eth batch, err: %v", err)
		return false
	}
	if lastVerifiedBatch.BatchNumber < lastVerifiedEthBatchNum {
		log.Infof("Waiting for the state to be synced, lastVerifiedBatchNum: %d, lastVerifiedEthBatchNum: %d",
			lastVerifiedBatch.BatchNumber, lastVerifiedEthBatchNum)
		return false
	}
	return true
}

func (a *Aggregator2) buildInputProver(ctx context.Context, batchToVerify *state.Batch) (*pb.InputProver, error) {
	previousBatch, err := a.State.GetBatchByNumber(ctx, batchToVerify.BatchNumber-1, nil)
	if err != nil && err != state.ErrStateNotSynchronized {
		return nil, fmt.Errorf("Failed to get previous batch, err: %v", err)
	}

	blockTimestampByte := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(blockTimestampByte, uint64(batchToVerify.Timestamp.Unix()))
	batchHashData := common.BytesToHash(keccak256.Hash(
		batchToVerify.BatchL2Data,
		batchToVerify.GlobalExitRoot[:],
		blockTimestampByte,
		batchToVerify.Coinbase[:],
	))
	inputProver := &pb.InputProver{
		PublicInputs: &pb.PublicInputs{
			OldStateRoot:     previousBatch.StateRoot.String(),
			OldLocalExitRoot: previousBatch.LocalExitRoot.String(),
			NewStateRoot:     batchToVerify.StateRoot.String(),
			NewLocalExitRoot: batchToVerify.LocalExitRoot.String(),
			SequencerAddr:    batchToVerify.Coinbase.String(),
			BatchHashData:    batchHashData.String(),
			BatchNum:         uint32(batchToVerify.BatchNumber),
			EthTimestamp:     uint64(batchToVerify.Timestamp.Unix()),
			AggregatorAddr:   a.Ethman.GetPublicAddress().String(),
			ChainId:          a.cfg.ChainID,
		},
		GlobalExitRoot:    batchToVerify.GlobalExitRoot.String(),
		BatchL2Data:       hex.EncodeToString(batchToVerify.BatchL2Data),
		Db:                map[string]string{},
		ContractsBytecode: map[string]string{},
	}

	return inputProver, nil
}

func (a *Aggregator2) compareInputHashes(ip *pb.InputProver, resGetProof *pb.GetProofResponse) {
	// Calc inputHash
	batchNumberByte := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(batchNumberByte, ip.PublicInputs.BatchNum)
	blockTimestampByte := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(blockTimestampByte, ip.PublicInputs.EthTimestamp)
	hash := keccak256.Hash(
		[]byte(ip.PublicInputs.OldStateRoot)[:],
		[]byte(ip.PublicInputs.OldLocalExitRoot)[:],
		[]byte(ip.PublicInputs.NewStateRoot)[:],
		[]byte(ip.PublicInputs.NewLocalExitRoot)[:],
		[]byte(ip.PublicInputs.SequencerAddr)[:],
		[]byte(ip.PublicInputs.BatchHashData)[:],
		batchNumberByte[:],
		blockTimestampByte[:],
	)
	// Prime field. It is the prime number used as the order in our elliptic curve
	const fr = "21888242871839275222246405745257275088548364400416034343698204186575808495617"
	frB, _ := new(big.Int).SetString(fr, encoding.Base10)
	inputHashMod := new(big.Int).Mod(new(big.Int).SetBytes(hash), frB)
	internalInputHash := inputHashMod.Bytes()

	// InputHash must match
	internalInputHashS := fmt.Sprintf("0x%064s", hex.EncodeToString(internalInputHash))
	publicInputsExtended := resGetProof.GetPublic()
	if resGetProof.GetPublic().InputHash != internalInputHashS {
		log.Error("inputHash received from the prover (", publicInputsExtended.InputHash,
			") doesn't match with the internal value: ", internalInputHashS)
		log.Debug("internalBatchHashData: ", ip.PublicInputs.BatchHashData, " externalBatchHashData: ", publicInputsExtended.PublicInputs.BatchHashData)
		log.Debug("inputProver.PublicInputs.OldStateRoot: ", ip.PublicInputs.OldStateRoot)
		log.Debug("inputProver.PublicInputs.OldLocalExitRoot:", ip.PublicInputs.OldLocalExitRoot)
		log.Debug("inputProver.PublicInputs.NewStateRoot: ", ip.PublicInputs.NewStateRoot)
		log.Debug("inputProver.PublicInputs.NewLocalExitRoot: ", ip.PublicInputs.NewLocalExitRoot)
		log.Debug("inputProver.PublicInputs.SequencerAddr: ", ip.PublicInputs.SequencerAddr)
		log.Debug("inputProver.PublicInputs.BatchHashData: ", ip.PublicInputs.BatchHashData)
		log.Debug("inputProver.PublicInputs.BatchNum: ", ip.PublicInputs.BatchNum)
		log.Debug("inputProver.PublicInputs.EthTimestamp: ", ip.PublicInputs.EthTimestamp)
	}
}

func waitTick(ctx context.Context, ticker *time.Ticker) {
	select {
	case <-ticker.C:
		// nothing
	case <-ctx.Done():
		return
	}
}
