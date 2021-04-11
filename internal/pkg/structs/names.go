package structs

type Name3 struct {
	Name string `struct:"[12]byte"`
}

type Name4 struct {
	Name string `struct:"[16]byte"`
}

type Name5 struct {
	Name string `struct:"[20]byte"`
}

type Name8 struct {
	Name string `struct:"[32]byte"`
}

// union Name256Byte
type Name256Byte struct {
	Name string `struct:"[256]byte"`
}
