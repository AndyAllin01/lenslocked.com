package hash

import (
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha256"
	"hash"
)

//NewHMAC creates and returns a new hmac object
func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{
		hmac: h,
	}

}

//HMAC wraps around the crypto/hmac package
type HMAC struct {
	hmac hash.Hash
}

//Hash hashes the input string using HMAC with the 
//secret key provided when the HMAC object was created
func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b:=h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}