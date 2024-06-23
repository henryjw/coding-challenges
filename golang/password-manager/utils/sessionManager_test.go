package utils_test

import (
	"github.com/stretchr/testify/assert"
	"password-manager/utils"
	"testing"
)

func resetState() {
	utils.ResetSession()
}

func TestCreateAndGetSession(t *testing.T) {
	t.Cleanup(resetState)

	vaultId := uint(1)
	vaultPasswordHash := "passwordHash"

	utils.CreateSession(vaultId, vaultPasswordHash)

	session := utils.GetSession()

	assert.Equal(t, vaultId, session.VaultId)
	assert.Equal(t, vaultPasswordHash, session.VaultPasswordHash)
}

func TestGetSessionWithoutCreating(t *testing.T) {
	t.Cleanup(resetState)

	session := utils.GetSession()
	assert.Nil(t, session)
}
