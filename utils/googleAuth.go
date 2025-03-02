package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuth struct {
	Name    string `json:"name"`
	SurName string `json:"surName"`
	Email   string `json:"email"`
}

var googleOauthConfig = oauth2.Config{
	ClientID:     "953257087625-dm82cn9b20a19526g33cmu1di1q34nju.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:3000",
	Endpoint:     google.Endpoint,
	Scopes:       []string{"profile", "email"},
}

func GeteCredentialsByCode(req *http.Request, code string) (GoogleAuth, error) {

	token, err := googleOauthConfig.Exchange(req.Context(), code)

	if err != nil {
		return GoogleAuth{}, err
	}

	client := googleOauthConfig.Client(req.Context(), token)

	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")

	if err != nil {
		return GoogleAuth{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return GoogleAuth{}, err
	}

	var credentials GoogleAuth
	err = json.Unmarshal(body, &credentials)

	return credentials, err

}
