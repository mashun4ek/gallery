package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	localctx "github.com/mashun4ek/webdevcalhoun/gallery/context"
	"github.com/mashun4ek/webdevcalhoun/gallery/models"
	"golang.org/x/oauth2"
)

type OAuths struct {
	os      models.OAuthService
	configs map[string]*oauth2.Config
}

func NewOAuths(os models.OAuthService, configs map[string]*oauth2.Config) *OAuths {
	return &OAuths{
		os:      os,
		configs: configs,
	}
}

func (o *OAuths) Connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}
	// create csrf token
	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OAuths) Callback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth Service", http.StatusInternalServerError)
		return
	}
	r.ParseForm()
	state := r.FormValue("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "Invalid state provided", http.StatusBadRequest)
		return
	}
	// got here -> state is correct, delete cookie
	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)

	// exchange code for a token
	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := localctx.User(r.Context())
	existing, err := o.os.Find(user.ID, service)
	if err == models.ErrNotFound {

	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		o.os.Delete(existing.ID)
	}
	userOAuth := models.OAuth{
		UserID:  user.ID,
		Token:   *token,
		Service: service,
	}
	err = o.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%+v", token)
}

func (o *OAuths) DropboxTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth Service", http.StatusInternalServerError)
		return
	}
	r.ParseForm()
	path := r.FormValue("path")

	user := localctx.User(r.Context())
	userOAuth, err := o.os.Find(user.ID, service)
	if err != nil {
		panic(err)
	}
	token := userOAuth.Token

	data := struct {
		Path string `json:"path"`
	}{
		Path: path,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	client := oauthConfig.Client(context.TODO(), &token)
	req, err := http.NewRequest(http.MethodPost, "https://api.dropboxapi.com/2/files/list_folder", bytes.NewReader(dataBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
	// folders, files, err := dbx.List(token.AccessToken, path)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Fprintln(w, "Folders=", folders)
	// fmt.Fprintln(w, "Files=", files)
}
