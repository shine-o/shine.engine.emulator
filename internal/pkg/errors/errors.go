package errors

import "fmt"

const (
	PersistenceErrDB ShineErrorCode = iota

	PersistenceErrItemInvalidAmount
	PersistenceErrItemInvalidShnId
	PersistenceErrItemInvalidCharacterId
	PersistenceErrItemDistinctShnID
	PersistenceErrItemSlotUpdate

	PersistenceErrInventoryFull

	PersistenceErrCharInvalidSlot
	PersistenceErrCharNameTaken
	PersistenceErrCharNoSlot
	PersistenceErrCharInvalidName
	PersistenceErrCharInvalidClassGender
	PersistenceErrCharNotExists

	ZoneItemEquipFailed
	ZoneItemMissingData
)

//go:generate stringer -type=ShineErrorCode
type ShineErrorCode int

type Err struct {
	Code    ShineErrorCode
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
