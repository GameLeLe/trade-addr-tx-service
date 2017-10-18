// Package hdwallet implements heirarchical deterministic Bitcoin wallets, as defined in BIP 32.
//
// BIP 32 - https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
//
// This package provides utilities for generating hierarchical deterministic Bitcoin wallets.
//
// Examples
//
//          // Generate a random 256 bit seed
//          seed, err := hdwallet.GenSeed(256)
//
//          // Create a master private key
//          masterprv := hdwallet.MasterKey(seed)
//
//          // Convert a private key to public key
//          masterpub := masterprv.Pub()
//
//          // Generate new child key based on private or public key
//          childprv, err := masterprv.Child(0)
//          childpub, err := masterpub.Child(0)
//
//          // Create bitcoin address from public key
//          address := childpub.Address()
//
//          // Convenience string -> string Child and Address functions
//          walletstring := childpub.String()
//          childstring, err := hdwallet.StringChild(walletstring,0)
//          childaddress, err := hdwallet.StringAddress(childstring)
//
// Extended Keys
//
// Hierarchical deterministic wallets are simply deserialized extended keys. Extended Keys can be imported and exported as base58-encoded strings. Here are two examples:
//          public key:   "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"
//          private key:  "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"
//
package hdwallet

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
)

var curve *btcec.KoblitzCurve = btcec.S256()

func hash160(data []byte) []byte {
	sha := sha256.New()
	ripe := ripemd160.New()
	sha.Write(data)
	ripe.Write(sha.Sum(nil))
	return ripe.Sum(nil)
}

func dblSha256(data []byte) []byte {
	sha1 := sha256.New()
	sha2 := sha256.New()
	sha1.Write(data)
	sha2.Write(sha1.Sum(nil))
	return sha2.Sum(nil)
}

func privToPub(key []byte) []byte {
	//return compress(curveETH.ScalarBaseMult(key))
	return compress(curve.ScalarBaseMult(key))
}

func onCurve(x, y *big.Int) bool {
	return curve.IsOnCurve(x, y)
}

func compress(x, y *big.Int) []byte {
	two := big.NewInt(2)
	rem := two.Mod(y, two).Uint64()
	rem += 2
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(rem))
	rest := x.Bytes()
	return append(b[1:], rest...)
}

//2.3.4 of SEC1 - http://www.secg.org/index.php?action=secg,docs_secg
func expand(key []byte) (*big.Int, *big.Int) {
	params := curve.Params()
	exp := big.NewInt(1)
	exp.Add(params.P, exp)
	exp.Div(exp, big.NewInt(4))
	x := big.NewInt(0).SetBytes(key[1:len(key)])
	y := big.NewInt(0).SetBytes(key[:1])
	beta := big.NewInt(0)
	beta.Exp(x, big.NewInt(3), nil)
	beta.Add(beta, big.NewInt(7))
	beta.Exp(beta, exp, params.P)
	if y.Add(beta, y).Mod(y, big.NewInt(2)).Int64() == 0 {
		y = beta
	} else {
		y = beta.Sub(params.P, beta)
	}
	return x, y
}

//Expand expand
func Expand(key []byte) (*big.Int, *big.Int) {
	params := curve.Params()
	exp := big.NewInt(1)
	exp.Add(params.P, exp)
	exp.Div(exp, big.NewInt(4))
	x := big.NewInt(0).SetBytes(key[1:len(key)])
	y := big.NewInt(0).SetBytes(key[:1])
	beta := big.NewInt(0)
	beta.Exp(x, big.NewInt(3), nil)
	beta.Add(beta, big.NewInt(7))
	beta.Exp(beta, exp, params.P)
	if y.Add(beta, y).Mod(y, big.NewInt(2)).Int64() == 0 {
		y = beta
	} else {
		y = beta.Sub(params.P, beta)
	}
	return x, y
}

func addPrivKeys(k1, k2 []byte) []byte {
	i1 := big.NewInt(0).SetBytes(k1)
	i2 := big.NewInt(0).SetBytes(k2)
	i1.Add(i1, i2)
	i1.Mod(i1, curve.Params().N)
	k := i1.Bytes()
	zero, _ := hex.DecodeString("00")
	return append(zero, k...)
}

