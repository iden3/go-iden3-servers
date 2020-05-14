package handlers

import (
	"fmt"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Fail(c *gin.Context, msg string, err error) {
	if err != nil {
		log.WithError(err).Error(msg)
		msg = fmt.Sprintf("%v: %v", msg, err)
	} else {
		log.Error(msg)
	}
	c.JSON(400, gin.H{
		"error": msg,
	})
	return
}

func HandleStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "online",
	})
}

// Generic
// func HandleGetRoot(c *gin.Context, srv *loaders.Server) {
// 	// get the contract data
// 	root, err := srv.StateWriter.GetRoot(srv.Issuer.ID())
// 	if err != nil {
// 		Fail(c, "error contract.GetRoot(C.Keys.Ethereum.KUpdateRoot)", err)
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"root":         srv.Manager.MT().RootKey().Hex(),
// 		"contractRoot": common3.HexEncode(root.Root[:]),
// 	})
// }

// TODO: Redo once IdenStateReader.Info() is implemented
// Admin
// func HandleInfo(c *gin.Context, im *idenmanager.IdenManager) {
// 	r := Adminservice.Info(im.ID())
//
// 	c.JSON(200, gin.H{
// 		"info": r,
// 	})
// }

// func HandleRawDump(c *gin.Context, srv *loaders.Server) {
// 	srv.AdminUtils.RawDump(c)
// }

// func HandleRawImport(c *gin.Context, srv *loaders.Server) {
// 	var data map[string]string
// 	err := c.ShouldBindJSON(&data)
// 	if err != nil {
// 		Fail(c, "json parsing error", err)
// 		return
// 	}
//
// 	count, err := srv.AdminUtils.RawImport(data)
// 	if err != nil {
// 		c.String(http.StatusBadRequest, err.Error())
// 	}
// 	c.String(http.StatusOK, "imported "+strconv.Itoa(count)+" key,value entries")
// }

// func HandleClaimsDump(c *gin.Context, srv *loaders.Server) {
// 	r := srv.AdminUtils.ClaimsDump()
// 	c.JSON(http.StatusOK, r)
// }
