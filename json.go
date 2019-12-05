// Copyright 2015 go-fuzz project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package json

import (
	"encoding/json"
	"fmt"
	"github.com/fwessels/simdjson-go"
	"strings"
)

func Fuzz(data []byte) (score int) {
	var dst map[string]interface{}
	var dstA []interface{}
	pj, err := simdjson.Parse(data, nil)
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
		if strings.Contains(jErr.Error(), "cannot unmarshal array into") {
			jErr2 := json.Unmarshal(data, &dstA)
			if jErr2 != nil {
				msg := fmt.Sprintf("no error reported, but json.Unmarshal (Array) reported: %v", jErr2)
				if false {
					panic(msg)
				} else {
					fmt.Println(msg)
				}
			}
		} else {
			msg := fmt.Sprintf("no error reported, but json.Unmarshal reported: %v", jErr)
			if false {
				panic(msg)
			} else {
				fmt.Println(msg)
			}
		}
	}
	// Check if we can convert back
	i := pj.Iter()
	if i.PeekNextTag() != simdjson.TagEnd {
		_, err = i.MarshalJSON()
		if err != nil {
			panic(err)
		}
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
