package cryptography

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
)

const (
	KEY_LAUNCHER = ""
	KEY_BROKER   = ""
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

// GunboundStaticDecrypt decrypts a block using either KEY_LAUNCHER or KEY_BROKER based on the encryption type
func GunboundStaticDecrypt(block []byte, encryptionType int) ([]byte, error) {
	var key []byte
	var err error

	switch encryptionType {
	case 1:
		key, err = hex.DecodeString(KEY_LAUNCHER)
	case 2:
		key, err = hex.DecodeString(KEY_BROKER)
	default:
		return nil, fmt.Errorf("invalid encryption type")
	}

	if err != nil {
		return nil, err
	}

	return AESDecryptBlock(block, key)
}
