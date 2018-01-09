package main

import (
	"testing"

	"github.com/yieldr/vulcand/plugin/oauth2"
	"github.com/yieldr/vulcand/registry"
)

func TestFromJSON(t *testing.T) {
	r, err := registry.GetRegistry()
	if err != nil {
		t.Fatal(err)
	}

	spec := r.GetSpec("oauth2")
	m, err := spec.FromJSON([]byte(`{
  "IssuerURL": "https://auth.example.com",
  "ClientID": "my-client",
  "ClientSecret": "1234567890",
  "RedirectURL": "https://example.com/foo/bar"
}`))
	if err != nil {
		t.Fatal(err)
	}

	o, ok := m.(*oauth2.OAuth2)
	if !ok {
		t.Fatalf("unexpected middleware struct %T", m)
	}

	assertEqual(t, o.IssuerURL, "https://auth.example.com")
	assertEqual(t, o.ClientID, "my-client")
	assertEqual(t, o.ClientSecret, "1234567890")
	assertEqual(t, o.RedirectURL, "https://example.com/foo/bar")
	assertEqual(t, o.RedirectURLPath, "/foo/bar")
}

func assertEqual(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("failed asserting %q equals %q", a, b)
	}
}
