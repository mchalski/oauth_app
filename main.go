package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var l = logrus.New()

var (
	// serve app on this port
	EnvPortName = "PORT"
	EnvPort     string

	// full hydra ip/port
	EnvHydraURLName = "HYDRA"
	EnvHydraURL     string

	// oauth client params
	// for the poc we most likely use the single root client both for the idp and this app
	EnvHydraClientIDName = "CLIENTID"
	EnvHydraClientID     string

	EnvHydraClientSecretName = "CLIENTSECRET"
	EnvHydraClientSecret     string
)

//
func getEnv() {
	// serve app on this port
	EnvPort = os.Getenv(EnvPortName)

	// full hydra ip/port
	EnvHydraURL = os.Getenv(EnvHydraURLName)

	// oauth app params
	EnvHydraClientID = os.Getenv(EnvHydraClientIDName)
	EnvHydraClientSecret = os.Getenv(EnvHydraClientSecretName)

	if EnvPort == "" {
		l.Fatal("Need env var ", EnvPortName)
	}

	if EnvHydraURL == "" {
		l.Fatal("Need env var ", EnvHydraURLName)
	}

	if EnvHydraClientID == "" {
		l.Fatal("Need env var ", EnvHydraClientIDName)
	}

	if EnvHydraClientSecret == "" {
		l.Fatal("Need env var ", EnvHydraClientSecretName)
	}
}

//
func main() {
	getEnv()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + EnvPort,
	}

	l.Fatal(srv.ListenAndServe())
}

// renders the home screen, if user unauthorized - redirect to hydra/idp to get a token
// step 1 of the hydra flow
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//TODO check user session/acc token
	l.Infof("GET /")

	challenge := r.FormValue("challenge")
	l.Infof("got consent challenge: %s", challenge)

	redirToHydra(w)
}

// TODO token accept callback for step 4 of the flow

// redirect to hydra, which then redirects to the IdP login
func redirToHydra(w http.ResponseWriter) {
	//TODO
	// - state min 8chars, generate random
	// - set some reasonable scope
	url := EnvHydraURL +
		`/oauth2/auth?client_id=` +
		EnvHydraClientID +
		`&response_type=code&scope=dummy&state=123456780`

	w.Header().Set(`Location`,
		url)

	w.WriteHeader(301)
}
