package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

type BlockChain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

// marshal： 把数据编码成json字符串
func (b Block) calculateHash() string {
	// data是解码json之后得到的结果 _表示err
	data, _ := json.Marshal(b.data)
	// string(data)把json转换为字符串
	// /*/todo strconv.Itoa这个package不知道做啥的 之后查一下
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	// 计算blockhash 所以blockhash是一堆blockdata生成的hash
	blockHash := sha256.Sum256([]byte(blockData))
	// %x是golang的格式占位符
	// %x  十六进制，小写字母，每字节两个字符		Printf("%x", "golang")
	// %X      十六进制，大写字母，每字节两个字符      Printf("%X", "golang")
	// sprintf这个用法挺神奇的
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}
