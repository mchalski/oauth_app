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

	if port == "" {
		port = "80"
	}

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + EnvPort,
	}

	l.Fatal(srv.ListenAndServe())
}

//
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//TODO check user session/acc token
	l.Infof("GET /")

	challenge := r.FormValue("challenge")
	l.Infof("got consent challenge: %s", challenge)

	redirToHydra(w)
}

// redirect to hydra, which then redirects to the IdP login
func redirToHydra(w http.ResponseWriter) {
	//TODO state min 8chars, generate random
	url := hydraIp +
		`/oauth2/auth?client_id=` +
		EnvHydraClientID +
		`&response_type=code&scope=foo&state=123456780`

	w.Header().Set(`Location`,
		url)

	w.WriteHeader(301)
}
