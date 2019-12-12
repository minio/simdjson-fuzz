// Copyright 2015 go-fuzz project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/fwessels/simdjson-go"
)

// FuzzCorrect will check for correctness and compare output to stdlib.
func FuzzCorrect(data []byte) (score int) {
	const (
		// fail if simdjson doesn't report error, but json.Unmarshal does
		failOnMissingError = false
		// Run input through json.Unmarshal/json.Marshal first
		filterRaw = true
	)

	var want map[string]interface{}
	var wantA []interface{}
	if !utf8.Valid(data) {
		return 0
	}
	if filterRaw {
		var tmp interface{}
		err := json.Unmarshal(data, &tmp)
		if err != nil {
			return 0
		}
		data, err = json.Marshal(tmp)
		if err != nil {
			panic(err)
		}
	}
	pj, err := simdjson.Parse(data, nil)
	jErr := json.Unmarshal(data, &want)
	if err != nil {
		if jErr == nil {
			b, _ := json.Marshal(want)
			msg := fmt.Sprintf("got error %v, but json.Unmarshal could unmarshal to %#v js: %s", err, want, string(b))
			// Disabled for now:
			if true {
				panic(msg)
			} else {
				fmt.Println(msg)
			}
		}
		// Don't continue
		return 0
	}
	if jErr != nil {
		want = nil
		if strings.Contains(jErr.Error(), "cannot unmarshal array into") {
			jErr2 := json.Unmarshal(data, &wantA)
			if jErr2 != nil {
				msg := fmt.Sprintf("no error reported, but json.Unmarshal (Array) reported: %v", jErr2)
				if failOnMissingError {
					panic(msg)
				} else {
					fmt.Println(msg)
				}
			}
		} else {
			msg := fmt.Sprintf("no error reported, but json.Unmarshal reported: %v", jErr)
			if failOnMissingError {
				panic(msg)
			} else {
				fmt.Println(msg)
			}
			// TODO: We should investigate these later.
			return 1
		}
	}
	// Check if we can convert back
	var got map[string]interface{}
	var gotA []interface{}

	i := pj.Iter()
	if i.PeekNextTag() == simdjson.TagEnd {
		if len(want)+len(wantA) > 0 {
			msg := fmt.Sprintf("stdlib returned data %#v, but nothing from simdjson (tap:%d, str:%d, err:%v)", want, len(pj.Tape), len(pj.Strings), err)
			panic(msg)
		}
		return 0
	}

	data, err = i.MarshalJSON()
	if err != nil {
		switch {
		// This is ok.
		case strings.Contains(err.Error(), "INF or NaN number found"):
		default:
			panic(err)
		}
	}
	if want != nil {
		// We should be able to unmarshal into msi
		err := json.Unmarshal(data, &got)
		if err != nil {
			msg := fmt.Sprintf("unable to marshal output: %v", jErr)
			panic(msg)
		}
	}
	if wantA != nil {
		// We should be able to unmarshal into msi
		err := json.Unmarshal(data, &gotA)
		if err != nil {
			msg := fmt.Sprintf("unable to marshal array output: %v", jErr)
			panic(msg)
		}
	}
	if !reflect.DeepEqual(want, got) {
		msg := fmt.Sprintf("Marshal data mismatch:\nstdlib: %#v\nsimdjson:%#v", want, got)
		panic(msg)
	}
	if !reflect.DeepEqual(wantA, gotA) {
		msg := fmt.Sprintf("Marshal array mismatch:\nstdlib: %#v\nsimdjson:%#v", wantA, gotA)
		panic(msg)
	}

	return 1
}
