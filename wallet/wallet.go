package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go-blockchain/utils"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {
	w := new(Wallet)
	// publicKey と privateKey の生成
	// https://pkg.go.dev/crypto/ecdsa#pkg-overview Example
	// P256 の楕円曲線と乱数でキーを生成。
	// 楕円曲線DSA ( Elliptic Curve Digital Signature Algorithm )
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey

	// Public Key を Address に変換する。
	// 観点。アドレスは短く。セキュアに。
	// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
	// How to create Bitcoin Address

	// 2 - Perform SHA-256 hashing on the public key
	// public keyから、SHA-256 でハッシュを作る
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// 3 - Perform RIPEMD-160 hashing on the result of SHA-256
	// SHA-256 で作ったハッシュから、RIPEMD-160 でハッシュを作る
	// SHA-256 (32 bytes) より RIPEMD-160 (20 bytes) の方が短いハッシュが作れる
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	// 4 - Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
	// メインネットの場合は 0x00 をアタマにつける
	vd4 := make([]byte, 21)   // 空の 21 byte を生成
	vd4[0] = 0x00             // 0x00 を1番目に
	copy(vd4[1:], digest3[:]) // 2 ~ 21番目に RIPEMD-160 (20 bytes) のハッシュを入れる

	// 5 - Perform SHA-256 hash on the extended RIPEMD-160 result
	// 4 の結果から、SHA-256 でハッシュを作る
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// 6 - Perform SHA-256 hash on the result of the previous SHA-256 hash
	// 5 の結果から、SHA-256 でハッシュを作る
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// 7 - Take the first 4 bytes of the second SHA-256 hash. This is the address checksum
	// 6 の結果の 1 ~ 4 番目はチェックサム(誤り検出符号)
	chsum := digest6[:4]

	// 8 - Add the 4 checksum bytes from stage 7 at the end of extended RIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.
	// 4 の結果のアタマにチェックサムを付ける
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])

	// 9 - Convert the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format
	// 8の結果を base58 で暗号化したものがビットコインのブロックチェーンアドレス
	address := base58.Encode(dc8)
	w.blockchainAddress = address

	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

type Transaction struct {
	senderPrivateKey           *ecdsa.PrivateKey
	senderPublicKey            *ecdsa.PublicKey
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, value float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}

// ECDSA署名の生成
func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)       // トランザクションを JSON に
	h := sha256.Sum256([]byte(m)) // トランザクションのハッシュを生成
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}
