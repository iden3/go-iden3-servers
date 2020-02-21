package cmd

import (
	"fmt"

	"github.com/iden3/go-iden3-core/components/httpclient"
	"github.com/iden3/go-iden3-servers/config"
	log "github.com/sirupsen/logrus"
)

func PostAdminApi(cfgServer *config.Server, path string, result interface{}) error {
	httpClient := httpclient.NewHttpClient(fmt.Sprintf("http://%s/api/unstable", cfgServer.AdminApi))
	log.WithFields(log.Fields{
		"path": path,
	}).Info("Posting admin api")
	if result == nil {
		m := make(map[string]interface{})
		result = &m
	}
	if err := httpClient.DoRequest(httpClient.NewRequest().Path(path).Post(""), result); err != nil {
		return fmt.Errorf("Failed http request: %w", err)
	}
	log.WithFields(log.Fields{
		"result": result,
	}).Info("Post admin api")
	return nil
}
