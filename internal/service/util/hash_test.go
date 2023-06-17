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

	err = hash.VerifyPassword(password, hashedPassword)
	assert.Nil(t, err)
}
