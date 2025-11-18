package shortener

import (
    "crypto/rand"
    "math/big"
)

type Shortener struct {
    length int
}

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_.~"

// New creates a new shortener that generates slugs of the given length
func New(length int) *Shortener {
    return &Shortener{length: length}
}

// Generate produces a random slug with configured length.
// It uses a larger alphabet to produce longer-looking slugs (URL-safe characters).
func (s *Shortener) Generate() (string, error) {
    bytes := make([]byte, s.length)
    alphaLen := big.NewInt(int64(len(alphabet)))
    for i := 0; i < s.length; i++ {
        n, err := rand.Int(rand.Reader, alphaLen)
        if err != nil {
            return "", err
        }
        bytes[i] = alphabet[n.Int64()]
    }
    return string(bytes), nil
}
