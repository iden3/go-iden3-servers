package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-core/components/idenadminutils"
	"github.com/iden3/go-iden3-core/components/idenmanager"
	"github.com/iden3/go-iden3-core/services/idenstatewriter"
	"github.com/iden3/go-iden3-servers/cmd/genericserver"

	log "github.com/sirupsen/logrus"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// serveServiceApi start service api calls.
func serveServiceApi() *http.Server {
	api, serviceapi := genericserver.NewServiceAPI("/api/unstable")
	serviceapi.POST("/claims", handlePostClaim)                  // Get relay claim proof
	serviceapi.GET("/claims/:hi/proof", handleGetClaimProofByHi) // Get relay claim proof

	serviceapisrv := &http.Server{Addr: genericserver.C.Server.ServiceApi, Handler: api}
	go func() {
		if err := genericserver.ListenAndServe(serviceapisrv, "Service"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

// serveAdminApi start admin api calls.
func serveAdminApi(stopch chan interface{}) *http.Server {
	api, adminapi := genericserver.NewAdminAPI("/api/unstable", stopch)
	if adminapi == nil {
		println("IGNORE ME")
	}
	// DEPRECATED
	// adminapi.POST("/claims/basic", handleAddClaimBasic)

	adminapisrv := &http.Server{Addr: genericserver.C.Server.AdminApi, Handler: api}
	go func() {
		if err := genericserver.ListenAndServe(adminapisrv, "Admin"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}

// Serve initilization all services and its corresponding api calls.
func Serve(rs idenstatewriter.IdenStateWriter, cs *idenmanager.IdenManager, as *idenadminutils.IdenAdminUtils) {

	genericserver.Claimservice = cs
	genericserver.Rootservice = rs
	genericserver.Adminservice = as

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

	// start servers.
	genericserver.Rootservice.Start()
	serviceapisrv := serveServiceApi()
	adminapisrv := serveAdminApi(stopch)

	// wait until shutdown signal.
	<-stopch
	log.Info("Shutdown Server ...")

	if err := serviceapisrv.Shutdown(context.Background()); err != nil {
		log.Error("ServiceApi Shutdown:", err)
	} else {
		log.Info("ServiceApi stopped")
	}

	if err := adminapisrv.Shutdown(context.Background()); err != nil {
		log.Error("AdminApi Shutdown:", err)
	} else {
		log.Info("AdminApi stopped")
	}

}
