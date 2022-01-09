package manager

type auth struct {
	Dir string // source directory
}

func (a *auth) Key() ([]byte, error) {
	return []byte{}, nil
}

func (a *auth) SetKey(data []byte) error {
	return nil
}

func (a *auth) Checksum() ([]byte, error) {
	return []byte{}, nil
}

func (a *auth) SetChecksum(data []byte) error {
	return nil
}

func (a *auth) Salt() ([]byte, error) {
	return []byte{}, nil
}

func (a *auth) SetSalt(data []byte) error {
	return nil
}
