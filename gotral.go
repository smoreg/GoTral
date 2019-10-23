/*
* github.com/codenoid - Developer
* code source : - https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/
*  				- https://golang.org/pkg/crypto/cipher
* 				- https://github.com/codenoid/GoTral
*
 */
package gotral

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// createHash : create md5 hash and return as string
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt : encrypt given data with passphrase
// Load your secret key from a safe place and reuse it across multiple
// Seal/Open calls.
func Encrypt(data []byte, passphrase string) []byte {
	// create aes.NewCipher from hashed md5 passphrase
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	//  NewGCM returns the given 128-bit, block cipher wrapped in
	// Galois Counter Mode with the standard nonce length.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	// initialize slice with length of nonce that must be passed to Seal and Open.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

// Decrypt : decrypt given data with passphrase
// Load your secret key from a safe place and reuse it across multiple
// Seal/Open calls.
func Decrypt(data []byte, passphrase string) []byte {
	// create md5 byte slice
	key := []byte(createHash(passphrase))
	// just `reverse` algoritm with passphrase until return
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}