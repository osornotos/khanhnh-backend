package hashing

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"khanhnh-backend/internal/domain/auth"
	"strings"
)

var _ auth.HashProvider = (*Argon2idHash)(nil)

var (
	// ErrInvalidHash in returned by ComparePasswordAndHash if the provided
	// hash isn't in the expected format.
	ErrInvalidHash = errors.New("argon2id: hash is not in the correct format")

	// ErrIncompatibleVariant is returned by ComparePasswordAndHash if the
	// provided hash was created using a unsupported variant of Argon2.
	// Currently only argon2id is supported by this package.
	ErrIncompatibleVariant = errors.New("argon2id: incompatible variant of argon2")

	// ErrIncompatibleVersion is returned by ComparePasswordAndHash if the
	// provided hash was created using a different version of Argon2.
	ErrIncompatibleVersion = errors.New("argon2id: incompatible version of argon2")

	ErrHashNotMatch = errors.New("hash doesn't match")
)

type Argon2idHash struct {
	// time represents the number of
	// passed over the specified memory.
	time uint32
	// cpu memory to be used.
	memory uint32
	// threads for parallelism aspect
	// of the algorithm.
	threads uint8
	// keyLen of the generate hash key.
	keyLen uint32
	// saltLen the length of the salt used.
	saltLen uint32
}

// NewArgon2idHash constructor function for
// Argon2idHash.
func NewArgon2idHash(time, saltLen uint32, memory uint32, threads uint8, keyLen uint32) *Argon2idHash {
	if threads < 1 { // panic
		threads = 1
	}
	return &Argon2idHash{
		time:    time,
		saltLen: saltLen,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
}

// Generate returns an auth.HashSalt hash of a plain-text password using the
// provided algorithm parameters. The returned hash follows the format used by
// the Argon2 reference C implementation and contains the base64-encoded Argon2id
// derived key prefixed by the salt and parameters. It looks like this:
//
//	$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
func (a *Argon2idHash) Generate(password string) (*auth.HashSalt, error) {
	hashSalt, err := a.createHash([]byte(password), nil)
	if err != nil {
		return nil, err
	}
	hashSalt.HashEncoded = a.combineHashWithSalt(hashSalt.Hash, hashSalt.Salt)
	return hashSalt, nil
}

// GenerateWithSalt same as Generate but fixed salt
func (a *Argon2idHash) GenerateWithSalt(password string, salt []byte) (*auth.HashSalt, error) {
	hashSalt, err := a.createHash([]byte(password), salt)
	if err != nil {
		return nil, err
	}
	hashSalt.HashEncoded = a.combineHashWithSalt(hashSalt.Hash, hashSalt.Salt)
	return hashSalt, nil
}

// Compare generated hash with store hash.
func (a *Argon2idHash) Compare(password, encodedHash string) error {
	salt, hashFromEncodedHash, err := a.decodeHash(encodedHash)
	if err != nil {
		return err
	}

	// Generate hash for comparison.
	hashFromPassword, err := a.createHash([]byte(password), salt)
	if err != nil {
		return err
	}
	// Compare the generated hash with the stored hash.
	if subtle.ConstantTimeEq(int32(len(hashFromEncodedHash)), int32(len(hashFromPassword.Hash))) == 0 {
		return ErrHashNotMatch
	}
	if subtle.ConstantTimeCompare(hashFromEncodedHash, hashFromPassword.Hash) == 1 {
		return nil
	}
	//if !bytes.Equal(hash, salt) {
	//	return errors.New("hash doesn't match")
	//}
	return ErrHashNotMatch
}

// createHash using the password and provided salt.
// If not salt value provided fallback to random value
// generated of a given length.
func (a *Argon2idHash) createHash(password, salt []byte) (*auth.HashSalt, error) {
	var err error
	// If salt is not provided generate a salt of
	// the configured salt length.
	if len(salt) == 0 {
		salt, err = randomSecret(a.saltLen)
	}
	if err != nil {
		return nil, err
	}
	// Generate hash
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)
	// Return the generated hash and salt used for storage.
	return &auth.HashSalt{Hash: hash, Salt: salt}, nil
}

// combineHashWithSalt return hash follows the format used by
// the Argon2 reference C implementation and contains the base64-encoded Argon2id
// derived key prefixed by the salt and parameters. It looks like this:
//
//	$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
func (a *Argon2idHash) combineHashWithSalt(hash, salt []byte) string {
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.memory, a.time, a.threads, b64Salt, b64Hash)
	return encoded
}

// decodeHash expects a hash created from this package, and parses it to return the params used to
// create it, as well as the salt and key (password hash).
func (a *Argon2idHash) decodeHash(encodedHash string) (salt, key []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, ErrInvalidHash
	}

	if vals[1] != "argon2id" {
		return nil, nil, ErrIncompatibleVariant
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, ErrIncompatibleVersion
	}

	var memoryFromHash, timeFromHash, threadsFromHash uint32
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memoryFromHash, &timeFromHash, &threadsFromHash)
	if err != nil {
		return nil, nil, err
	}
	if memoryFromHash != a.memory || timeFromHash != a.time || threadsFromHash != uint32(a.threads) {
		return nil, nil, ErrIncompatibleVariant
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, err
	}

	return salt, key, nil
}
