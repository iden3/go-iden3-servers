package endpoint

import (
	"bytes"
	"errors"

	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/services/claimsrv"
	"github.com/iden3/go-iden3-servers/cmd/genericserver"

	"github.com/iden3/go-iden3-core/merkletree"
)

// handleCommitNewIdRoot handles a request to set the root key of a user tree
// though a set root claim.
func handleCommitNewIdRoot(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	var setRootMsg claimsrv.SetRootMsg
	err = c.BindJSON(&setRootMsg)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	// make sure that the given id from the post url matches with the id from the post data
	if !bytes.Equal(id.Bytes(), setRootMsg.Id.Bytes()) {
		genericserver.Fail(c, "error on CommitNewIdRoot, id not match",
			errors.New("CommitNewIdRoot id not match"))
		return
	}

	// add the root through genericserver.Claimservice
	setRootClaim, err := genericserver.Claimservice.CommitNewIdRoot(id,
		&setRootMsg.KSignPk.PublicKey, *setRootMsg.Root, setRootMsg.Timestamp, setRootMsg.Signature)
	if err != nil {
		genericserver.Fail(c, "error on CommitNewIdRoot", err)
		return
	}

	// Don't return the claim with proof because it uses a Relay root that
	// is not on the blockchain and may never be there.
	// proofRelayClaim, err := genericserver.Claimservice.GetClaimProofByHi(setRootClaim.Entry().HIndex())
	// if err != nil {
	// 	genericserver.Fail(c, "error on GetClaimByHi", err)
	// 	return
	// }
	c.JSON(200, gin.H{
		"setRootClaim": setRootClaim,
	})
}

// handleUpdateSetRootClaim handles a request to add a new set root claim to the relay.
func handleUpdateSetRootClaim(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	var setRootReq claimsrv.SetRoot0Req
	err = c.BindJSON(&setRootReq)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	// make sure that the given id from the post url matches with the id from the post data
	if !id.Equal(setRootReq.ProofClaimAuthKOp.Id) {
		genericserver.Fail(c, "error on handleUpdateSetRootClaim, id not match with genesis proof",
			errors.New("handleUpdateSetRootClaim id not match with genesis proof"))
		return
	}

	// add the root through genericserver.Claimservice
	setRootClaim, err := genericserver.Claimservice.UpdateSetRootClaim(&id, setRootReq)
	if err != nil {
		genericserver.Fail(c, "error on UpdateSetRootClaim", err)
		return
	}

	c.JSON(200, gin.H{
		"setRootClaim": setRootClaim.Entry(),
	})
}

// DEPRECATED
// handlePostClaim handles the request to add a claim to a user tree.
// func handlePostClaim(c *gin.Context) {
// 	idHex := c.Param("id")
// 	id, err := core.IDFromString(idHex)
// 	if err != nil {
// 		genericserver.Fail(c, "error on parse id", err)
// 		return
// 	}
// 	var bytesSignedMsg claimsrv.BytesSignedMsg
// 	err = c.BindJSON(&bytesSignedMsg)
// 	if err != nil {
// 		genericserver.Fail(c, "json parsing error", err)
// 		return
// 	}
//
// 	bytesValue, err := common3.HexDecode(bytesSignedMsg.ValueHex)
// 	if err != nil {
// 		genericserver.Fail(c, "error on parsing bytesSignedMsg.HexValue to bytes", err)
// 		return
// 	}
//
// 	// bytesValue to Element data
// 	var dataBytes [128]byte
// 	copy(dataBytes[:], bytesValue)
// 	data := merkletree.NewDataFromBytes(dataBytes)
// 	entry := merkletree.Entry{
// 		Data: *data,
// 	}
//
// 	claimValueMsg := claimsrv.ClaimValueMsg{
// 		ClaimValue: entry,
// 		Signature:  bytesSignedMsg.Signature,
// 		KSignPk:    bytesSignedMsg.KSignPk,
// 	}
// 	err = genericserver.Claimservice.AddUserIdClaim(&id, claimValueMsg)
// 	if err != nil {
// 		genericserver.Fail(c, "error on AddUserIdClaim", err)
// 		return
// 	}
// 	// return claim with proofs
// 	proofClaim, err := genericserver.Claimservice.GetClaimProofUserByHi(id, entry.HIndex())
// 	if err != nil {
// 		genericserver.Fail(c, "error on GetClaimByHi", err)
// 		return
// 	}
//
// 	c.JSON(200, gin.H{
// 		"proofClaim": proofClaim,
// 	})
// 	return
// }

