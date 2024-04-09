package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"log"
	"math/big"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w *Wallet) UnmarshalJSON(data []byte) error {
	var raw struct {
		PrivateKey struct {
			PublicKey struct {
				X *big.Int `json:"X"`
				Y *big.Int `json:"Y"`
			} `json:"PublicKey"`
			D *big.Int `json:"D"`
		} `json:"PrivateKey"`
		PublicKey []byte `json:"PublicKey"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	w.PrivateKey.PublicKey.Curve = elliptic.P256()
	w.PrivateKey.PublicKey.X = raw.PrivateKey.PublicKey.X
	w.PrivateKey.PublicKey.Y = raw.PrivateKey.PublicKey.Y
	w.PrivateKey.D = raw.PrivateKey.D
	w.PublicKey = raw.PublicKey

	return nil
}
func (w *Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	fmt.Printf("puc key: %x\n", w.PublicKey)
	fmt.Printf("puc hash: %x\n", pubHash)
	fmt.Printf("address: %x\n", address)

	return address
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}
