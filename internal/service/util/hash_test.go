package hash_test

import (
	hash "art-sso/internal/service/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	hashedPassword, err := hash.HashPassword(password)
	assert.Nil(t, err)

	verified := hash.VerifyPassword(password, hashedPassword)
	assert.True(t, verified)
}
