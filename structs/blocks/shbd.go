package blocks

type SHBD struct {
	X uint32
	Y uint32
	Data []byte `struct-size:"X * Y"`
}

