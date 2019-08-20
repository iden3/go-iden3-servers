package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-servers/cmd/genericserver"
)

func handleInfo(c *gin.Context) {

	c.JSON(200, gin.H{
		"status": "everything up",
	})
}

type CreateIdentityReq struct {
	ClaimAuthKOp       *merkletree.Entry   `json:"claimAuthKOp" binding:"required"`
	ExtraGenesisClaims []*merkletree.Entry `json:"extraGenesisClaims"`
}

func handlePostIdentity(c *gin.Context) {
	var createIdentityReq CreateIdentityReq
	err := c.BindJSON(&createIdentityReq)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	id, proofKOp, err := ia.NewIdentity(createIdentityReq.ClaimAuthKOp, createIdentityReq.ExtraGenesisClaims)
	if err != nil {
		genericserver.Fail(c, "error on NewIdentity, probably Identity already exists", err)
		return
	}
	c.JSON(200, gin.H{
		"id":         id.String(),
		"proofOkKey": proofKOp,
	})
}

type ClaimReq struct {
	Claim *merkletree.Entry `json:"claim" binding:"required"`
}

func handlePostClaim(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	var claimReq ClaimReq
	err = c.BindJSON(&claimReq)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	err = ia.AddClaim(&id, claimReq.Claim)
	if err != nil {
		genericserver.Fail(c, "AddClaim error", err)
		return
	}

	c.JSON(200, gin.H{})
}

type ClaimsReq struct {
	Claims []*merkletree.Entry `json:"claims" binding:"required"`
}

func handlePostClaims(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	var claimsReq ClaimsReq
	err = c.BindJSON(&claimsReq)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	err = ia.AddClaims(&id, claimsReq.Claims)
	if err != nil {
		genericserver.Fail(c, "AddClaims error", err)
		return
	}

	c.JSON(200, gin.H{})
}

func handleGetAllReceivedClaims(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	idStorages, err := ia.LoadIdStorages(&id)
	if err != nil {
		genericserver.Fail(c, "error loading IdStorages", err)
		return
	}
	claims, err := ia.GetAllReceivedClaims(&id, idStorages)
	if err != nil {
		genericserver.Fail(c, "GetAllReceivedClaims error", err)
		return
	}
	c.JSON(200, gin.H{
		"receivedClaims": core.ClaimObjArrayToHexArray(claims),
	})
}

func handleGetAllEmittedClaims(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	idStorages, err := ia.LoadIdStorages(&id)
	if err != nil {
		genericserver.Fail(c, "error loading IdStorages", err)
		return
	}
	claims, err := ia.GetAllEmittedClaims(&id, idStorages)
	if err != nil {
		genericserver.Fail(c, "GetAllEmittedClaims error", err)
		return
	}
	c.JSON(200, gin.H{
		"emittedClaims": core.ClaimObjArrayToHexArray(claims),
	})
}

func handleGetAllClaims(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	emittedClaims, receivedClaims, err := ia.GetAllClaims(&id)
	if err != nil {
		genericserver.Fail(c, "GetAllClaims error", err)
		return
	}
	c.JSON(200, gin.H{
		"emittedClaims":  core.ClaimObjArrayToHexArray(emittedClaims),
		"receivedClaims": core.ClaimObjArrayToHexArray(receivedClaims),
	})
}

func handleGetFullMT(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	mt, err := ia.GetFullMT(&id)
	if err != nil {
		genericserver.Fail(c, "GetFullMT error", err)
		return
	}
	c.JSON(200, gin.H{
		"mt": mt,
	})
}
