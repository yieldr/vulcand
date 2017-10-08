package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrowdAuth(t *testing.T) {
	crowd := NewCrowdProvider("https://crowd.ydworld.com", "impulse", "rjX7f9l44Y8W46T")
	user, err := crowd.Authenticate("chucknorris", "infinity")
	assert.NoError(t, err)
	assert.Equal(t, "chucknorris", user.Username())
	assert.Equal(t, "Chuck Norris", user.FullName())
	assert.Equal(t, "cn@ydworld.com", user.Email())
}

func TestCrowdAuthFail(t *testing.T) {
	crowd := NewCrowdProvider("https://crowd.ydworld.com", "impulse", "rjX7f9l44Y8W46T")
	user, err := crowd.Authenticate("troll", "i can haz access?")
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Invalid username or password.")
}
