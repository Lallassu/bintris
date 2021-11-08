// This code uses OpenAL Soft via gomobile/exp/audio/al.
// See included license file.
package main

import (
	"fmt"

	"golang.org/x/mobile/exp/audio/al"
)

type Sound struct {
}

func (s *Sound) Init() {
	if err := al.OpenDevice(); err != nil {
		fmt.Printf("ERROR: failed to open sound device\n")
	}
}
