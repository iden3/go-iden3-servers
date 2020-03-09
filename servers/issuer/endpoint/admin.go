package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-servers/handlers"
	"github.com/iden3/go-iden3-servers/loaders"
	// "errors"
	// "net/http"
	// "github.com/gin-gonic/gin"
	// common3 "github.com/iden3/go-iden3-core/common"
	// "github.com/iden3/go-iden3-servers/cmd/genericserver"
)

type addClaimBasicMsg struct {
	IndexData string
	Data      string
}

func handleSyncIdenStatePublic(c *gin.Context, srv *loaders.Server) {
	if err := srv.Issuer.SyncIdenStatePublic(); err != nil {
		handlers.Fail(c, "SyncIdenStatePublic", err)
		return
	}
	c.JSON(200, gin.H{})
}

// DEPRECATED
// func handleAddClaimBasic(c *gin.Context) {
// 	var m addClaimBasicMsg
// 	err := c.ShouldBindJSON(&m)
// 	if err != nil {
// 		genericserver.Fail(c, "json parsing error", err)
// 		return
// 	}
//
// 	indexData, err := common3.HexDecode(m.IndexData)
// 	if err != nil {
// 		genericserver.Fail(c, "error on handlePostClaim", err)
// 		return
// 	}
// 	data, err := common3.HexDecode(m.Data)
// 	if err != nil {
// 		genericserver.Fail(c, "error on handlePostClaim", err)
// 		return
// 	}
//
// 	if len(indexData) < 400/8 {
// 		genericserver.Fail(c, "error on handlePostClaim", errors.New("indexData smaller than 400/8"))
// 		return
// 	}
// 	if len(data) < 496/8 {
// 		genericserver.Fail(c, "error on handlePostClaim", errors.New("data smaller than 496/8"))
// 		return
// 	}
//
// 	var indexSlot [400 / 8]byte
// 	var dataSlot [496 / 8]byte
// 	copy(indexSlot[:], indexData[:400/8])
// 	copy(dataSlot[:], data[:496/8])
// 	proofOfClaim, err := genericserver.Adminservice.AddClaimBasic(indexSlot, dataSlot)
// 	if err != nil {
// 		c.String(http.StatusBadRequest, err.Error())
// 	}
// 	c.JSON(http.StatusOK, proofOfClaim)
// }
