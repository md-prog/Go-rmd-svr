package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type config struct {
	Key string
}

//IntercomService TODO
type IntercomService struct {
	Key string
}

//CalculateHash TODO
func (i *IntercomService) CalculateHash(user_id string) string {
	key := []byte(i.Key)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(user_id))
	return hex.EncodeToString(h.Sum(nil))
}

//InitIntercomService TODO
func InitIntercomService(key string) *IntercomService {
	return &IntercomService{Key: key}
}
