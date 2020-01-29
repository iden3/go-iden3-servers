package handlers

import (
	"net/http"
	"strconv"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-servers/loaders"
	log "github.com/sirupsen/logrus"
)

func Fail(c *gin.Context, msg string, err error) {
	if err != nil {
		log.WithError(err).Error(msg)
	} else {
		log.Error(msg)
	}
	c.JSON(400, gin.H{
		"error": msg,
	})
	return
}

// Generic
func HandleGetRoot(c *gin.Context, iden *loaders.Identity) {
	// get the contract data
	root, err := iden.StateWriter.GetRoot(iden.Manager.ID())
	if err != nil {
		Fail(c, "error contract.GetRoot(C.Keys.Ethereum.KUpdateRoot)", err)
		return
	}
	c.JSON(200, gin.H{
		"root":         iden.Manager.MT().RootKey().Hex(),
		"contractRoot": common3.HexEncode(root.Root[:]),
	})
}

// TODO: Redo once IdenStateReader.Info() is implemented
// Admin
// func HandleInfo(c *gin.Context, im *idenmanager.IdenManager) {
// 	r := Adminservice.Info(im.ID())
//
// 	c.JSON(200, gin.H{
// 		"info": r,
// 	})
// }

func HandleRawDump(c *gin.Context, iden *loaders.Identity) {
	iden.AdminUtils.RawDump(c)
}

func HandleRawImport(c *gin.Context, iden *loaders.Identity) {
	var data map[string]string
	err := c.BindJSON(&data)
	if err != nil {
		Fail(c, "json parsing error", err)
		return
	}

	count, err := iden.AdminUtils.RawImport(data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.String(http.StatusOK, "imported "+strconv.Itoa(count)+" key,value entries")
}

func HandleClaimsDump(c *gin.Context, iden *loaders.Identity) {
	r := iden.AdminUtils.ClaimsDump()
	c.JSON(http.StatusOK, r)
}
