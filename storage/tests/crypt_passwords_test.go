package tests

import (
	"github.com/stretchr/testify/require"
	"storage/internal/cryptPasswords"
	"testing"
)

func TestCreateAndCompare(t *testing.T) {
	hash, err := cryptPasswords.GeneratePasswordHash("password")
	require.NoError(t, err)

	err = cryptPasswords.ComparePasswordWithHash(hash, "password")
	require.NoError(t, err)

	err = cryptPasswords.ComparePasswordWithHash(hash, "password1")
	require.Error(t, err)

	hash, err = cryptPasswords.GeneratePasswordHash("passwordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpasswordpassword")
	require.Error(t, err)
}