func addPubKeys(k1, k2 []byte) []byte {
	x1, y1 := expand(k1)
	x2, y2 := expand(k2)
	return compress(curve.Add(x1, y1, x2, y2))
}

func uint32ToByte(i uint32) []byte {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, i)
	return a
}

func uint16ToByte(i uint16) []byte {
	a := make([]byte, 2)
	binary.BigEndian.PutUint16(a, i)
	return a[1:]
}

func byteToUint16(b []byte) uint16 {
	if len(b) == 1 {
		zero := make([]byte, 1)
		b = append(zero, b...)
	}
	return binary.BigEndian.Uint16(b)
}

//WalletToFile write the HDWallet bytes to file
func WalletToFile(filename string, wallet *HDWallet) error {
	var err error
	masterpub := wallet
	buf := new(bytes.Buffer)
	tmpBytes := make([]byte, 2)
	//write vBytes
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.Vbytes)))
	_, err = buf.Write(tmpBytes)
	if err != nil {
		return err
	}
	_, err = buf.Write(masterpub.Vbytes)
	if err != nil {
		return err
	}
	//write Fingerprint
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.Fingerprint)))
	_, err = buf.Write(tmpBytes)
	if err != nil {
		return err
	}
	_, err = buf.Write(masterpub.Fingerprint)
	if err != nil {
		return err
	}
	//write I
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.I)))
	_, err = buf.Write(tmpBytes)
	if err != nil {
		return err
	}
	_, err = buf.Write(masterpub.I)
	if err != nil {
		return err
	}
	//write Chaincode
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.Chaincode)))
	_, err = buf.Write(tmpBytes)
	if err != nil {
		return err
	}
	_, err = buf.Write(masterpub.Chaincode)
	if err != nil {
		return err
	}
	//write Key
	binary.BigEndian.PutUint16(tmpBytes, uint16(len(masterpub.Key)))
	_, err = buf.Write(tmpBytes)
	if err != nil {
		return err
	}
	_, err = buf.Write(masterpub.Key)
	if err != nil {
		return err
	}
	//write depth
	binary.BigEndian.PutUint16(tmpBytes, masterpub.Depth)
	_, err = buf.Write(tmpBytes)
	if err != nil {
		return err
	}
	//write to file
	err = ioutil.WriteFile(filename, buf.Bytes(), 0666)
	return err
}

//ReadWalletFromFile read the hdwallet object in the file
func ReadWalletFromFile(filename string) (*HDWallet, error) {
	var err error
	//read from file
	fi, _ := os.Open(filename)
	defer fi.Close()
	masterPubInFile := &HDWallet{}
	r := bufio.NewReader(fi)
	tmpBytes := make([]byte, 2)
	//read vbytes
	r.Read(tmpBytes)
	length := binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.Vbytes = make([]byte, length)
	_, err = r.Read(masterPubInFile.Vbytes)
	if err != nil {
		return nil, err
	}
	//read Fingerprint
	r.Read(tmpBytes)
	length = binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.Fingerprint = make([]byte, length)
	_, err = r.Read(masterPubInFile.Fingerprint)
	if err != nil {
		return nil, err
	}
	//read I
	_, err = r.Read(tmpBytes)
	if err != nil {
		return nil, err
	}
	length = binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.I = make([]byte, length)
	_, err = r.Read(masterPubInFile.I)
	if err != nil {
		return nil, err
	}
	//read Chaincode
	_, err = r.Read(tmpBytes)
	if err != nil {
		return nil, err
	}
	length = binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.Chaincode = make([]byte, length)
	_, err = r.Read(masterPubInFile.Chaincode)
	if err != nil {
		return nil, err
	}
	//read Chaincode
	_, err = r.Read(tmpBytes)
	if err != nil {
		return nil, err
	}
	length = binary.BigEndian.Uint16(tmpBytes)
	masterPubInFile.Key = make([]byte, length)
	_, err = r.Read(masterPubInFile.Key)
	if err != nil {
		return nil, err
	}
	//read depth
	_, err = r.Read(tmpBytes)
	if err != nil {
		return nil, err
	}
	masterPubInFile.Depth = binary.BigEndian.Uint16(tmpBytes)

	return masterPubInFile, nil
}
