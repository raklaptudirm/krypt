package pass

type Manager interface {
	Password([]byte) ([]byte, error)
	Passwords() ([][]byte, error)
	Write([]byte) error
	Delete([]byte) error
}
