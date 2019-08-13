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

type NewIdentityMsg struct {
	ClaimAuthKOpHex       string   `json:"claimAuthKOp"`
	ExtraGenesisClaimsHex []string `json:"extraGenesisClaims"`
}

func handlePostIdentity(c *gin.Context) {
	var newIdentityMsg NewIdentityMsg
	err := c.BindJSON(&newIdentityMsg)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	claimKOp, err := core.HexToClaim(newIdentityMsg.ClaimAuthKOpHex)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	var extraGenesisClaims []merkletree.Claim
	for _, claimHex := range newIdentityMsg.ExtraGenesisClaimsHex {
		claim, err := core.HexToClaim(claimHex)
		if err != nil {
			genericserver.Fail(c, "json parsing error", err)
			return
		}
		extraGenesisClaims = append(extraGenesisClaims, claim)
	}

	id, proofKOp, err := ia.NewIdentity(claimKOp, extraGenesisClaims)
	if err != nil {
		genericserver.Fail(c, "error on NewIdentity, probably Identity already exists", err)
		return
	}
	c.JSON(200, gin.H{
		"id":         id.String(),
		"proofOkKey": proofKOp,
	})
}
