package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"

	"github.com/iden3/go-iden3-servers/config"
	"github.com/iden3/go-iden3-servers/loaders"
	"github.com/iden3/go-iden3-servers/serve"

	log "github.com/sirupsen/logrus"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func serveServiceApi(addr string, iden *loaders.Identity) *http.Server {
	api, serviceapi := serve.NewServiceAPI("/api/unstable", iden)

	// NEW Agent API
	serviceapi.GET("/claims/:hi/proof0", serve.WithIden(iden, handleGetClaimProofByHiBlockchain)) // Get relay claim proof (to Blockchain)

	// NEW Agent API
	serviceapi.GET("/ids/:id/setrootclaim", serve.WithIden(iden, handleGetSetRootClaim))
	// NEW Agent API
	serviceapi.POST("/ids/:id/setrootclaim", serve.WithIden(iden, handleUpdateSetRootClaim))
	serviceapi.GET("/ids/:id/root", serve.WithIden(iden, handleGetIdRoot))
	serviceapi.POST("/ids/:id/root", serve.WithIden(iden, handleCommitNewIdRoot))

	serviceapisrv := &http.Server{Addr: addr, Handler: api}
	go func() {
		if err := serve.ListenAndServe(serviceapisrv, "Service"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

func serveAdminApi(addr string, stopch chan interface{}, iden *loaders.Identity) *http.Server {
	api, adminapi := serve.NewAdminAPI("/api/unstable", stopch, iden)
	adminapi.POST("/mimc7", serve.WithIden(iden, handleMimc7))

	adminapisrv := &http.Server{Addr: addr, Handler: api}
	go func() {
		if err := serve.ListenAndServe(adminapisrv, "Admin"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}

func Serve(cfg *config.Config, iden *loaders.Identity) {

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
	iden.StateWriter.Start()
	serviceapisrv := serveServiceApi(cfg.Server.ServiceApi, iden)
	adminapisrv := serveAdminApi(cfg.Server.AdminApi, stopch, iden)

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
