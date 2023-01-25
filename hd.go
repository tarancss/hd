// Package hd provides hierarchical deterministic wallet (HD wallet) functionality according to BIP39, BIP32 and BIP44.
// The initialization of the wallet requires a 64-byte seed. It is recommended to generate seeds using BIP39 out of a
// 24 word mnemonic and passphrase which are easy to remember.
// Once the HdWallet is initialized, you can easily generate any address.
// For a full description of what a HD wallet is, please read: https://en.bitcoinwiki.org/wiki/Deterministic_wallet
package hd

import (
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// External addresses are used for inputs or deposits.
	External uint8 = 0x00
	// Change addresses are used for internal use.
	Change uint8 = 0x01

	purpose  uint32 = 44 // BIP44
	coin     uint32 = 60 // Ethereum
	hardened uint32 = 0x80000000
)

var (
	// ErrInternal wraps errors reported by deps.
	ErrInternal error = errors.New("hd internal error")
	// ErrInvalidSeedLen will be reported when setting a wrong length (recommended seed length is 64-byte).
	ErrInvalidSeedLen error = errors.New("hd: length of seed is invalid")
	// ErrUnusableSeed will be reported if the seed cannot be used.
	ErrUnusableSeed error = errors.New("hd: the master key cannot be used")
)

// HdWallet is a composed type.
type HdWallet struct { //nolint:golint // changing would break compatibility
	*hdkeychain.ExtendedKey // HD wallet branch from which account/addresses are generated
}

// Init initializes the HD wallet for Ethereum for the given seed.
func Init(seed []byte) (*HdWallet, error) {
	// generate a master wallet
	master, err := getHdMaster(seed)
	if err != nil {
		return nil, err
	}
	// generate a BIP44 and Ethereum branch
	tmpW, err := master.Derive(hdkeychain.HardenedKeyStart + purpose)
	if err != nil {
		return nil, fmt.Errorf("%s: %w ", ErrInternal, err)
	}

	tmpW, err = tmpW.Derive(hdkeychain.HardenedKeyStart + coin)
	if err != nil {
		return nil, fmt.Errorf("%s: %w ", ErrInternal, err)
	}

	return &HdWallet{ExtendedKey: tmpW}, nil
}

// Address generates an address for 'wallet', flg should be either external or change and address number.
func (w *HdWallet) Address(wallet uint32, flg uint8, addrNum uint32,
) (addr, key []byte, prv ecdsa.PrivateKey, err error) {
	var tmpW *hdkeychain.ExtendedKey
	// get account
	tmpW, err = w.Derive(hdkeychain.HardenedKeyStart + wallet)
	if err != nil {
		return
	}
	// get external
	tmpW, err = tmpW.Derive(uint32(flg & Change))
	if err != nil {
		return
	}
	// get index to be used as address
	tmpW, err = tmpW.Derive(hdkeychain.HardenedKeyStart + addrNum)
	if err != nil {
		return
	}

	privateKey, _ := tmpW.ECPrivKey()
	prv = *privateKey.ToECDSA()

	return crypto.PubkeyToAddress(prv.PublicKey).Bytes(), crypto.FromECDSA(&prv), prv, nil
}

// getHdMaster generates a Hd master wallet that can be used for many coins.
func getHdMaster(seed []byte) (*hdkeychain.ExtendedKey, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if len(seed) < hdkeychain.MinSeedBytes || len(seed) > hdkeychain.MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "Bitcoin seed", Data = S)
	hmac512 := hmac.New(sha512.New, []byte("Bitcoin seed"))
	if _, err := hmac512.Write(seed); err != nil {
		return nil, fmt.Errorf("%s: %w ", ErrInternal, err)
	}

	lr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences where left is master secret key and right is master chain code
	secretKey := lr[:len(lr)/2]
	chainCode := lr[len(lr)/2:]

	// Ensure the key is usable.
	secretKeyNum := new(big.Int).SetBytes(secretKey)
	if secretKeyNum.Cmp(btcec.S256().N) >= 0 || secretKeyNum.Sign() == 0 {
		return nil, ErrUnusableSeed
	}

	parentFP := []byte{0x00, 0x00, 0x00, 0x00}

	return hdkeychain.NewExtendedKey([]byte{0x04, 0x88, 0xad, 0xe4}, secretKey, chainCode,
		parentFP, 0, 0, true), nil // starts with xprv -> {0x04, 0x88, 0xad, 0xe4}
}
