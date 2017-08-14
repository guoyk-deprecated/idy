package idy

import "time"
import "math/rand"
import "testing"

func TestSliceEdgeSituation(t *testing.T) {
	c := SliceConfig{
		ShardNo:            1,
		ShardCount:         1,
		SliceEffectiveSize: 1,
		SliceSize:          1,
	}
	s := NewSlice(c)
	s.NewSeed()
	s.UpdateElements()
	l := rand.Intn(50) + 10
	for i := 0; i < l; i = i + 1 {
		_, _ = s.NextId()
	}
}

func TestSliceEquality(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	r1 := rand.Intn(10) + 1
	r2 := rand.Intn(2048) + 2048
	c := SliceConfig{
		ShardNo:            r1,
		ShardCount:         rand.Intn(10) + r1,
		SliceEffectiveSize: r2,
		SliceSize:          r2 + rand.Intn(2048),
	}
	s1 := NewSlice(c)
	s1.NewSeed()
	s2 := NewSlice(c)
	s2.Seed = s1.Seed
	s1.UpdateElements()
	s2.UpdateElements()

	for i := 0; i < c.SliceEffectiveSize-1; i = i + 1 {
		v1, _ := s1.NextId()
		v2, _ := s2.NextId()

		if v1 != v2 {
			t.Error("sequence of slice1 does not equals to slice2", "at:", i, v1, "!=", v2, "config:", c)
		}
	}
}
