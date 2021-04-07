package crypto

import (
	crand "crypto/rand"
	"encoding/binary"
	"github.com/google/logger"
	"math/rand"
	"time"
)

type Source struct{}

func (s Source) Seed(seed int64) {}

func (s Source) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s Source) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		logger.Fatal(err)
	}
	return v
}

func RandomUint16Between(min, max uint16) uint16 {
	var src Source
	rand.Seed(time.Now().UnixNano())
	rnd := rand.New(src)
	return uint16(rnd.Intn(int(max-min)) + int(min))
}

func RandomIntBetween(min, max int) int {
	var src Source
	rand.Seed(time.Now().UnixNano())
	rnd := rand.New(src)
	return rnd.Intn(max-min) + min
}

func RandomUint32Between(min, max uint32) uint32 {
	var src Source
	rand.Seed(time.Now().UnixNano())
	rnd := rand.New(src)
	return uint32(rnd.Intn(int(max-min)) + int(min))
}

// RandomXorKey generate a random number between 0 and the defined xorLimit
func RandomXorKey(xorLimit uint16) uint16 {
	return RandomUint16Between(0, xorLimit)
}

// XorCipher decrypt bytes using captured xorKey and xorTable
func XorCipher(data []byte, xorKey []byte, xorPos *uint16, xorLimit uint16) {
	for i := range data {
		data[i] ^= xorKey[*xorPos]
		*xorPos++
		if *xorPos >= xorLimit {
			*xorPos = 0
		}
	}
}
