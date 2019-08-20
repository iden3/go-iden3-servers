package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-core/services/identityagentsrv"
	"github.com/iden3/go-iden3-servers/cmd/genericserver"

	log "github.com/sirupsen/logrus"
)

var ia *identityagentsrv.Service

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func serveServiceApi() *http.Server {
	prefix := "/api/unstable"
	api := gin.Default()
	api.Use(cors.Default())
	serviceapi := api.Group(prefix)

	serviceapi.GET("/info", handleInfo)
	serviceapi.POST("/identity", handlePostIdentity)
	serviceapi.POST("/id/:id/claim", handlePostClaim)
	serviceapi.POST("/id/:id/claims", handlePostClaims)
	serviceapi.GET("/id/:id/claims", handleGetAllClaims)
	serviceapi.GET("/id/:id/claims/emitted", handleGetAllEmittedClaims)
	serviceapi.GET("/id/:id/claims/received", handleGetAllReceivedClaims)
	serviceapi.GET("/id/:id/mt", handleGetFullMT)

	serviceapisrv := &http.Server{Addr: genericserver.C.Server.ServiceApi, Handler: api}
	go func() {
		if err := genericserver.ListenAndServe(serviceapisrv, "Service"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

func Serve(iaSrv *identityagentsrv.Service) {
	ia = iaSrv

	stopch := make(chan interface{})

	// catch ^C to send the stop signal
	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)
	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				stopch <- nil
			}
		}
	}()

	// start servers
	serviceapisrv := serveServiceApi()

	// wait until shutdown signal
	<-stopch
	log.Info("Shutdown Server ...")

	if err := serviceapisrv.Shutdown(context.Background()); err != nil {
		log.Error("ServiceApi Shutdown:", err)
	} else {
		log.Info("ServiceApi stopped")
	}
}
