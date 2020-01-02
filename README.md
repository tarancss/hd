# hd
Hierarchical deterministic wallet for Ethereum.

[![GoDoc](https://godoc.org/github.com/tarancss/hd?status.svg)](https://godoc.org/github.com/tarancss/hd)
[![Go Report Card](https://goreportcard.com/badge/github.com/tarancss/hd)](https://goreportcard.com/report/github.com/tarancss/hd)

#### Usage
This package provides hierarchical deterministic wallet ("HD wallet") functionality according to BIP39, BIP32 and BIP44.

Once the HdWallet is initialized, you can easily generate any address by requesting the wallet number, either Change or External and the id of the address (a number between 1 and 2e32-1). See test file for same code.

For a full description of what a HD wallet is, please read [here](https://en.bitcoinwiki.org/wiki/Deterministic_wallet).

#### Configuration
The initialization of the wallet requires a 64-byte seed. It is recommended to generate seeds using BIP39 out of a 24 word mnemonic and passphrase which are easy to remember. You should always keep private keys and seed safe.

