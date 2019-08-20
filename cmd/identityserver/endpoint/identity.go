package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/identityagentsrv"
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
	err := c.ShouldBindJSON(&createIdentityReq)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	id, proofKOp, err := ia.CreateIdentity(createIdentityReq.ClaimAuthKOp, createIdentityReq.ExtraGenesisClaims)
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

func loadAgent(c *gin.Context) *identityagentsrv.Agent {
	id, err := core.IDFromString(c.Param("id"))
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return nil
	}
	agent, err := ia.NewAgent(&id)
	if err != nil {
		genericserver.Fail(c, "error loading Agent", err)
		return nil
	}
	return agent
}

func handlePostClaim(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}

	var claimReq ClaimReq
	err := c.ShouldBindJSON(&claimReq)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	err = agent.AddClaim(claimReq.Claim)
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
	agent := loadAgent(c)
	if agent == nil {
		return
	}

	var claimsReq ClaimsReq
	err := c.ShouldBindJSON(&claimsReq)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	err = agent.AddClaims(claimsReq.Claims)
	if err != nil {
		genericserver.Fail(c, "AddClaims error", err)
		return
	}

	c.JSON(200, gin.H{})
}

func handleGetAllReceivedClaims(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}

	claims, err := agent.GetAllReceivedClaims()
	if err != nil {
		genericserver.Fail(c, "GetAllReceivedClaims error", err)
		return
	}
	c.JSON(200, gin.H{
		"receivedClaims": claims,
	})
}

func handleGetAllEmittedClaims(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}

	claims, err := agent.GetAllEmittedClaims()
	if err != nil {
		genericserver.Fail(c, "GetAllEmittedClaims error", err)
		return
	}
	c.JSON(200, gin.H{
		"emittedClaims": claims,
	})
}

func handleGetFullMT(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}

	mt, err := agent.GetFullMT()
	if err != nil {
		genericserver.Fail(c, "GetFullMT error", err)
		return
	}
	c.JSON(200, gin.H{
		"mt": mt,
	})
}
