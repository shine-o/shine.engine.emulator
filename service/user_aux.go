// auxiliar structs and functions for Network Commands
package service

//import networking "github.com/shine-o/shine.engine.networking"

/* 3256 */
// union Name5
// {
//	char n5_name[20];
//	unsigned int n5_code[5];
// };
type ComplexName1 struct {
	Name     [20]byte
	NameCode [5]uint32
}

/* 3801 */
// seems like a utility map struct for names, maybe related with NIF files
// union Name4
// {
//	char n4_name[16];
//	unsigned int n4_code[4];
// };
type Name4 struct {
	Name     [16]byte
	NameCode [4]uint16
}
