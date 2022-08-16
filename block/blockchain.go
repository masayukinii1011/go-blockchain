package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go-blockchain/utils"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFICULTY = 3                // 計算に使うハッシュの頭の桁数。難易度の調整に使用。
	MINING_SENDER    = "THE BLOCKCHAIN" // マイニング報酬の送信者
	MINING_REWARD    = 1.0              // マイニング報酬
)

/*
func main() {
	// マイナーのアドレス
	address := "c"

	// ブロックチェーンの作成。初回一つだけトランザクションが空のブロックが追加されている
	blockChain := NewBlockchain(address)

	// プールにトランザクションを追加する
	blockChain.AddTransaction("A", "B", 1.0)

	// マイニングする
	blockChain.Mining()
}
*/
type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string // 報酬の送信先アドレス
	port              uint16
}

// ブロックチェーンを作成する
// 一つだけトランザクションが空のブロックを追加する
func NewBlockchain(blockchainAddress string, port uint16) *Blockchain {
	b := new(Block)
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

// nonce、一つ前のブロックから生成したハッシュ、ブロックチェーンのプールに入っているトランザクション、を使ってブロックを作成する
// 作成したブロックをチェーンに追加する
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	return b
}

// 最後のブロックを返す
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// トランザクションを作成してプールに追加する
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	// 送信者がマイナーであれば(=報酬を与えるとき)
	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	// トランザクションの署名が認証されれば
	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {

		//送信者の残高が足らなければ
		/*
			if bc.CalcurateTotalAmount(sender) < value {
				log.Println("ERROR: Not enough balance in a wallet")
				return false
			}
		*/

		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
		return false
	}
}

// トランザクションの署名の認証
func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)       // トランザクションを JSON に
	h := sha256.Sum256([]byte(m)) // トランザクションのハッシュを生成
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

// プールのトランザクションをコピーする
func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(
				t.senderBlockchainAddress,
				t.recipientBlockchainAddress,
				t.value,
			),
		)
	}
	return transactions
}

/*
  nonce + previousHash + transaction = 000........
	となるような nonce が見つかると、コンセンサスが取れて、ブロックが追加できる
	「000」の桁数は difficulty で調整。
*/

// ValidProof が true になるまで nonce を変更して試す
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFICULTY) {
		nonce += 1
	}
	return nonce
}

// 入力した nonce が正しいか判定する
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)                  // difficulty 桁分の "0" の文字列
	guessBlock := Block{0, nonce, previousHash, transactions} // nonce からブロックを作成
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())      // 作成したブロックからハッシュを生成
	return guessHashStr[:difficulty] == zeros                 // ハッシュの頭が "000" であれば、入力した nonce は正しい
}

// マイニング
func (bc *Blockchain) Mining() bool {
	// MINING_SENDER から、登録されているアドレスに、報酬を送る、トランザクションを追加
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)

	// 計算して nonce を求める
	nonce := bc.ProofOfWork()

	// 最後のブロックからハッシュを生成する
	previousHash := bc.LastBlock().Hash()

	// 最後のブロックから生成したハッシュと nonce を使ってブロックを作成
	bc.CreateBlock(nonce, previousHash)

	return true
}

// アドレスのトランザクションの合計を計算する
func (bc *Blockchain) CalcurateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

type Block struct {
	timestamp    int64          // ブロックの作成時の UnixTime
	nonce        int            // (number used once)。使い捨ての32bit値。マイニングで使用。
	previousHash [32]byte       // 一つ前のブロックから生成したハッシュ
	transactions []*Transaction // トランザクションのまとまり
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}

// ブロックの Json からハッシュを生成する
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

// ブロックの private フィールドを public に変換して Json 化する
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

type Transaction struct {
	senderBlockchainAddress    string  // 送信者アドレス
	recipientBlockchainAddress string  //受信者アドレス
	value                      float32 // 値
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string
		Recipient string
		Value     float32
	}{
		Sender:    t.senderBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value:     t.value,
	})
}
