package waygate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AuthRequest struct {
	ClientId    string
	RedirectUri string
	Scope       string
	State       string
}

type Oauth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type ServerDatabase interface {
	GetWaygateTokenData(string) (TokenData, error)
	GetWaygate(string) (Waygate, error)
	GetWaygateTokenByCode(code string) (string, error)
}

type Server struct {
	db  ServerDatabase
	mux *http.ServeMux
}

func NewServer(db ServerDatabase) *Server {

	s := &Server{
		db: db,
	}

	mux := &http.ServeMux{}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		s.token(w, r)
	})

	mux.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		s.open(w, r)
	})

	s.mux = mux

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func ExtractAuthRequest(r *http.Request) (*AuthRequest, error) {
	r.ParseForm()

	clientId := r.Form.Get("client_id")
	if clientId == "" {
		return nil, errors.New("Missing client_id param")
	}

	redirectUri := r.Form.Get("redirect_uri")
	if redirectUri == "" {
		return nil, errors.New("Missing redirect_uri param")
	}

	if !strings.HasPrefix(redirectUri, clientId) && redirectUri != "urn:ietf:wg:oauth:2.0:oob" {
		return nil, errors.New("redirect_uri must be on the same domain as client_id")
	}

	scope := r.Form.Get("scope")
	if scope == "" {
		return nil, errors.New("Missing scope param")
	}

	state := r.Form.Get("state")
	if state == "" {
		return nil, errors.New("state param can't be empty")
	}

	req := &AuthRequest{
		ClientId:    clientId,
		RedirectUri: redirectUri,
		Scope:       scope,
		State:       state,
	}

	return req, nil
}

func (s *Server) token(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	code := r.Form.Get("code")

	token, err := s.db.GetWaygateTokenByCode(code)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, err.Error())
		return
	}

	resp := Oauth2TokenResponse{
		AccessToken: token,
		TokenType:   "bearer",
	}

	jsonStr, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonStr)
}

func (s *Server) open(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Invalid method")
		return
	}

	token, err := extractToken("token", r)
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, err.Error())
		return
	}

	tokenData, err := s.db.GetWaygateTokenData(token)
	if err != nil {
		w.WriteHeader(403)
		fmt.Fprintf(w, err.Error())
		return
	}

	//waygate, err := s.db.GetWaygate(tokenData.WaygateId)
	_, err = s.db.GetWaygate(tokenData.WaygateId)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	r.ParseForm()

	tunnelType := r.Form.Get("type")

	if tunnelType != "wireguard" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "No supported tunnel types")
		return
	}

	//if s.SshConfig == nil {
	//	w.WriteHeader(500)
	//	fmt.Fprintf(w, "No SSH config set")
	//	return
	//}

	//deleteFromAuthorizedKeys(s.SshConfig.AuthorizedKeysPath, tokenData.WaygateId)

	//tunnelPort, err := randomOpenPort()
	//if err != nil {
	//	w.WriteHeader(500)
	//	fmt.Fprintf(w, err.Error())
	//	return
	//}

	//privKey, err := addToAuthorizedKeys(tokenData.WaygateId, tunnelPort, false, s.SshConfig.AuthorizedKeysPath)
	//if err != nil {
	//	w.WriteHeader(500)
	//	fmt.Fprintf(w, err.Error())
	//	return
	//}

	//tun := SSHTunnel{
	//	TunnelType:       "ssh",
	//	Domains:          waygate.Domains,
	//	ServerAddress:    s.SshConfig.ServerAddress,
	//	ServerPort:       s.SshConfig.ServerPort,
	//	ServerTunnelPort: tunnelPort,
	//	Username:         s.SshConfig.Username,
	//	ClientPrivateKey: privKey,
	//}

	//json.NewEncoder(w).Encode(tun)
}

// Looks for auth token in query string, then headers, then cookies
func extractToken(tokenName string, r *http.Request) (string, error) {

	query := r.URL.Query()

	queryToken := query.Get(tokenName)
	if queryToken != "" {
		return queryToken, nil
	}

	tokenHeader := r.Header.Get(tokenName)
	if tokenHeader != "" {
		return tokenHeader, nil
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		tokenHeader := strings.Split(authHeader, " ")[1]
		return tokenHeader, nil
	}

	tokenCookie, err := r.Cookie(tokenName)
	if err == nil {
		return tokenCookie.Value, nil
	}

	return "", errors.New("No token found")
}
