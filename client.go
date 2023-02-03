package waygate

import (
	"context"
	"fmt"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"net/http"
)

var ClientStoreFactory = func() ClientStore {
	return NewClientJsonStore()
}

type ClientStore interface {
	SetState(string)
	GetState() string
	GetAccessToken() (string, error)
	SetAccessToken(token string)
}

func FlowToken(serverAddr, bindAddr string) (string, error) {

	db := ClientStoreFactory()

	token, err := db.GetAccessToken()
	if err != nil {
		outOfBand := false
		oauthConf := buildOauthConfig(serverAddr, outOfBand, bindAddr)
		requestId, _ := GenRandomCode()
		db.SetState(requestId)
		oauthUrl := oauthConf.AuthCodeURL(requestId, oauth2.AccessTypeOffline)

		fmt.Println(oauthUrl)

		browser.OpenURL(oauthUrl)

		mux := http.NewServeMux()

		srv := &http.Server{
			Addr:    bindAddr,
			Handler: mux,
		}

		mux.HandleFunc("/waygate/callback", func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				go srv.Shutdown(context.Background())
			}()

			r.ParseForm()

			code := r.Form.Get("code")
			if code == "" {
				w.WriteHeader(400)
				fmt.Fprintf(w, "Missing code param")
				return
			}

			state := r.Form.Get("state")
			pendingState := db.GetState()

			db.SetState("")

			if state != pendingState {
				w.WriteHeader(400)
				fmt.Fprintf(w, "State does not match")
				return
			}

			ctx := context.Background()
			tok, err := oauthConf.Exchange(ctx, code)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			token = tok.AccessToken
			db.SetAccessToken(token)
		})

		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}

	return token, nil
}

func buildOauthConfig(providerUri string, outOfBand bool, bindAddr string) *oauth2.Config {

	oauthConf := &oauth2.Config{
		ClientID:     bindAddr,
		ClientSecret: "fake-secret",
		Scopes:       []string{"tunnel"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://%s/waygate/authorize", providerUri),
			TokenURL: fmt.Sprintf("https://%s/waygate/token", providerUri),
		},
	}

	if outOfBand {
		oauthConf.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"
	} else {
		oauthConf.RedirectURL = fmt.Sprintf("%s/waygate/callback", bindAddr)
	}

	return oauthConf
}
