package main

import (
	_ "embed"
	"encoding/json"

	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-app-servicenow/config"
	"github.com/mattermost/mattermost-app-servicenow/routers/mattermost"
	"github.com/mattermost/mattermost-app-servicenow/routers/oauth"
	"github.com/mattermost/mattermost-plugin-apps/apps"
)

const (
	baseURLPosition = 1
	addressPosition = 2
)

//go:embed manifest.json
var manifestSource []byte

func main() {
	var manifest apps.Manifest

	err := json.Unmarshal(manifestSource, &manifest)
	if err != nil {
		panic("failed to load manfest: " + err.Error())
	}

	localMode := os.Getenv("LOCAL") == "true"

	// Init routers
	r := mux.NewRouter()
	mattermost.Init(r, &manifest, localMode)
	oauth.Init(r)

	http.Handle("/", r)

	if localMode {
		if len(os.Args) > baseURLPosition {
			config.SetBaseURL(os.Args[baseURLPosition])
		}

		addr := ":3000"
		if len(os.Args) > addressPosition {
			addr = os.Args[addressPosition]
		}

		manifest.HTTPRootURL = addr
		manifest.Type = apps.AppTypeHTTP

		_ = http.ListenAndServe(addr, nil)

		return
	}

	lambda.Start(httpadapter.New(r).Proxy)
}
