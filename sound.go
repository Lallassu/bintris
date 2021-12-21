// This code uses OpenAL Soft via gomobile/exp/audio/al.
// See included license file.
package main

import (
	"fmt"
	"io"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/audio/al"
)

type Sound struct {
	sources   []al.Source
	buffers   []al.Buffer
	sounds    map[string]int
	initiated bool
}

func (s *Sound) Init() {
	if err := al.OpenDevice(); err != nil {
		fmt.Printf("ERROR: failed to open sound device: %v", err)
		return
	}
	s.initiated = true
	s.sounds = make(map[string]int)
}

func (s *Sound) Load(name, file string, format uint32, hz int32) {
	// So we still can play w/o audio
	if !s.initiated {
		return
	}

	f, err := asset.Open(file)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	s.sources = append(s.sources, al.GenSources(1)...)
	s.buffers = append(s.buffers, al.GenBuffers(1)...)
	id := len(s.buffers) - 1
	s.sounds[name] = id

	s.buffers[id].BufferData(format, data, hz)
	s.sources[id].QueueBuffers(s.buffers[id])
}

func (s *Sound) Play(name string) {
	if !s.initiated {
		return
	}

	id, ok := s.sounds[name]
	if !ok {
		fmt.Printf("Sound %v not found\n", name)
		return
	}

	if len(s.sources) < id {
		fmt.Printf("Error: Sound is not loaded? Len: %d Id: %d\n", len(s.sources), id)
		return
	}
	al.PlaySources(s.sources[id])
}

func (s *Sound) Stop(name string) {
	if !s.initiated {
		return
	}

	id, ok := s.sounds[name]
	if !ok {
		fmt.Printf("Sound %v not found\n", name)
		return
	}

	if len(s.sources) < id {
		fmt.Printf("Error: Sound is not loaded? Len: %d Id: %d\n", len(s.sources), id)
		return
	}
	al.StopSources(s.sources[id])
}

func (s *Sound) Close() {
	if !s.initiated {
		return
	}
	for i := range s.sources {
		al.DeleteSources(s.sources[i])
	}

	for i := range s.buffers {
		al.DeleteBuffers(s.buffers[i])
	}
	al.CloseDevice()
}
