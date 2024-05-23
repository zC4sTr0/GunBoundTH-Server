package main

import (
	"encoding/hex"
	"fmt"

	"github.com/zC4sTr0/GunBoundTH-Server/cryptography"
)

func main() {
	data, err := hex.DecodeString("E4AE6422B374D8779DC2F3695810650F")
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return
	}

	decrypted, err := cryptography.GunboundStaticDecrypt(data)
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return
	}

	fmt.Println(cryptography.StringDecode(decrypted))
}
