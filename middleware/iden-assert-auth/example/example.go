package main

import (
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/eth"
	"github.com/iden3/go-iden3-core/services/discoverysrv"
	"github.com/iden3/go-iden3-core/services/ethsrv"
	"github.com/iden3/go-iden3-core/services/nameresolversrv"
	"github.com/iden3/go-iden3-core/services/signedpacketsrv"
	"github.com/iden3/go-iden3-servers/middleware/iden-assert-auth"
)

func handleGetHello(c *gin.Context) {
	user := auth.GetUser(c)
	c.JSON(200, gin.H{
		"id":      common3.HexEncode(user.Id[:]),
		"ethName": user.EthName,
		"text":    "Hello World.",
	})
}

func main() {
	nonceDb := core.NewNonceDb()
	domain := "test.eth"

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	nameResolverService, err := nameresolversrv.New("/tmp/go-iden3/names.json")
	if err != nil {
		log.Fatal(err)
	}
	discoveryService, err := discoverysrv.New("/tmp/go-iden3/identitites.json")
	if err != nil {
		log.Fatal(err)
	}
	// Example ethereum client that won't work because it doesn't use any keystore nor gateway
	client := eth.NewClient2(nil, nil, nil)
	ethService := ethsrv.New(client, ethsrv.ContractAddresses{
		RootCommits: common.HexToAddress("0xAF8DC5663cf3890DF4E236cdA0718f4ecB8b42f5")},
	)
	signedPacketVerifier := signedpacketsrv.NewSignedPacketVerifier(discoveryService, nameResolverService, ethService)
	authapi, err := auth.AddAuthMiddleware(&r.RouterGroup, domain, nonceDb, []byte("password"),
		signedPacketVerifier)
	if err != nil {
		log.Fatal(err)
	}

	authapi.GET("/hello", handleGetHello)

	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}
