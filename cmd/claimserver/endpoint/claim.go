package endpoint

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/components/idenmanager"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/crypto"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-servers/cmd/genericserver"
)

// IdData struct representing user data that claim server will manage afterwards.
type IdData struct {
	Id          core.ID `json:"id"`
	NotifSrvUrl string  `json:"notifSrvUrl"`
}

type IdDataB64 IdData

// UnmarshalText retrieve data from an array of bytes.
func (d *IdDataB64) UnmarshalText(text []byte) error {
	idDataJSON, err := base64.URLEncoding.WithPadding(base64.NoPadding).
		DecodeString(string(text))
	if err != nil {
		return err
	}
	var idData IdData
	if err := json.Unmarshal(idDataJSON, &idData); err != nil {
		return err
	}
	*d = IdDataB64(idData)
	return nil
}

// claimData struct representing data needed in order to be accepted by handlePostClaim function.
type claimData struct {
	IdData IdDataB64 `json:"idData" binding:"required"`
	Cert   string    `json:"data" binding:"required"`
}

// handlePostClaim handles the request to add a claim to a user tree.
func handlePostClaim(c *gin.Context) {
	var m claimData
	if err := c.BindJSON(&m); err != nil {
		genericserver.Fail(c, "cannot parse json body", err)
		return
	}

	hash := claims.ClearMostSigByte(crypto.HashBytes([]byte(m.Cert)))
	// Pending to update according new data received by the server
	auxData := claims.ClearMostSigByte(crypto.HashBytes([]byte(m.Cert)))
	objectType := claims.ObjectTypeCertificate
	indexObject := uint16(0)
	claim, err := claims.NewClaimLinkObjectIdentity(objectType, indexObject,
		m.IdData.Id, hash, auxData)
	if err != nil {
		genericserver.Fail(c, "error on NewClaimLinkObjectIdentity", err)
		return
	}

	// If necessary store the claim with a version higher than an existing
	// claim to invalidate the later.
	version, err := idenmanager.GetNextVersion(genericserver.Claimservice.MT(), claim.Entry().HIndex())
	if err != nil {
		genericserver.Fail(c, "error on GetNextVersion", err)
		return
	}
	claim.Version = version

	// Add claim to claim server merke tree.
	err = genericserver.Claimservice.AddClaim(claim)
	if err != nil {
		genericserver.Fail(c, "error on AddLinkObjectClaim", err)
		return
	}

	// TODO
	// return claim with proofs.
	// proofClaim, err := genericserver.Claimservice.GetClaimProofByHiBlockchain(claim.Entry().HIndex())
	// if err != nil {
	// 	genericserver.Fail(c, "error on GetClaimProofByHi", err)
	// 	return
	// }

	c.JSON(200, gin.H{
		"status": "ok",
	})
	return
}

// handleGetClaimProofByHi handles the request to query the claim proof of a
// server claim (by hIndex).
func handleGetClaimProofByHi(c *gin.Context) {
	hihex := c.Param("hi")
	hiBytes, err := common3.HexDecode(hihex)
	if err != nil {
		genericserver.Fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	proofOfClaim, err := genericserver.Claimservice.GetClaimProofByHiBlockchain(hi)
	if err != nil {
		genericserver.Fail(c, "error on GetClaimProofByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofOfClaim,
	})
	return
}
