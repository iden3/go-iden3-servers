package endpoint

import (
	"bytes"
	"errors"

	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3-core/common"
	idenmanagermsgs "github.com/iden3/go-iden3-core/components/idenmanager/messages"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/crypto"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-servers/handlers"
	"github.com/iden3/go-iden3-servers/loaders"
)

// SetRootMsg contains the data to set the SetRootClaim with its signature in Hex
type SetRootMsg struct {
	Root      *merkletree.Hash        `binding:"required"`
	Id        *core.ID                `binding:"required"`
	KSignPk   *crypto.PublicKey       `binding:"required"`
	Timestamp int64                   `binding:"required"`
	Signature *crypto.SignatureEthMsg `binding:"required"`
}

// handleCommitNewIdRoot handles a request to set the root key of a user tree
// though a set root claim.
func handleCommitNewIdRoot(c *gin.Context, iden *loaders.Identity) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		handlers.Fail(c, "error on parse id", err)
		return
	}

	var setRootMsg SetRootMsg
	err = c.BindJSON(&setRootMsg)
	if err != nil {
		handlers.Fail(c, "json parsing error", err)
		return
	}

	// make sure that the given id from the post url matches with the id from the post data
	if !bytes.Equal(id.Bytes(), setRootMsg.Id.Bytes()) {
		handlers.Fail(c, "error on CommitNewIdRoot, id not match",
			errors.New("CommitNewIdRoot id not match"))
		return
	}

	// add the root through iden.Manager.Claimservice
	setRootClaim, err := iden.Manager.CommitNewIdRoot(id,
		&setRootMsg.KSignPk.PublicKey, *setRootMsg.Root, setRootMsg.Timestamp, setRootMsg.Signature)
	if err != nil {
		handlers.Fail(c, "error on CommitNewIdRoot", err)
		return
	}

	// Don't return the claim with proof because it uses a Relay root that
	// is not on the blockchain and may never be there.
	// proofRelayClaim, err := iden.Manager.GetClaimProofByHi(setRootClaim.Entry().HIndex())
	// if err != nil {
	// 	handlers.Fail(c, "error on GetClaimByHi", err)
	// 	return
	// }
	c.JSON(200, gin.H{
		"setRootClaim": setRootClaim,
	})
}

// handleUpdateSetRootClaim handles a request to add a new set root claim to the relay.
func handleUpdateSetRootClaim(c *gin.Context, iden *loaders.Identity) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		handlers.Fail(c, "error on parse id", err)
		return
	}

	var setRootReq idenmanagermsgs.SetRoot0Req
	err = c.BindJSON(&setRootReq)
	if err != nil {
		handlers.Fail(c, "json parsing error", err)
		return
	}

	// make sure that the given id from the post url matches with the id from the post data
	if !id.Equal(setRootReq.ProofClaimAuthKOp.Id) {
		handlers.Fail(c, "error on handleUpdateSetRootClaim, id not match with genesis proof",
			errors.New("handleUpdateSetRootClaim id not match with genesis proof"))
		return
	}

	// add the root through iden.Manager
	setRootClaim, err := iden.Manager.UpdateSetRootClaim(&id, setRootReq)
	if err != nil {
		handlers.Fail(c, "error on UpdateSetRootClaim", err)
		return
	}

	c.JSON(200, gin.H{
		"setRootClaim": setRootClaim.Entry(),
	})
}

// handleGetIdRoot handles a request to query the root key of a user tree.
func handleGetIdRoot(c *gin.Context, iden *loaders.Identity) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		handlers.Fail(c, "error on parse id", err)
		return
	}
	idRoot, idRootProof, err := iden.Manager.GetIdRoot(&id)
	if err != nil {
		handlers.Fail(c, "error on GetIdRoot", err)
		return
	}
	c.JSON(200, gin.H{
		"root":        iden.Manager.MT().RootKey().Hex(), // relay root
		"idRoot":      idRoot.Hex(),                      // user id root
		"proofIdRoot": common3.HexEncode(idRootProof),    // user id root proof in the relay merkletree
	})
	return
}

// handleGetSetRootClaim handles a request to query the last SetRootKey claim
// of an id with a proof to the root in the blockchain.
func handleGetSetRootClaim(c *gin.Context, iden *loaders.Identity) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		handlers.Fail(c, "error on parse id", err)
		return
	}

	proofClaim, err := iden.Manager.GetSetRootClaim(&id)
	if err != nil {
		handlers.Fail(c, "error on GetClaimProofByHiBlockchain", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofClaim,
	})
}

// handleGetClaimProofByHiBlockchain handles the request to query the claim proof of a
// relay claim (by hIndex) to a root published in the blockchain.
func handleGetClaimProofByHiBlockchain(c *gin.Context, iden *loaders.Identity) {
	hihex := c.Param("hi")
	hiBytes, err := common3.HexDecode(hihex)
	if err != nil {
		handlers.Fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	proofClaim, err := iden.Manager.GetClaimProofByHiBlockchain(hi)
	if err != nil {
		handlers.Fail(c, "error on GetClaimProofByHiBlockchain", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofClaim,
	})
	return
}