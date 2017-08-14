package idy

import "io"
import "encoding/json"

// Database persist struct
type Database struct {
	Version uint   `json:"version"`
	Shard   string `json:"shard"`
	Seed    string `json:"seed"`
	Start   string `json:"start"`
	Index   int    `json:"index"`
}

func (d Database) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(d)
}

func DecodeDatabase(d *Database, r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(d)
}
