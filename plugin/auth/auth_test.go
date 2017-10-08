package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockProvider struct {
	user     User
	password string
}

func (m mockProvider) Authenticate(username, password string) (User, error) {
	if username == m.user.Username() && password == m.password {
		return m.user, nil
	}
	return nil, errors.New("Authentication failed")
}

func NewMockProvider(user User, password string) Provider {
	return mockProvider{user: user, password: password}
}

type mockCache struct{}

func (m *mockCache) Set(key string, value interface{}, duration time.Duration) {}

func (m *mockCache) Get(key string) (interface{}, bool) { return nil, false }

func NewMockCache() *mockCache {
	return new(mockCache)
}

func TestAuthHandler(t *testing.T) {
	h := &AuthHandler{
		NewMockProvider(
			NewMockUser(
				"johndoe",
				"John Doe",
				"john.doe@yieldr.com",
				[]string{"acc1", "acc2"},
				[]string{"role1", "role1"},
			),
			"s3cretp4ssw0rd",
		),
		NewMockCache(),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "johndoe", r.Header.Get("X-Auth-Username"))
			assert.Equal(t, "John Doe", r.Header.Get("X-Auth-FullName"))
			assert.Equal(t, "john.doe@yieldr.com", r.Header.Get("X-Auth-Email"))
			assert.Equal(t, "acc1,acc2", r.Header.Get("X-Auth-Accounts"))
			assert.Equal(t, "role1,role1", r.Header.Get("X-Auth-Roles"))
		}),
	}

	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.SetBasicAuth("johndoe", "s3cretp4ssw0rd")

	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandlerFail(t *testing.T) {
	h := &AuthHandler{
		NewMockProvider(NewMockUser("johndoe", "", "", nil, nil), "s3cretp4ssw0rd"),
		NewMockCache(),
		http.NotFoundHandler(),
	}

	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	r.SetBasicAuth("troll", "i can haz access?")

	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMask(t *testing.T) {
	for _, in := range []string{
		"much password",
		"such sensitive",
		"wow!",
	} {
		out := Mask(in, '*')
		assert.Len(t, out, len(in))
		for _, r := range out {
			assert.Equal(t, '*', r)
		}
	}
}

func TestBase64(t *testing.T) {
	for s, expected := range map[string]string{
		"foo:bar": "Zm9vOmJhcg==",
		"bar:baz": "YmFyOmJheg==",
		"baz:foo": "YmF6OmZvbw==",
	} {
		assert.Equal(t, expected, Base64(s))
	}
}
