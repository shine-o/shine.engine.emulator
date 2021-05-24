package structs

// struct SERVERMENU
type ServerMenu struct {
	Reply   byte
	Content string `struct:"[32]byte"`
}
