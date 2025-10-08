package main

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	util "github.com/AVENTER-UG/util/util"
	cfg "github.com/AVENTER-UG/weave/prog/weave-client/types"
	logrus "github.com/sirupsen/logrus"
)

var config cfg.Config

func init() {
	config.APIHost = util.Getenv("WEAVE_BACKEND_API_HOST", "root")
	config.APIToken = util.Getenv("WEAVE_BACKEND_API_TOKEN", "")
	config.SSL = stringToBool(util.Getenv("WEAVE_BACKEND_SSL", "false"))

	setWeaveConfiguration()
}

// setWeaveConfiguration get the configuration from backen and set
// the config for weaver
func setWeaveConfiguration() {
	body, err := httpRequest("GET", "/api/weave/v0/network")

	if err != nil {
		logrus.WithField("func", "weave-client.getWeaveConfiguration").Error(err.Error())
		return
	}

	var net cfg.WeaveNetwork
	if err := json.Unmarshal(body, &net); err != nil {
		panic(err)
	}

	// read all peers
	var peers strings.Builder
	for i, p := range net.Peers {
		if i > 0 {
			peers.WriteByte(' ')
		}
		peers.WriteString(p.Hostname)
	}

	os.Setenv("WEAVE_PASSWORD", net.EncryptionPassword)
	os.Setenv("IPALLOC_RANGE", net.OverlayCIDR)
	os.Setenv("WEAVE_PEERS=", peers.String())
}

// httpRequest is the global function to connect with the backend
func httpRequest(method, url string) ([]byte, error) {
	protocol := "https"
	if !config.SSL {
		protocol = "http"
	}

	client := &http.Client{}
	// #nosec G402
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest(method, protocol+"://"+config.APIHost+url, nil)
	if len(config.APIToken) > 0 {
		req.Header.Set("Authorization", config.APIToken)
	}

	res, err := client.Do(req)

	if err != nil {
		logrus.WithField("func", "init").Error("Could not connect to backend: ", err.Error())
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		defer res.Body.Close()
	}

	return io.ReadAll(res.Body)
}

// stringToBool convert true/false string to boolean
func stringToBool(par string) bool {
	return strings.Compare(par, "true") == 0
}

func main() {
	cmd := exec.Command("/home/weave/weaver", os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logrus.WithField("func", "main").Error("waver error", err.Error())
		os.Exit(1)
	}
}
