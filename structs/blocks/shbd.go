package blocks

type SHBD struct {
	X    int    `struct:"int32"`
	Y    int    `struct:"int32"`
	Data []byte `struct-size:"X * Y"`
}
