// Package hd
// This is the testing of functions of the HD wallet
// Tests are run and compared against values provided by https://iancoleman.io/bip39/
package hd

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestHdWallet(t *testing.T) {
	var seed []byte
	seed, _ = hex.DecodeString("642ce4e20f09c9f4d285c2b336063eaafbe4cb06dece8134f3a64bdd8f8c0c24df73e1a2e7056359b6db61e179ff45e5ada51d14f07b30becb6d92b961d35df4") // seed has been generated with mnemonic "tuna song credit master earn feature dutch nurse yellow ship caution relief ten drip trip couch increase nominee salt drift nation oval exhaust baby" and passphrase "password"

	w, err := Init(seed)
	if err != nil {
		t.Errorf("Init %e", err)
	}

	// We generate 3 addresses for wallet 2, hd.External, indices from 0 to 2.
	var expected [][]string = [][]string{
		{"0xD43E2870777916Ede1f5Cc43F14f8C0741e11f96", "0x735e6eec7fbd869aafa61e50921b101eebc1d6961b8019a76bcf27cade1304b7"},
		{"0xF4cEFC8d1AfaA51d5A5E7f57d214B60429cA4378", "0xfa7d6a67439ec17e07c10f10a4a9007e46583b5219cb909c8b474398b7216917"},
		{"0x8A1847459c5FCD66f0B29012a21A2D5A314Ef1D0", "0x99c59090d814b1a3eb2b1f1715e7e7a09cbf8a770495a9a7836fb02226f8fd44"},
	}
	var addr, key, addrExp, keyExp []byte
	for i := uint32(0); i < uint32(3); i++ {
		addr, key, _, err = w.Address(uint32(2), External, i)
		if err != nil {
			t.Errorf("Address %d :%e", i, err)
		}
		addrExp, _ = hex.DecodeString(expected[i][0][2:])
		if bytes.Compare(addr, addrExp) != 0 {
			t.Errorf("Address %d does not match. Got:%x, expected:%x", i, addr, addrExp)
		}
		keyExp, _ = hex.DecodeString(expected[i][1][2:])
		if bytes.Compare(addr, addrExp) != 0 {
			t.Errorf("Key %d does not match. Got:%x, expected:%x", i, key, keyExp)
		}
	}
	return
}
