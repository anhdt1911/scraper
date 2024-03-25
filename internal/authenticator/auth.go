package authenticator

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"

	"github.com/anhdt1911/scraper/internal/config"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func New() (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+config.AuthDomain+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     config.AuthClientID,
		ClientSecret: config.AuthClientSecret,
		RedirectURL:  config.AuthCallbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

func (a *Authenticator) VerifyIDToken(c context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oath2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(c, rawIDToken)
}

func (a *Authenticator) Login(c *gin.Context) {
	// To prevent CSRF Attacks.
	state, err := generateRandomState()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	session := sessions.Default(c)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		c.String(500, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, a.AuthCodeURL(state))
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func (a *Authenticator) Callback(c *gin.Context) {
	session := sessions.Default(c)
	if c.Query("state") != session.Get("state") {
		c.String(http.StatusBadRequest, "Invalid state parameter")
		return
	}

	token, err := a.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
		return
	}

	idToken, err := a.VerifyIDToken(c.Request.Context(), token)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to verify ID Token")
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)
	if err := session.Save(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/dashboard")
}

func (a *Authenticator) GetUser(c *gin.Context) {
	session := sessions.Default(c)
	profile := session.Get("profile")

	c.JSON(http.StatusOK, gin.H{"data": profile})
}

func (a *Authenticator) Logout(c *gin.Context) {
	logoutUrl, err := url.Parse("https://" + config.AuthDomain + "/v2/logout")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + c.Request.Host)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", config.AuthClientID)
	logoutUrl.RawQuery = parameters.Encode()

	c.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
