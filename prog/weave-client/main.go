package main

import (
	"fmt"
	"os"
	"os/exec"
	"encoding/json"
	"io"
	"net/http"
	"crypto/tls"
	"strings"

	logrus "github.com/sirupsen/logrus"
	util "github.com/AVENTER-UG/util/util"
	cfg "github.com/AVENTER-UG/weave/prog/weave-client/types"
)

var config cfg.Config

func init() {
	config.APIURL = util.Getenv("WEAVE_BACKEND_API", "root")
	config.APIToken = util.Getenv("WEAVE_BACKEND_API_TOKEN", "")
	config.SSL = stringToBool(util.Getenv("WEAVE_BACKEND_SSL", "false"))

	protocol := "https"
	if !config.SSL {
		protocol = "http"
	}

	client := &http.Client{}
	// #nosec G402
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest("GET", protocol+"://"+config.APIURL, nil)
	if len(config.APIToken) > 0 {
		req.Header.Set("Authorization", config.APIToken)
	}

	res, err := client.Do(req)

	if err != nil {
		logrus.WithField("func", "init").Error("Could not connect to backend: ", err.Error())
		panic(err)
	}

	if res.StatusCode == http.StatusOK {
		defer res.Body.Close()
	}

	body, _ := io.ReadAll(res.Body)

	// TODO
	var net cfg.WeaveNetwork
	if err := json.Unmarshal(body, &net); err != nil {
		panic(err)
	}

	fmt.Printf("Weave Network: %+v\n", net)
}


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
