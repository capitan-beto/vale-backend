package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func VerifySignature(hash string, manifest string) error {
	secret := os.Getenv("WEBHOOK_SECRET")
	hmac := hmac.New(sha256.New, []byte(secret))
	hmac.Write([]byte(manifest))

	sha := hex.EncodeToString(hmac.Sum(nil))

	if sha == hash {
		fmt.Println("HMAC verification passed")
		return nil
	} else {
		return fmt.Errorf("HMAC verification failed")
	}
}
