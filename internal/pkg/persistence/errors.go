package persistence

import "fmt"

const (
	ErrDB CustomErrorCode = iota

	ErrItemInvalidAmount
	ErrItemInvalidShnId
	ErrItemInvalidCharacterId
	ErrInventoryFull

	ErrCharInvalidSlot
	ErrCharNameTaken
	ErrCharNoSlot
	ErrCharInvalidName
	ErrCharInvalidClassGender
	ErrCharNotExists
)

//go:generate stringer -type=CustomErrorCode
type CustomErrorCode int

type Err struct {
	Code    CustomErrorCode
	Message string
	Details ErrDetails
}

type ErrDetails map[string]interface{}

func (e Err) Error() string {
	if len(e.Details) > 0 {
		var res = e.Code.String() + " > "
		if e.Message != "" {
			res += e.Message + " > "
		}
		for k, v := range e.Details {
			res += fmt.Sprintf("%v = %v, ", k, v)
		}
		return res
	}

	if e.Message == "" {
		return e.Code.String()
	}

	return e.Code.String() + " > " + e.Message
}

//// ErrInvalidSlot happens if the client tries to bypass client side verification
//var ErrInvalidSlot = &ErrCharacter{
//	Code:    0,
//	Message: "invalid slot",
//}
//
//// ErrNameTaken name is reserved or in use
//var ErrNameTaken = &ErrCharacter{
//	Code:    1,
//	Message: "name taken",
//}
//
//// ErrNoSlot happens if the client tries to bypass client side verification
//var ErrNoSlot = &ErrCharacter{
//	Code:    2,
//	Message: "no slot available",
//}
//
//// ErrInvalidName happens if the client tries to bypass client side verification
//var ErrInvalidName = &ErrCharacter{
//	Code:    3,
//	Message: "invalid name",
//}
//
//// ErrInvalidClassGender happens if the client tries to bypass client side verification
//var ErrInvalidClassGender = &ErrCharacter{
//	Code:    4,
//	Message: "invalid class gender data",
//}
