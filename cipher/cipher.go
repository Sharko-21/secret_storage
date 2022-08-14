package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"main/logger"
)

const defaultKey = "cRfUjXn2r4u7x!A%D*G-KaPdSgVkYp3s"

func Encrypt(text []byte, userKey []byte) []byte {
	key := userKey
	if len(key) == 0 {
		key = []byte(defaultKey)
	}
	cphr, err := aes.NewCipher(key)
	if err != nil {
		logger.ErrorLog.Printf("%+v", err)
	}
	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		logger.ErrorLog.Printf("%+v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.ErrorLog.Printf("%+v", err)
	}
	return gcm.Seal(nonce, nonce, text, nil)
}

func Decrypt(cipherText []byte, userKey []byte) []byte {
	key := userKey
	if len(key) == 0 {
		key = []byte(defaultKey)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		logger.ErrorLog.Printf("%+v", err)
	}
	gcmDecrypt, err := cipher.NewGCM(c)
	if err != nil {
		logger.ErrorLog.Printf("%+v", err)
	}
	nonceSize := gcmDecrypt.NonceSize()
	if len(cipherText) < nonceSize {
		logger.ErrorLog.Printf("%+v", err)
	}
	nonce, encryptedMessage := cipherText[:nonceSize], cipherText[nonceSize:]
	plaintext, _ := gcmDecrypt.Open(nil, nonce, encryptedMessage, nil)
	return plaintext
}
