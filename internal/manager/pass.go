package manager

var Pass = &pass{}

type pass struct {
	Dir string // source directory
}

func (p *pass) Password(hash []byte) ([]byte, error) {
	return []byte{}, nil
}

func (p *pass) Passwords() ([][]byte, error) {
	return [][]byte{}, nil
}

func (p *pass) Write(data []byte) error {
	return nil
}

func (p *pass) Delete(hash []byte) error {
	return nil
}
