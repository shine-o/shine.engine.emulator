package login

import (
	"time"
)

// User model for schema: accounts
type User struct {
	tableName struct{} `pg:"accounts.users"`
	ID        uint64
	UserName  string
	Password  string
	DeletedAt time.Time `pg:"soft_delete"`
}

//
//func md5Hash(text string) string {
//	hasher := md5.New()
//	hasher.Write([]byte(text))
//	return hex.EncodeToString(hasher.Sum(nil))
//}