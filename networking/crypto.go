package networking

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

// RandomXorKey generate a random number between 0 and the defined xorLimit
func RandomXorKey() uint16 {
	var src cryptoSource
	rnd := rand.New(src)
	return uint16(rnd.Intn(int(xorLimit)))
}

// XorCipher decrypt bytes using captured xorKey and xorTable
func XorCipher(data []byte, xorPos *uint16) {
	for i := range data {
		data[i] ^= xorKey[*xorPos]
		*xorPos++
		if *xorPos >= xorLimit {
			*xorPos = 0
		}
	}
}
