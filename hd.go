// Package hd provides hierarchical deterministic wallet ("HD wallet") functionality according to BIP39, BIP32 and BIP44.
// The initialization of the wallet requires a 64-byte seed. It is recommended to generate seeds using BIP39 out of a 24 word mnemonic and passphrase which are easy to remember.
// Once the HdWallet is initialized, you can easily generate any address.
// For a full description of what a HD wallet is, please read here: https://en.bitcoinwiki.org/wiki/Deterministic_wallet
package hd

import (
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	External uint8  = 0x00 // Use this flag to generate addresses used for external inputs or deposits.
	Change   uint8  = 0x01 // Internal chain or change addresses.
	purpose  uint32 = 44   // BIP44
	coin     uint32 = 60   // Ethereum
	hardened uint32 = 0x80000000
)

var (
	ErrInvalidSeedLen error = errors.New("ErrInvalidSeedLen: length of seed is invalid")
	ErrUnusableSeed   error = errors.New("ErrUnusableSeed: the master key cannot be used")
)

// HdWallet is a composed type
type HdWallet struct {
	*hdkeychain.ExtendedKey // HD wallet branch from which account/addresses are generated
}

// Init initializes the HD wallet for Ethereum for the given seed.
func Init(seed []byte) (*HdWallet, error) {
	// generate a master wallet
	master, err := getHdMaster(seed)
	if err != nil {
		return &HdWallet{}, err
	}
	// generate a BIP44 and Ethereum branch
	tmpW, err := master.Child(hardened + purpose)
	if err != nil {
		return nil, err
	}
	tmpW, err = tmpW.Child(hardened + coin)
	return &HdWallet{tmpW}, err
}

// Address generates an address for 'wallet', flg should be either external or change and address number.
func (w *HdWallet) Address(wallet uint32, flg uint8, addrNum uint32) (addr, key []byte, prv ecdsa.PrivateKey, err error) {
	var tmpW *hdkeychain.ExtendedKey
	// get account
	tmpW, err = w.Child(hardened + wallet)
	if err != nil {
		return
	}
	// get external
	tmpW, err = tmpW.Child(uint32(flg & Change))
	if err != nil {
		return
	}
	// get index to be used as address
	tmpW, err = tmpW.Child(hardened + addrNum)
	if err != nil {
		return
	}
	private_key, _ := tmpW.ECPrivKey()
	prv = ecdsa.PrivateKey(*private_key)
	return crypto.PubkeyToAddress(private_key.PublicKey).Bytes(), crypto.FromECDSA(&prv), prv, nil
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
	hmac512.Write(seed)
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
