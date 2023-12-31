package handler

import (
	"github.com/Hari-Kiri/goalRsa"
)

var (
	rsaPrivateKey string
	rsaPublicKey  string
)

func init() {
	var errorGenerateRsaKeyPair error
	rsaPrivateKey, rsaPublicKey, errorGenerateRsaKeyPair = goalRsa.NewPemFormatRsaKeyPair(4096)
	if errorGenerateRsaKeyPair != nil {
		panic(errorGenerateRsaKeyPair)
	}
}
