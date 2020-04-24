package structs

import "strings"

// /* 2839 */
// union Name3
// {
//  char n3_name[12];
//  unsigned int n3_code[3];
// };
type Name3 struct {
	Name [12]byte `struct:"[12]byte"`
}

func (n *Name3) String() string {
	b := n.Name[:]
	return strings.TrimRight(string(b), "\x00")
}

func NewName3(n string) Name3 {
	n3 := Name3{}
	copy(n3.Name[:], n)
	return n3
}

// union Name4
// {
//	char n4_name[16];
//	unsigned int n4_code[4];
// };
type Name4 struct {
	Name [16]byte `struct:"[16]byte"`
}

func (n *Name4) String() string {
	b := n.Name[:]
	return strings.TrimRight(string(b), "\x00")
}

func NewName4(n string) Name4 {
	n4 := Name4{}
	copy(n4.Name[:], n)
	return n4
}

// /* 3256 */
// union Name5
// {
//  char n5_name[20];
//  unsigned int n5_code[5];
// };
type Name5 struct {
	//Name [20]byte `struct:"[20]byte"`
	Name [20]byte `struct:"[20]byte"`
}

func (n *Name5) String() string {
	b := n.Name[:]
	return strings.TrimRight(string(b), "\x00")
}

func NewName5(n string) Name5 {
	n5 := Name5{}
	copy(n5.Name[:], n)
	return n5
}

///* 2787 */
// union Name8
// {
//	char n8_name[32];
//	unsigned int n8_code[8];
// };
type Name8 struct {
	Name [32]byte `struct:"[32]byte"`
}

func (n *Name8) String() string {
	b := n.Name[:]
	return strings.TrimRight(string(b), "\x00")
}

func NewName8(n string) Name8 {
	n8 := Name8{}
	copy(n8.Name[:], n)
	return n8
}

// union Name256Byte
// {
//	char n256_name[256];
//	unsigned __int64 n256_code[32];
// };
type Name256Byte struct {
	Name [256]byte `struct:"[256]byte"`
}

func (n *Name256Byte) String() string {
	b := n.Name[:]
	return strings.TrimRight(string(b), "\x00")
}

func NewName256Byte(n string) Name256Byte {
	n256 := Name256Byte{}
	copy(n256.Name[:], n)
	return n256
}
