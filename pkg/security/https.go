package security

import (
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
)

var tlsEnabled = false 

// UseTLS retourne vrai si TLS est activé.
func UseTLS() bool {
	return tlsEnabled
}

// ConfigureTLS retourne la configuration TLS pour le serveur HTTP.
func ConfigureTLS(mux http.Handler) *autocert.Manager {
	if !tlsEnabled {
		log.Println("TLS désactivé, démarrage en HTTP")
		return nil
	}
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("domaine.com"),
		Cache:      autocert.DirCache("certs"),
	}
	return m
}
