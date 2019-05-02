package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

//Bytes helps generate n random bytes
//Used for remember tokens
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//NBytes returns the number of bytes used in the base64 URL encoded string
func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}

//String generates a byts slice of size n
//returns a base64 URL encoded version of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

//RememberToken is a helper func to return tokens
//of a specified byte size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
