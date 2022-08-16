package main

import (
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"
)

//var cache = make(map[string]*block.Blockchain)
const tempDir = "wallet_server/templates"

type WalletServer struct {
	port    uint16
	gateway string // Gateway アドレス
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		t.Execute(w, "")
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}

/*
func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	// キャッシュに有ればキャッシュのブロックチェーンを使う
	bc, ok := cache["blockchain"]

	// 無ければ新規作成してキャッシュに入れる
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minersWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
		log.Printf("private_key %v", minersWallet.PrivateKeyStr())
		log.Printf("public_key %v", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address %v", minersWallet.BlockchainAddress())
	}

	return bc
}
*/
