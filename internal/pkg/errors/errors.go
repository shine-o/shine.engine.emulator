package errors

import "fmt"

const (
	PersistenceErrDB ShineErrorCode = iota

	PersistenceItemInvalidAmount
	PersistenceItemInvalidShnId
	PersistenceItemInvalidCharacterId
	PersistenceItemDistinctShnID
	PersistenceItemSlotUpdate

	PersistenceInventoryFull

	PersistenceCharInvalidSlot
	PersistenceCharNameTaken
	PersistenceCharNoSlot
	PersistenceCharInvalidName
	PersistenceCharInvalidClassGender
	PersistenceCharNotExists
	PersistenceUnknownInventory
	PersistenceOutOfRangeSlot

	ZoneItemEquipFailed
	ZoneItemMissingData
	ZoneInventorySlotOccupied
	ZoneMapNotFound
	ZoneUnknownItemClass
	ZoneItemSlotChangeNoItem
	ZoneItemSlotEquipNoItem
	ZoneItemSlotIsBound
	ZoneItemSlotInUse
	ZoneItemSlotGTS //Guild Tournament Storage
	ZoneItemUnknownInventoryType
	ZoneItemDeleteNoItem
	ZoneItemNoItemInSlot
	ZoneItemSlotIsOccupied
	ZoneItemEquipBadType
	ZoneItemSlotChangeConstraint
	ZoneMissingPlayer
	ZoneUnexpectedEvent
	ZoneMapCollisionDetected
	ZoneUnknownNpcRole
	ZoneMissingMapData
	ZoneMissingNpcData
	UnitTestError
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
