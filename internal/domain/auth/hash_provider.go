package auth

type HashSalt struct {
	Hash, Salt  []byte
	HashEncoded string
}

type HashProvider interface {
	Generate(password string) (*HashSalt, error)
	Compare(password, hash string) error
}
