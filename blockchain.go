package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// block中的data试一个map（键是string 值是interface（可以是任意的东西）==》是用来存储交易的信息（from to amount
// pow也就是工作量证明 你跑了多少遍挖矿函数 pow就是几
type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

// marshal： 把数据编码成json字符串
// 传入区块 返回blockhash的函数 调用方法：b.calculateHash()
func (b Block) calculateHash() string {
	// data是解码json之后得到的结果 _表示err
	data, _ := json.Marshal(b.data)
	// string(data)把json转换为字符串
	// iota函数：
	// strconv.Itoa函数的参数int，它可以将数字转换成对应的string的数字
	// 所以就是把pow这个int转换为string然后和前面的字符串加起来得到blockdata
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	// 计算blockhash 所以blockhash是一堆blockdata生成的hash
	// 解释一下hash的特点：
	// 每个字符串都对应着唯一的hash
	// 不同长度的字符串 输出都是的hash都是相同长度的
	// 可以由x字符串很容易来得到hesh(x) 但是从hash(x)反推x几乎不可能
	blockHash := sha256.Sum256([]byte(blockData))
	// %x是golang的格式占位符
	// %x  十六进制，小写字母，每字节两个字符		Printf("%x", "golang")
	// %X      十六进制，大写字母，每字节两个字符      Printf("%X", "golang")
	// sprintf这个用法挺神奇的
	return fmt.Sprintf("%x", blockHash)
}

// 挖一个区块的意思是你要生成一个满足区块难度要求的blockhash
// 比如区块难度是3 区块hash就要以000开头
//

// hasprefix源码：
// // HasPrefix tests whether the string s begins with prefix.
// func HasPrefix(s, prefix string) bool {
// 	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
// }

// hasprefix函数检查blockhash的开头是不是满足difficulty的要求
// 如果不满足就一直试hash函数 直到满足
// 这里的pow就是矿工要解出来的随机数（应为这个随机数设置为简单的递增函数 所以随机数也代表你调用挖矿函数的次数==》工作量证明
// 传入区块和难度值 一直试一直试 给区块的block填入符合区块难度要求的hash和pow值（挖矿函数
// 调用方法：b.mine
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

// 创建创世区块 返回的blockchain是前面已经定义的（包含genesis block， chain数组，difficulty
// 这个函数是返回了一个以当前时间作为时间戳 区块数组中只有genesis block的blockchain 区块难度是手动设置的
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

// 增加新区块
// 传入整个区块链 from to amount
// blockdata是新区块里需要写入的交易数据
func (b *Blockchain) addBlock(from, to string, amount float64) {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	// lastblock是拿到了传入的整条区块链中的最后一个block
	lastBlock := b.chain[len(b.chain)-1]
	// 构建新区块
	// 传入填好交易信息的data 构建一个block
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	// 挖矿给新区块填入hash
	newBlock.mine(b.difficulty)
	// 把这个新的区块加到整个区块链里
	b.chain = append(b.chain, newBlock)
}

// 然后上面这些代码就让我们可以不断地挖矿、写入交易、增加区块形成区块链了
// 但是还有一个非常重要的事情： 检验交易的合法性

// 传入整个区块链 这个函数会检验交易是否合法 来返回一个布尔值（所以传统的区块链在检验交易合法性的时候是整条整条检验的==》也就是要同步全节点==》效率好低==》sharding celestia iota solona这些项目都在试图用不同方法解决这个bug
func (b Blockchain) isValid() bool {

	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		// 拿区块里的数据（包含矿工解出来的pow随机数重新算一遍 看这个hash和block里存储的hash是否相等
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := CreateBlockchain(2)

	// record transactions on the blockchain for Alice, Bob, and John
	// 一个区块一个交易
	blockchain.addBlock("Alice", "Bob", 5)
	blockchain.addBlock("John", "Bob", 2)

	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.isValid())
}
