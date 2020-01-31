package endpoint

import (

	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-crypto/babyjub"
	// "github.com/iden3/go-iden3-core/utils"
)

type handleIdGenesis struct {
	KOp         *babyjub.PublicKey `json:"operationalPk" binding:"required"`
	KDis        common.Address     `json:"kdisable" binding:"required"`
	KReen       common.Address     `json:"kreenable" binding:"required"`
	KUpdateRoot common.Address     `json:"kupdateRoot" binding:"required"`
}

// handlePostIdRes is the response of a creation of a new user tree in the relay.
type handlePostIdRes struct {
	Id         core.ID           `json:"id"`
	ProofClaim *proof.ProofClaim `json:"proofClaim"`
}
