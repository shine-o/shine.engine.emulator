package structs

//struct SERVERMENU
//{
//  char reply;
//  char string[32];
//};
type ServerMenu struct {
	Reply   byte
	Content string `struct:"[32]byte"`
}
