package lib

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func randomXorKey(xorLimit int) uint16 {
	var src cryptoSource
	rnd := rand.New(src)
	return uint16(rnd.Intn(xorLimit))
}
