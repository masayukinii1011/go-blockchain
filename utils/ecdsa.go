package utils

import "math/big"

// https://zoom-blc.com/what-is-ecdsa
type Signature struct {
	R *big.Int // 公開鍵で言うところのX座標のようなもの
	S *big.Int // R,ハッシュ,秘密鍵を組み合わせた特定の計算結果
}
