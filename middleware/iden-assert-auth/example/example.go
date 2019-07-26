package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/services/discoverysrv"
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
	discoveryservice, err := discoverysrv.New("/tmp/go-iden3/identitites.json")
	if err != nil {
		log.Fatal(err)
	}
	signedPacketVerifier := signedpacketsrv.NewSignedPacketVerifier(discoveryservice, nameResolverService)
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
