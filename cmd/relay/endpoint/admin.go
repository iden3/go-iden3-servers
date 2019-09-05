package endpoint

import (
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-servers/cmd/genericserver"
)

func handleMimc7(c *gin.Context) {
	var elements []*big.Int
	err := c.BindJSON(&elements)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	r, err := genericserver.Adminservice.Mimc7(elements)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.String(http.StatusOK, r.String())
}

type addClaimBasicMsg struct {
	Namespace string
	IndexData string
	Data      string
}

// DEPRECATED
// func handleAddClaimBasic(c *gin.Context) {
// 	var m addClaimBasicMsg
// 	err := c.BindJSON(&m)
// 	if err != nil {
// 		genericserver.Fail(c, "json parsing error", err)
// 		return
// 	}
//
// 	if len(m.IndexData) != 400/8 {
// 		c.String(http.StatusBadRequest, "indexData smaller than 400/8")
// 		return
// 	}
// 	if len(m.Data) != 496/8 {
// 		c.String(http.StatusBadRequest, "data smaller than 496/8")
// 		return
// 	}
//
// 	var indexSlot [400 / 8]byte
// 	var dataSlot [496 / 8]byte
// 	copy(indexSlot[:], m.IndexData[:400/8])
// 	copy(dataSlot[:], m.Data[:496/8])
// 	proofClaim, err := genericserver.Adminservice.AddClaimBasic(indexSlot, dataSlot)
// 	if err != nil {
// 		c.String(http.StatusBadRequest, err.Error())
// 	}
// 	c.JSON(http.StatusOK, proofClaim)
// }
