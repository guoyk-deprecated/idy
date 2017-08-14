package idy

import srand "crypto/rand"
import "strconv"

type Slice struct {
	Config   SliceConfig
	Seed     int64
	Start    uint64
	Index    int
	Elements []uint64
}

func NewSlice(c SliceConfig) *Slice {
	return &Slice{
		Config:   c,
		Elements: make([]uint64, c.SliceSize),
	}
}

func (s *Slice) NextId() (uint64, bool) {
	// if slice exceeded
	if s.Index >= s.Config.SliceEffectiveSize-1 {
		// move to next slice
		s.Start = s.Start + uint64(s.Config.SliceSize*s.Config.ShardCount)
		// reset index to 0
		s.Index = 0
		// create a new Seed
		s.NewSeed()
		// shuffle elements
		s.UpdateElements()
		return s.Elements[s.Index], true
	} else {
		// just increase index
		s.Index = s.Index + 1
		return s.Elements[s.Index], false
	}
}

func (s *Slice) NewSeed() {
	// read 7 Secure PRNG bytes
	seeds := make([]byte, 7)
	srand.Read(seeds)
	// convert to int64
	var val uint64
	for i := 0; i < len(seeds); i = i + 1 {
		val = val + uint64(seeds[i])<<uint64(8*i)
	}
	s.Seed = int64(val)
}

func (s *Slice) UpdateElements() {
	// fill elements
	for i := 0; i < s.Config.SliceSize; i = i + 1 {
		s.Elements[i] = s.Start + uint64((s.Config.ShardNo-1)*s.Config.SliceSize+i)
	}
	// shuffle elements with Seed
	fisherYatesShuffle(s.Elements, s.Seed)
}

func (s *Slice) toDatabase() Database {
	return Database{
		Version: 1,
		Shard:   s.Config.Encode(),
		Seed:    strconv.FormatInt(s.Seed, 10),
		Start:   strconv.FormatUint(s.Start, 10),
		Index:   s.Index,
	}
}
