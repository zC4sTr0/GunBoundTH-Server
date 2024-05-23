package cryptography

import (
	"crypto/aes"
	"encoding/hex"
)

const (
	KEY_BROKER = "FFB3B3BEAE97AD83B9610E23A43C2EB0"
)

// StringDecode decodes a byte array to a string, stopping at the first null byte
func StringDecode(input []byte) string {
	result := ""
	for _, b := range input {
		if b != 0 {
			result += string(b)
		} else {
			break
		}
	}
	return result
}

// AESDecryptBlock decrypts a block using AES in ECB mode
func AESDecryptBlock(block []byte, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(block))
	cipher.Decrypt(decrypted, block)
	return decrypted, nil
}

// GunboundStaticDecrypt decrypts a block using either KEY_BROKER based on the encryption type
func GunboundStaticDecrypt(block []byte) ([]byte, error) {
	var key []byte
	var err error

	key, err = hex.DecodeString(KEY_BROKER)
	if err != nil {
		return nil, err
	}

	return AESDecryptBlock(block, key)
}
