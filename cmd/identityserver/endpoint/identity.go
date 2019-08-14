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

type ClaimMsg struct {
	ClaimHex string `json:"claim"`
}

func handlePostClaim(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	var claimMsg ClaimMsg
	err = c.BindJSON(&claimMsg)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	claim, err := core.HexToClaim(claimMsg.ClaimHex)
	if err != nil {
		genericserver.Fail(c, "claim parsing error", err)
		return
	}

	err = ia.AddClaim(&id, claim)
	if err != nil {
		genericserver.Fail(c, "AddClaim error", err)
		return
	}

	c.JSON(200, gin.H{})
}

type ClaimsMsg struct {
	ClaimsHex []string `json:"claims"`
}

func handlePostClaims(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	var claimsMsg ClaimsMsg
	err = c.BindJSON(&claimsMsg)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	claims, err := core.HexArrayToClaimArray(claimsMsg.ClaimsHex)
	if err != nil {
		genericserver.Fail(c, "claims parsing error", err)
		return
	}
	err = ia.AddClaims(&id, claims)
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
