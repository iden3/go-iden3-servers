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

// serveServiceApi start service api calls.
func serveServiceApi(addr string, srv *loaders.Server) *http.Server {
	// api, serviceapi := serve.NewServiceAPI("/api/unstable", srv)
	// serviceapi.POST("/claims", serve.WithServer(srv, handlePostClaim))                  // Get relay claim proof
	// serviceapi.GET("/claims/:hi/proof", serve.WithServer(srv, handleGetClaimProofByHi)) // Get relay claim proof

	// serviceapisrv := &http.Server{Addr: addr, Handler: api}
	// go func() {
	// 	if err := serve.ListenAndServe(serviceapisrv, "Service"); err != nil &&
	// 		err != http.ErrServerClosed {
	// 		log.Fatalf("listen: %s\n", err)
	// 	}
	// }()
	// return serviceapisrv
	return nil
}

// serveAdminApi start admin api calls.
func serveAdminApi(addr string, stopch chan interface{}, srv *loaders.Server) *http.Server {
	api, adminapi := serve.NewAdminAPI("/api/unstable", stopch, srv)
	if adminapi == nil {
		println("IGNORE ME")
	}
	// DEPRECATED
	// adminapi.POST("/claims/basic", serve.WithServer(srv, handleAddClaimBasic))

	adminapisrv := &http.Server{Addr: addr, Handler: api}
	go func() {
		if err := serve.ListenAndServe(adminapisrv, "Admin"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}

// Serve initilization all services and its corresponding api calls.
func Serve(cfg *config.Config, srv *loaders.Server) {

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
	// srv.StateWriter.Start()
	serviceapisrv := serveServiceApi(cfg.Server.ServiceApi, srv)
	adminapisrv := serveAdminApi(cfg.Server.AdminApi, stopch, srv)

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
