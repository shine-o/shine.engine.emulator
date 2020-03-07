package lib

import (
	crand "crypto/rand"
	"encoding/binary"
	"github.com/spf13/viper"
	"math/rand"
)

var xorKey []byte

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

// decrypt encrypted bytes using captured xorKey and xorTable
func xorCipher(eb []byte, xorPos *uint16) {
	xorLimit := uint16(viper.GetInt("crypt.xorLimit"))
	for i, _ := range eb {
		eb[i] ^= xorKey[*xorPos]
		*xorPos++
		if *xorPos >= xorLimit {
			*xorPos = 0
		}
	}
}
