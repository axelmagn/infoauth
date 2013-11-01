package infoauth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"github.com/axelmagn/envcfg"
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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue(UserIDKey)
	if idStr == "" {
		http.Error(w, "No User ID Specified.", http.StatusBadRequest)
		return
	}

	id64, err := strconv.ParseUint(idStr, 10, 0)
	id := uint(id64)
	if err != nil {
		http.Error(w, "Error parsing User id.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	u, err := GetUser(id)
	if err != nil {
		http.Error(w, "Error retrieving User.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	if u == nil {
		http.Error(w, "User does not exist.", http.StatusBadRequest)
		return
	}

	raw, err := u.Value()
	if err != nil {
		http.Error(w, "Error encoding User.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	w.Write(raw)
}

func SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	raw := r.FormValue(UserContentKey)
	if raw == "" {
		http.Error(w, "No user submitted", http.StatusBadRequest)
		return
	}

	u, err := DecodeUser([]byte(raw))
	if err != nil {
		http.Error(w, "Could not parse User.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	if u.ID == 0 {
		u.ID, err = NewUserID()
		if err != nil {
			http.Error(w, "Could not assign User ID.\n"+Debug(err), http.StatusInternalServerError)
			return
		}
	}

	err = u.Save()
	if err != nil {
		http.Error(w, "Could not save User.\n"+Debug(err), http.StatusInternalServerError)
		return
	}

	v, err := u.Value()
	if err != nil {
		http.Error(w, "Could not decode user after saving.\n"+Debug(err), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

// takes a regex specifying path groups
// assumes that 2nd group is the user id
// returns a handler function for users
func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Printf("GET:\t%s", r.URL.Path)
		GetUserHandler(w, r)
	case "POST":
		log.Printf("POST:\t%s", r.URL.Path)
		SaveUserHandler(w, r)
	case "PUT":
		log.Printf("PUT:\t%s", r.URL.Path)
		SaveUserHandler(w, r)
	}
}

func GetGoogleAuthURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET:\t%s", r.URL.Path)
	url, err := NewGoogleAuthURL()
	if err != nil {
		http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(url))
}

func GetLinkedInAuthURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET:\t%s", r.URL.Path)
	url, err := NewLinkedInAuthURL()
	if err != nil {
		http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(url))
}

func ExchangeCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET:\t%s", r.URL.Path)
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

	// TODO: rewrite this to get a handshake so that it knows what service it's using
	token, h, err := ExchangeCode(code, state)
	if err != nil {
		http.Error(w, "Could not exchange token: " + err.Error(), http.StatusBadRequest)
		return
	}

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		http.Error(w, "Error serializing token.\n" + Debug(err), http.StatusInternalServerError)
		return
	}

	var service string
    switch h.Config {
    case C_GOOGLE:
        service = GoogleServiceName
    case C_LINKEDIN:
        service = GoogleServiceName
    default:
    	service = "Unknown"
    }

    // No redirect specified
    if h.OuterRedirect == "" {
	    w.Write([]byte("\nService:\t"+service+"\n"))
		w.Write(tokenJSON)
	} else {
		redirect, err := url.Parse(h.OuterRedirect)
		if err != nil {
			http.Error(w, "Could not parse redirect URL.\n" + Debug(err), http.StatusInternalServerError)
			return
		}

		params := redirect.Query()
		params.Set("access_token", token.AccessToken)
		params.Set("service", service)
		redirect.RawQuery = params.Encode()
		http.Redirect(w, r, redirect.String(), http.StatusSeeOther)
	}

}

func GoogleAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET:\t%s", r.URL.Path)
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
	outerRedirect := r.FormValue(OuterRedirectKey)
	if outerRedirect == "" {	
		authUrl, err := NewLinkedInAuthURL()
		if err != nil {
			http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
			return
		}
	    http.Redirect(w, r, authUrl, http.StatusOK)
	    return
	} else {
		authUrl, err := NewLinkedInAuthURLWithRedirect(outerRedirect)
		if err != nil {
			http.Error(w, "Could not generate Authentication URL.\n"+Debug(err), http.StatusInternalServerError)
			return
		}
	    http.Redirect(w, r, authUrl, http.StatusSeeOther)
	    return

	}
}

