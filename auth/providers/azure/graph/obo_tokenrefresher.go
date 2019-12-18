package graph

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/glog"
	"github.com/moul/http2curl"
	"github.com/pkg/errors"
)

type oboTokenRefresher struct {
	name         string
	client       *http.Client
	clientID     string
	clientSecret string
	scope        string
	loginURL     string
}

// NewOBOTokenRefresher returns a TokenRefresher that implements OAuth On-Behalf-Of flow on Azure Active Directory
// https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-on-behalf-of-flow
func NewOBOTokenRefresher(clientID, clientSecret, loginURL, scope string) TokenRefresher {
	return &oboTokenRefresher{
		name:         "OBOTokenRefresher",
		client:       &http.Client{},
		clientID:     clientID,
		clientSecret: clientSecret,
		scope:        scope,
		loginURL:     loginURL,
	}
}

func (u *oboTokenRefresher) Name() string { return u.name }

func (u *oboTokenRefresher) Refresh(token string) (AuthResponse, error) {
	var authResp = AuthResponse{}
	form := url.Values{}
	form.Set("client_id", u.clientID)
	form.Set("client_secret", u.clientSecret)
	form.Set("assertion", token)
	form.Set("requested_token_use", "on_behalf_of")
	form.Set("scope", u.scope)
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")

	req, err := http.NewRequest(http.MethodPost, u.loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return authResp, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if glog.V(10) {
		cmd, _ := http2curl.GetCurlCommand(req)
		glog.V(10).Infoln(cmd)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return authResp, errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		return authResp, errors.Errorf("request %s failed with status code: %d and response: %s", req.URL.Path, resp.StatusCode, string(data))
	}
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return authResp, errors.Wrapf(err, "failed to decode response for request %s", req.URL.Path)
	}

	return authResp, nil
}
