package commands

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/iden3/go-iden3-core/services/identityagentsrv"
	"github.com/iden3/go-iden3-servers/cmd/genericserver"
	"github.com/iden3/go-iden3-servers/cmd/identityserver/endpoint"
)

var ServerCommands = []cli.Command{
	{
		Name:    "start",
		Aliases: []string{},
		Usage:   "start the server",
		Action:  cmdStart,
	},
	{
		Name:    "stop",
		Aliases: []string{},
		Usage:   "stops the server",
		Action:  cmdStop,
	},
	{
		Name:    "info",
		Aliases: []string{},
		Usage:   "server status",
		Action:  cmdInfo,
	},
}

func cmdStart(c *cli.Context) error {

	if err := genericserver.MustRead(c); err != nil {
		return err
	}

	storage := genericserver.LoadStorage()
	identityagentService := identityagentsrv.New(storage)

	endpoint.Serve(identityagentService)

	return nil
}

func postAdminApi(command string) (string, error) {

	hostport := strings.Split(genericserver.C.Server.AdminApi, ":")
	if hostport[0] == "0.0.0.0" {
		hostport[0] = "127.0.0.1"
	}
	url := "http://" + hostport[0] + ":" + hostport[1] + "/" + command

	var body bytes.Buffer
	resp, err := http.Post(url, "text/plain", &body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func cmdStop(c *cli.Context) error {
	if err := genericserver.MustRead(c); err != nil {
		return err
	}
	output, err := postAdminApi("stop")
	if err == nil {
		log.Info("Server response: ", output)
	}
	return err
}

func cmdInfo(c *cli.Context) error {
	if err := genericserver.MustRead(c); err != nil {
		return err
	}
	output, err := postAdminApi("info")
	if err == nil {
		log.Info("Server response: ", output)
	}
	return err
}
