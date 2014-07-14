package infoauth

import (
	"encoding/json"
	"fmt"
	"github.com/axelmagn/envcfg"
	"log"
	"net/http"
	"net/url"
)

const UserContentKey = "user"
const UserIDKey = "id"

const OauthCodeKey = "code"
const OauthStateKey = "state"
const OuterRedirectKey = "redirect"

func Debug(e error) string {
	debug := GetSetting(S_DEBUG)
	if debug == envcfg.TRUE {
		return e.Error()
	}
	return ""
}

func GetGoogleAuthURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET:\t%s", r.URL.Path)
	defer PanicToError(w)
	url, err := NewGoogleAuthURL()
	if err != nil {
		http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(url))
}

func GetLinkedInAuthURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET:\t%s", r.URL.Path)
	defer PanicToError(w)
	url, err := NewLinkedInAuthURL()
	if err != nil {
		http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(url))
}

func ExchangeCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET:\t%s", r.URL.Path)
	defer PanicToError(w)

	code := r.FormValue(OauthCodeKey)
	if code == "" {
		http.Error(w, "No auth code specified.", http.StatusBadRequest)
		return
	}

	state := r.FormValue(OauthStateKey)
	if state == "" {
		http.Error(w, "No state token specified", http.StatusBadRequest)
		return
	}

	token, h, err := ExchangeCode(code, state)
	if err != nil {
		http.Error(w, "Could not exchange token: "+err.Error(), http.StatusBadRequest)
		return
	}

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		http.Error(w, "Error serializing token.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	var service string
	switch h.Config {
	case C_GOOGLE:
		service = GoogleServiceName
	case C_LINKEDIN:
		service = LinkedInServiceName
	default:
		service = "Unknown"
	}

	// No redirect specified
	if h.OuterRedirect == "" {
		w.Write([]byte("\nService:\t" + service + "\n"))
		w.Write(tokenJSON)
	} else {
		redirect, err := url.Parse(h.OuterRedirect)
		if err != nil {
			http.Error(w, "Could not parse redirect URL.\n"+Debug(err), http.StatusInternalServerError)
			return
		}

		params := redirect.Query()
		params.Set("access_token", token.AccessToken)
		params.Set("service", service)
		params.Set("success", "true")
		redirect.RawQuery = params.Encode()
		http.Redirect(w, r, redirect.String(), http.StatusSeeOther)
	}

}

func GoogleAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET:\t%s", r.URL.Path)
	defer PanicToError(w)
	outerRedirect := r.FormValue(OuterRedirectKey)
	var authUrl string
	var err error
	if outerRedirect == "" {
		authUrl, err = NewGoogleAuthURL()
		if err != nil {
			http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
			return
		}
	} else {
		authUrl, err = NewGoogleAuthURLWithRedirect(outerRedirect)
		if err != nil {
			http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
			return
		}

	}
	http.Redirect(w, r, authUrl, http.StatusSeeOther)
}

func LinkedInAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET:\t%s", r.URL.Path)
	defer PanicToError(w)
	outerRedirect := r.FormValue(OuterRedirectKey)
	var authUrl string
	var err error
	if outerRedirect == "" {
		authUrl, err = NewLinkedInAuthURL()
		if err != nil {
			http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
			return
		}
	} else {
		authUrl, err = NewLinkedInAuthURLWithRedirect(outerRedirect)
		if err != nil {
			http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, authUrl, http.StatusSeeOther)
}

func PanicToError(w http.ResponseWriter) {
	r := recover()
	if r != nil {
		errStr := fmt.Sprintf("Server Panic.\n%v", r)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}
}
