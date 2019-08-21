package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-core/services/adminsrv"
	"github.com/iden3/go-iden3-core/services/claimsrv"
	"github.com/iden3/go-iden3-core/services/counterfactualsrv"
	"github.com/iden3/go-iden3-core/services/identitysrv"
	"github.com/iden3/go-iden3-core/services/rootsrv"
	"github.com/iden3/go-iden3-servers/cmd/genericserver"

	log "github.com/sirupsen/logrus"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func serveServiceApi() *http.Server {
	api, serviceapi := genericserver.NewServiceAPI("/api/unstable")

	// TODO: Deprecate handleGetClaimProofByHi
	serviceapi.GET("/claims/:hi/proof", handleGetClaimProofByHi) // Get relay claim proof
	// NEW Agent API
	serviceapi.GET("/claims/:hi/proof0", handleGetClaimProofByHiBlockchain) // Get relay claim proof (to Blockchain)

	serviceapi.POST("/ids", handleCreateIdGenesis)
	serviceapi.POST("/counterfactuals", handleCreateCounterfactual)
	serviceapi.GET("/counterfactuals/:ethaddr", handleGetCounterfactual)
	serviceapi.POST("/counterfactuals/:ethaddr/deploy", handleDeployCounterfactual)
	serviceapi.POST("/counterfactuals/:ethaddr/forward", handleForwardCounterfactual)
	// NEW Agent API
	serviceapi.GET("/ids/:id/setrootclaim", handleGetSetRootClaim)
	serviceapi.GET("/ids/:id/root", handleGetIdRoot)
	// NEW Agent API
	serviceapi.POST("/ids/:id/setrootclaim", handleUpdateSetRootClaim)
	serviceapi.POST("/ids/:id/root", handleCommitNewIdRoot)
	serviceapi.POST("/ids/:id/claims", handlePostClaim)
	serviceapi.GET("/ids/:id/claims/:hi/proof", handleGetClaimProofUserByHi) // Get user claim proof

	serviceapisrv := &http.Server{Addr: genericserver.C.Server.ServiceApi, Handler: api}
	go func() {
		if err := genericserver.ListenAndServe(serviceapisrv, "Service"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

func serveAdminApi(stopch chan interface{}) *http.Server {
	api, adminapi := genericserver.NewAdminAPI("/api/unstable", stopch)
	adminapi.POST("/mimc7", handleMimc7)
	adminapi.POST("/claims/basic", handleAddClaimBasic)

	adminapisrv := &http.Server{Addr: genericserver.C.Server.AdminApi, Handler: api}
	go func() {
		if err := genericserver.ListenAndServe(adminapisrv, "Admin"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}

func Serve(rs rootsrv.Service, cs claimsrv.Service, ids identitysrv.Service, counterfs counterfactualsrv.Service, as adminsrv.Service) {

	genericserver.Idservice = ids
	genericserver.Counterfactualservice = counterfs
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

	// start servers
	genericserver.Rootservice.Start()
	serviceapisrv := serveServiceApi()
	adminapisrv := serveAdminApi(stopch)

	// wait until shutdown signal
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
