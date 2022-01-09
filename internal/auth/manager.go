package auth

type Manager interface {
	Key() ([]byte, error)
	SetKey([]byte) error

	Checksum() ([]byte, error)
	SetChecksum([]byte) error

	Salt() ([]byte, error)
	SetSalt([]byte) error
}
