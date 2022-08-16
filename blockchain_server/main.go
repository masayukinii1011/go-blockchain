// go run blockchain_server/main.go blockchain_server/blockchain_server.go
package main

import (
	"flag"
)

func main() {
	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Server")
	flag.Parse()
	app := NewBlockchainServer(uint16(*port))
	app.Run()

	/*
		// wallet
		walletMiner := wallet.NewWallet()
		walletA := wallet.NewWallet()
		walletB := wallet.NewWallet()

		// transaction。A から B へ 1.0送る
		t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0)

		// block chain
		blockChain := block.NewBlockchain(walletMiner.BlockchainAddress(),)

		// transaction を block chain に追加
		isAdded := blockChain.AddTransaction(walletA.BlockchainAddress(), walletB.BlockchainAddress(), 1.0, walletA.PublicKey(), t.GenerateSignature())
		fmt.Println(isAdded)

		// マイニングする
		blockChain.Mining()
	*/
}
