package oauth2

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/alexkappa/errors"
	"github.com/codegangsta/cli"
	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/vulcand/vulcand/plugin"
	"golang.org/x/oauth2"
)

var SessionStore sessions.Store

func init() {
	SessionStore = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
}

func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      "oauth2",
		FromOther: FromOther,
		FromCli:   FromCli,
		CliFlags:  CliFlags(),
	}
}

func FromOther(o OAuth2) (plugin.Middleware, error) {
	return New(o.Domain, o.ClientID, o.ClientSecret, o.RedirectURL)
}

func FromCli(c *cli.Context) (plugin.Middleware, error) {
	return New(
		c.String("domain"),
		c.String("clientId"),
		c.String("clientSecret"),
		c.String("redirectUrl"))
}

func CliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "domain", Usage: "oauth2 idp domain"},
		cli.StringFlag{Name: "clientId", Usage: "oauth2 client id"},
		cli.StringFlag{Name: "clientSecret", Usage: "oauth2 client secret"},
		cli.StringFlag{Name: "redirectUrl", Usage: "oauth2 redirect url"},
	}
}

type OAuth2 struct {
	// Domain holds the authorization servers domain. This can be any OAuth2
	// compatible server however this plugin has only been tested to work with
	// Auth0.
	Domain string
	// RedirectURLPath holds the URL path of the redirect URL. This path will be
	// reserved for handling the OAuth2 redirect callback therefore it is
	// important to use a path that does not conflict with upstream services
	// routing.
	RedirectURLPath string

	*oauth2.Config
}

func New(domain, clientID, clientSecret, redirectURL string) (*OAuth2, error) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid redirectUrl")
	}
	return &OAuth2{
		Domain:          domain,
		RedirectURLPath: u.Path,
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  u.String(),
			Scopes:       []string{"openid", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://" + domain + "/authorize",
				TokenURL: "https://" + domain + "/oauth/token",
			},
		},
	}, nil
}

func (o *OAuth2) NewHandler(next http.Handler) (http.Handler, error) {
	return context.ClearHandler(&OAuth2Handler{
		oauth2: o,
		next:   next,
	}), nil
}

func (o *OAuth2) String() string {
	return fmt.Sprintf(
		"domain=%s, client-id=%s client-secret=%s redirect=%s",
		o.Domain,
		o.ClientID,
		o.ClientSecret,
		o.RedirectURL,
	)
}

type OAuth2Handler struct {
	oauth2 *OAuth2
	next   http.Handler
}

func (h *OAuth2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == h.oauth2.RedirectURLPath {
		h.Callback(w, r)
	} else {
		h.All(w, r)
	}
}

func (h *OAuth2Handler) Callback(w http.ResponseWriter, r *http.Request) {

	t, err := h.oauth2.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s, err := SessionStore.Get(r, "oauth2-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !t.Valid() {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	s.Values["access_token"] = t.AccessToken
	s.Values["id_token"] = t.Extra("id_token")

	if err = s.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnTo := "/"
	if _, ok := s.Values["return_to"]; ok {
		returnTo = s.Values["return_to"].(string)
	}

	http.Redirect(w, r, returnTo, http.StatusTemporaryRedirect)
}

func (h *OAuth2Handler) All(w http.ResponseWriter, r *http.Request) {

	s, err := SessionStore.Get(r, "oauth2-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := s.Values["access_token"]; !ok {

		s.Values["return_to"] = r.URL.Path

		if err = s.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		q := make(url.Values)
		q.Set("client", h.oauth2.ClientID)
		q.Set("redirect_uri", h.oauth2.RedirectURL)
		q.Set("protocol", "oauth2")
		q.Set("response_type", "code")

		http.Redirect(w, r, "https://"+h.oauth2.Domain+"/login?"+q.Encode(), 302)
		return
	}

	h.next.ServeHTTP(w, r)
}
