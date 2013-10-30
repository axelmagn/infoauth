package infoauth

import (
	"net/http"
	"strconv"
	"github.com/axelmagn/envcfg"
	"log"
)

var UserContentKey = "user"
var UserIDKey = "id"

func Debug(e error) string {
	debug := GetSetting(S_DEBUG)
	if debug == envcfg.TRUE {
		return e.Error()
	}
	return ""
}

func GetUserHandler(w http.ResponseWriter, r *http.Request)  {
	idStr := r.FormValue(UserIDKey)
	if idStr == "" {
		http.Error(w, "No User ID Specified.", http.StatusBadRequest)
		return
	}


	id64, err := strconv.ParseUint(idStr, 10, 0)
	id := uint(id64)
	if err != nil {
		http.Error(w, "Error parsing User id.\n" + Debug(err), http.StatusInternalServerError)
		return
	} 

	u, err := GetUser(id)
	if err != nil {
		http.Error(w, "Error retrieving User.\n" + Debug(err), http.StatusInternalServerError)
		return
	} 

	if u == nil {
		http.Error(w, "User does not exist.", http.StatusBadRequest)
		return
	} 

	raw, err := u.Value()
	if err != nil {
		http.Error(w, "Error encoding User.\n" + Debug(err), http.StatusInternalServerError)
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
		http.Error(w, "Could not parse User.\n" + Debug(err), http.StatusInternalServerError)
		return
	}

	if u.ID == 0 {
		u.ID, err = NewUserID()
		if err != nil {
			http.Error(w, "Could not assign User ID.\n" + Debug(err), http.StatusInternalServerError)
			return
		}
	}

	err = u.Save()
	if err != nil {
		http.Error(w, "Could not save User.\n" + Debug(err), http.StatusInternalServerError)
		return
	}

	v, err := u.Value()
	if err != nil {
		http.Error(w, "Could not decode user after saving.\n" + Debug(err), http.StatusInternalServerError)
		return
	}
	w.Write(v)
}

// takes a regex specifying path groups
// assumes that 2nd group is the user id
// returns a handler function for users
func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch(r.Method) {
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