// handleGetIdRoot handles a request to query the root key of a user tree.
func handleGetIdRoot(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}
	idRoot, idRootProof, err := genericserver.Claimservice.GetIdRoot(&id)
	if err != nil {
		genericserver.Fail(c, "error on GetIdRoot", err)
		return
	}
	c.JSON(200, gin.H{
		"root":        genericserver.Claimservice.MT().RootKey().Hex(), // relay root
		"idRoot":      idRoot.Hex(),                                    // user id root
		"proofIdRoot": common3.HexEncode(idRootProof),                  // user id root proof in the relay merkletree
	})
	return
}

// handleGetSetRootClaim handles a request to query the last SetRootKey claim
// of an id with a proof to the root in the blockchain.
func handleGetSetRootClaim(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	proofClaim, err := genericserver.Claimservice.GetSetRootClaim(&id)
	if err != nil {
		genericserver.Fail(c, "error on GetClaimProofByHiBlockchain", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofClaim,
	})
}

// handleGetClaimProofUserByHi handles the request to query the claim proof of
// a user claim (by hIndex).
// DEPRECATED
// func handleGetClaimProofUserByHi(c *gin.Context) {
// 	idHex := c.Param("id")
// 	hiHex := c.Param("hi")
// 	hiBytes, err := common3.HexDecode(hiHex)
// 	if err != nil {
// 		genericserver.Fail(c, "error on HexDecode of Hi", err)
// 		return
// 	}
// 	hi := &merkletree.Hash{}
// 	copy(hi[:], hiBytes)
// 	id, err := core.IDFromString(idHex)
// 	if err != nil {
// 		genericserver.Fail(c, "error on parse id", err)
// 		return
// 	}
// 	proofClaim, err := genericserver.Claimservice.GetClaimProofUserByHi(id, hi)
// 	if err != nil {
// 		genericserver.Fail(c, "error on GetClaimByHi", err)
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"proofClaim": proofClaim,
// 	})
// 	return
// }

// TODO: Deprecate handleGetClaimProofByHi
// handleGetClaimProofByHi handles the request to query the claim proof of a
// relay claim (by hIndex).
// func handleGetClaimProofByHi(c *gin.Context) {
// 	hihex := c.Param("hi")
// 	hiBytes, err := common3.HexDecode(hihex)
// 	if err != nil {
// 		genericserver.Fail(c, "error on HexDecode of Hi", err)
// 		return
// 	}
// 	hi := &merkletree.Hash{}
// 	copy(hi[:], hiBytes)
// 	proofClaim, err := genericserver.Claimservice.GetClaimProofByHi(hi)
// 	if err != nil {
// 		genericserver.Fail(c, "error on GetClaimProofByHi", err)
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"proofClaim": proofClaim,
// 	})
// 	return
// }

// handleGetClaimProofByHiBlockchain handles the request to query the claim proof of a
// relay claim (by hIndex) to a root published in the blockchain.
func handleGetClaimProofByHiBlockchain(c *gin.Context) {
	hihex := c.Param("hi")
	hiBytes, err := common3.HexDecode(hihex)
	if err != nil {
		genericserver.Fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	proofClaim, err := genericserver.Claimservice.GetClaimProofByHiBlockchain(hi)
	if err != nil {
		genericserver.Fail(c, "error on GetClaimProofByHiBlockchain", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofClaim,
	})
	return
}
