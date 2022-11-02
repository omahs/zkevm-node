package state

import "github.com/0xPolygonHermez/zkevm-node/proverclient/pb"

// Proof struct
type Proof2 struct {
	BatchNumber      uint64
	BatchNumberFinal uint64
	Proof            *pb.GetProofResponse
	InputProver      *pb.InputProver
	ProofID          *string
	Prover           *string
	Aggregating      bool
}
