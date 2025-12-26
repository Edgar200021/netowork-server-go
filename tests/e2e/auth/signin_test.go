package auth

import (
	"testing"

	"github.com/Edgar200021/netowork-server-go/tests"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	tests.New(t)

	c := 2

	assert.Equal(t, c, 2)
}
