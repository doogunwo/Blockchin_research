package pow

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewPow (b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b,target}
	return pow
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, num)

	return buff.Bytes()
}

func (pow *ProofOfWork) prepareData (nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.prevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
			data := pow.prepareData(nonce)
			hash = sha256.Sum256(data)
			fmt.Printf("\r%x", hash)

			hashInt.SetBytes(hash[:])
			if hashInt.Cmp(pow.target) == -1 {
					break
			} else {
					nonce++
			}
	}
	fmt.Print("mining complete")

	return nonce, hash[:]
}


func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid = hashInt.Cmp(pow.target) == -1
	return isValid
}