// Copyright 2015 go-fuzz project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package json

import (
	"bytes"
	"fmt"
	"github.com/minio/simdjson-go"
)

// FuzzCorrect will check for correctness and compare output to stdlib.
func FuzzSerialize(data []byte) (score int) {
	// Create a tape from the input and ensure that the output of JSON matches.
	pj, err := simdjson.Parse(data, nil)
	if err != nil {
		pj, err = simdjson.ParseND(data, pj)
		if err != nil {
			// Don't continue
			return 0
		}
	}
	i := pj.Iter()
	want, err := i.MarshalJSON()
	if err != nil {
		panic(err)
	}
	// Check if we can convert back
	s := simdjson.NewSerializer()
	got := make([]byte,0, len(want))
	var dst []byte
	var target *simdjson.ParsedJson
	for _, comp := range []simdjson.CompressMode{simdjson.CompressNone, simdjson.CompressFast, simdjson.CompressDefault, simdjson.CompressBest}{
		level := fmt.Sprintf("level-%d:", comp)
		s.CompressMode(comp)
		dst := s.Serialize(dst[:0], *pj)
		target, err = s.Deserialize(dst, target)
		if err != nil {
			panic(level+err.Error())
		}
		i := target.Iter()
		got, err = i.MarshalJSONBuffer(got[:0])
		if err != nil {
			panic(level+err.Error())
		}
		if !bytes.Equal(want, got) {
			err := fmt.Sprintf("%s JSON mismatch:\nwant: %s\ngot :%s", level, string(want), string(got))
			err += fmt.Sprintf("\ntap0:%x", pj.Tape)
			err += fmt.Sprintf("\ntap1:%x", target.Tape)
			panic(err)
		}
	}
	return 1
}
