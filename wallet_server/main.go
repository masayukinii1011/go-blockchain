// go run wallet_server/main.go wallet_server/wallet_server.go
package main

import (
	"flag"
)

func main() {
	port := flag.Uint("port", 8080, "TCP Port Number for Wallet Server")
	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Blockchain Gateway")
	flag.Parse()

	app := NewWalletServer(uint16(*port), *gateway)
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
