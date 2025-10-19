package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const contentType = "Content-Type"

// NormalHeaders are the regular Headers used by an HTTP Server for
// request authentication.
var NormalHeaders = &Headers{
	Authenticate:      "WWW-Authenticate",
	Authorization:     "Authorization",
	AuthInfo:          "Authentication-Info",
	UnauthCode:        http.StatusUnauthorized,
	UnauthContentType: "text/plain",
	UnauthResponse:    fmt.Sprintf("%d %s\n", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)),
}

type SecretProvider func(user, realm string) string
type Headers struct {
	Authenticate      string // WWW-Authenticate
	Authorization     string // Authorization
	AuthInfo          string // Authentication-Info
	UnauthCode        int    // 401
	UnauthContentType string // text/plain
	UnauthResponse    string // Unauthorized.
}

type BasicAuth struct {
	Realm   string
	Secrets SecretProvider
	// Headers used by authenticator. Set to ProxyHeaders to use with
	// proxy server. When nil, NormalHeaders are used.
	Headers *Headers
}

func NewBasicAuthenticator(realm string, secrets SecretProvider) *BasicAuth {
	return &BasicAuth{Realm: realm, Secrets: secrets}
}

// V returns NormalHeaders when h is nil, or h otherwise. Allows to
// use uninitialized *Headers values in structs.
func (h *Headers) V() *Headers {
	if h == nil {
		return NormalHeaders
	}
	return h
}

func (a *BasicAuth) CheckAuth(r *http.Request) string {
	s := strings.SplitN(r.Header.Get(a.Headers.V().Authorization), " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return ""
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return ""
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return ""
	}
	user, password := pair[0], pair[1]
	secret := a.Secrets(user, a.Realm)
	if secret == "" {
		return ""
	}
	compare := bcrypt.CompareHashAndPassword
	if compare([]byte(secret), []byte(password)) != nil {
		// fmt.Println("[77] error", secret, password)
		return ""
	}
	// fmt.Println("[80] return success")
	return pair[0]
}

func (a *BasicAuth) RequireAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, a.Headers.V().UnauthContentType)
	w.Header().Set(a.Headers.V().Authenticate, `Basic realm="`+a.Realm+`"`)
	w.WriteHeader(a.Headers.V().UnauthCode)
	w.Write([]byte(a.Headers.V().UnauthResponse))
}
