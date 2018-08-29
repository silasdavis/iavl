package benchmarks

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/tendermint/iavl"
	db "github.com/tendermint/tendermint/libs/db"
)

var dbChoices = []db.DBBackendType{"memdb", "goleveldb"}

func MakeTree(keySize int, dataSize int, numItems int, dataBase db.DB, cacheSize int) {
	rand.Seed(123456789)
	t := iavl.NewMutableTree(dataBase, cacheSize)
	for i := 0; i < numItems; i++ {
		t.Set(randBytes(keySize), randBytes(dataSize))
	}
	t.SaveVersion()
}

// Benchmarks inserting a specific key/value pair into a mutable tree loaded from a db
func benchmarkInsert(b *testing.B, key []byte, value []byte, dataBase db.DB, cacheSize int) {
	tree := iavl.NewMutableTree(dataBase, cacheSize)
	tree.Load()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Set(key, value)
		tree.Rollback()
	}
}

func makeRange(start uint, end uint, step uint) []uint {
	var res []uint
	for i := start; i < end; i += step {
		res = append(res, i)
	}
	return res
}

func makeMultiRangeHelper(ranges [][]uint, cur [][]uint) [][]uint {
	if len(ranges) == 0 {
		fmt.Printf("returning cur: %v\n", cur)
		return cur
	}
	fmt.Printf("making range: %v\n", ranges[0])
	parts := makeRange(ranges[0][0], ranges[0][1], ranges[0][2])
	fmt.Printf("resulting parts: %v\n", parts)
	//nxt := make([][]uint, len(cur)*len(parts))
	var nxt [][]uint
	for i := 0; i < len(parts); i++ {
		for j := 0; j < len(cur); j++ {
			//for _, c := range cur {
			fmt.Printf("nxt before making new nc: %v\n", nxt)
			nc := append(cur[j], parts[i])
			fmt.Printf("made nc: %v\n", nc)
			fmt.Printf("nxt after making new nc, before appending: %v\n", nxt)
			nxt = append(nxt, nc)
			fmt.Printf("new nxt: %v\n", nxt)
			//nxt[(int(i)*len(cur))+j]
		}
	}
	fmt.Printf("recursing with nxt: %v\n", nxt)
	fmt.Printf("and ranges[1:]: %v\n", ranges[1:])
	return makeMultiRangeHelper(ranges[1:], nxt)
}

func MakeMultiRange(ranges [][]uint) [][]uint {
	return makeMultiRangeHelper(ranges, [][]uint{{}})
}

func BenchmarkTryMultiRange(b *testing.B) {
	fmt.Println(MakeMultiRange([][]uint{{5, 16, 5}, {5, 16, 5}, {10, 11, 5}, {10, 11, 5}, {10, 11, 5}, {10, 11, 5}, {0, 2, 1}}))
}

func BenchmarkInsert(b *testing.B) {
	logTreeSizeRange := []uint{5, 16, 5}
	logCacheSizeRange := []uint{5, 16, 5}
	logInsertedKeySize := []uint{10, 11, 5}
	logInsertedDataSize := []uint{10, 11, 5}
	logCurrentKeySize := []uint{10, 11, 5}
	logCurrentDataSize := []uint{10, 11, 5}
	choiceDataBase := []uint{0, 2, 1}
	fmt.Println("tree, cache, insert key, insert data, current key, current data, db")
	params := MakeMultiRange([][]uint{
		logTreeSizeRange, logCacheSizeRange, logInsertedKeySize,
		logInsertedDataSize, logCurrentKeySize, logCurrentDataSize,
		choiceDataBase,
	})
	fmt.Printf("number of datapoints to collect: %d\n", len(params))
	fmt.Println(params)
	for x, p := range params {
		lts := p[0]
		lcs := p[1]
		liks := p[2]
		lids := p[3]
		lcks := p[4]
		lcds := p[5]
		cdb := dbChoices[p[6]]

		dirName := "BenchInsert/test.db"

		dataBase := db.NewDB("test", cdb, dirName)
		MakeTree(1<<lcks, 1<<lcds, 1<<lts, dataBase, 1<<lcs)
		key := randBytes(1 << liks)
		data := randBytes(1 << lids)
		b.Run(fmt.Sprintf("Insert - %d (%d %d %d %d %d %d %s)", x, lts, lcs, liks, lids, lcks, lcds, cdb), func(b *testing.B) {
			benchmarkInsert(b, key, data, dataBase, 1<<lcs)
		})
		dataBase.Close()
		os.RemoveAll(dirName)
	}
}

/*

Actions to benchmark:
	Insert(tree, cache, keyi, datai, keyc, datac)
	Update(tree, cache, keyc, datac)
	Remove(tree, cache, keyc, datac)
	QueryHit(tree, cache, keyc, datac)
	QueryMiss(tree, cache, keyc, datac, keyq)
	Iterate Range(tree, cache, keyc, datac, range)

	Update many and SaveVersion
	DeleteVersion


Numeric Parameters:
	Num items in tree
	Key length of item in question
	Data length of item in question
	Key length distribution in tree
	Data length distribution in tree

Choices
	DB backend



*/
