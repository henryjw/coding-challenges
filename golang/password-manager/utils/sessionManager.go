package utils

var session *Session

type Session struct {
	VaultId           uint
	VaultPasswordHash string
}

// GetSession Returns `nil` if no session exists
func GetSession() *Session {
	if session == nil {
		return nil
	}

	return &Session{
		VaultId:           session.VaultId,
		VaultPasswordHash: session.VaultPasswordHash,
	}
}

func ResetSession() {
	session = nil
}

func CreateSession(vaultId uint, vaultPasswordHash string) Session {
	session = &Session{
		VaultId:           vaultId,
		VaultPasswordHash: vaultPasswordHash,
	}

	return *session
}
