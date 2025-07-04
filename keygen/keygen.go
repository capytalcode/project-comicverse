package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

func main() {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Private Key: %x\n", priv)
	fmt.Printf("Public Key: %x\n", pub)

	pub64 := base64.URLEncoding.EncodeToString(pub)
	priv64 := base64.URLEncoding.EncodeToString(priv)

	fmt.Printf("Private Key: %s\n", priv64)
	fmt.Printf("Public Key: %s\n", pub64)
}
