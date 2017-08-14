package idy

import "errors"
import "strconv"
import "strings"
import "fmt"

type SliceConfig struct {
	ShardNo            int
	ShardCount         int
	SliceEffectiveSize int
	SliceSize          int
}

func DecodeSliceConfig(s string, c *SliceConfig) error {
	ss := strings.Split(s, ":")
	if len(ss) != 4 {
		return errors.New("bad slice config, need 4 components")
	}

	c.ShardNo, _ = strconv.Atoi(ss[0])
	c.ShardCount, _ = strconv.Atoi(ss[1])
	c.SliceEffectiveSize, _ = strconv.Atoi(ss[2])
	c.SliceSize, _ = strconv.Atoi(ss[3])

	return c.Validate()
}

func (s SliceConfig) Validate() error {
	if s.ShardNo < 1 {
		return errors.New("bad shard number (first value in 'shard' option)")
	}
	if s.ShardCount < s.ShardNo {
		return errors.New("bad shard count (second value in 'shard' option)")
	}
	if s.SliceEffectiveSize < 1 {
		return errors.New("bad effective size (third value in 'shard' option)")
	}
	if s.SliceSize < s.SliceEffectiveSize {
		return errors.New("bad slice size (fourth value in 'shard' option)")
	}
	return nil
}

func (s SliceConfig) Encode() string {
	return fmt.Sprintf("%d:%d:%d:%d", s.ShardNo, s.ShardCount, s.SliceEffectiveSize, s.SliceSize)
}
