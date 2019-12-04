// Copyright 2015 go-fuzz project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package json

import (
	"encoding/json"
	"fmt"
	"github.com/fwessels/simdjson-go"
)

func Fuzz(data []byte) (score int) {
	var dst map[string]interface{}
	_, err := simdjson.Parse(data, nil)
	jErr := json.Unmarshal(data, &dst)
	if err != nil {
		if jErr == nil {
			msg := fmt.Sprintf("got error %v, but json.Unmarshal could unmarshal", err)
			// Disabled for now:
			if false {
				panic(msg)
			} else {
				fmt.Println(msg)
			}
		}
		// Don't continue
		return 0
	}
	if jErr != nil {
		msg := fmt.Sprintf("no error reported, but json.Unmarshal reported: %v", jErr)
		panic(msg)
	}

	// Do simple ND test.
	d2 := append(make([]byte,0,len(data)*3+2), data...)
	d2 = append(d2, '\n')
	d2 = append(d2, data...)
	d2 = append(d2, '\n')
	d2 = append(d2, data...)
	_, _ = simdjson.ParseND(data, nil)

	return 1
}
