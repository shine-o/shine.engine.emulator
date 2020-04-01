package structs

// /* 2839 */
// union Name3
// {
//  char n3_name[12];
//  unsigned int n3_code[3];
// };
type Name3 struct {
	Name [12]byte `struct:"[12]byte"`
}

/* 3801 */
// union Name4
// {
//	char n4_name[16];
//	unsigned int n4_code[4];
// };
type Name4 struct {
	Name [16]byte `struct:"[16]byte"`
}

// /* 3256 */
// union Name5
// {
//  char n5_name[20];
//  unsigned int n5_code[5];
// };
type Name5 struct {
	Name [20]byte `struct:"[20]byte"`
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

// union Name256Byte
// {
//	char n256_name[256];
//	unsigned __int64 n256_code[32];
// };
type Name256Byte struct {
	Name [256]byte `struct:"[256]byte"`
}
