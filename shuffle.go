package idy

import "github.com/yanke-guo/idy/rand"

func fisherYatesShuffle(slice []uint64, seed int64) {
	rnd := rand.New(rand.NewSource(seed))
	for i := len(slice) - 1; i > 0; i = i - 1 {
		j := rnd.Intn(i)
		k := slice[i]
		slice[i] = slice[j]
		slice[j] = k
	}
}
