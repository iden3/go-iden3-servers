package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/claimsrv"
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

func handleGetRoot(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}
	c.JSON(200, gin.H{
		"root": agent.GetCurrentRoot(),
	})
}

func handlePostRoot(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}
	var setRootReq claimsrv.SetRoot0Req
	err := c.ShouldBindJSON(&setRootReq)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}
	if err := agent.RootUpdate(setRootReq); err != nil {
		genericserver.Fail(c, "RootUpdate error", err)
		return
	}
	c.JSON(200, gin.H{})
}

func handleGetClaimsReceived(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}

	claims, err := agent.ClaimsReceived()
	if err != nil {
		genericserver.Fail(c, "ClaimsReceived error", err)
		return
	}
	c.JSON(200, gin.H{
		"receivedClaims": claims,
	})
}

func handleGetClaimsEmitted(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}

	claims, err := agent.ClaimsEmitted()
	if err != nil {
		genericserver.Fail(c, "ClaimsEmitted error", err)
		return
	}
	c.JSON(200, gin.H{
		"emittedClaims": claims,
	})
}

func handleGetClaimsGenesis(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}

	claims, err := agent.ClaimsGenesis()
	if err != nil {
		genericserver.Fail(c, "ClaimsGenesis error", err)
		return
	}
	c.JSON(200, gin.H{
		"genesisClaims": claims,
	})
}

func handleGetMT(c *gin.Context) {
	agent := loadAgent(c)
	if agent == nil {
		return
	}

	mt, err := agent.ExportMT()
	if err != nil {
		genericserver.Fail(c, "ExportMT error", err)
		return
	}
	c.JSON(200, gin.H{
		"mt": mt,
	})
}
