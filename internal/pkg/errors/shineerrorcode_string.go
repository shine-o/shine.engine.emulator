// Code generated by "stringer -type=ShineErrorCode"; DO NOT EDIT.

package errors

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PersistenceErrDB-0]
	_ = x[PersistenceErrItemInvalidAmount-1]
	_ = x[PersistenceErrItemInvalidShnId-2]
	_ = x[PersistenceErrItemInvalidCharacterId-3]
	_ = x[PersistenceErrItemDistinctShnID-4]
	_ = x[PersistenceErrItemSlotUpdate-5]
	_ = x[PersistenceErrInventoryFull-6]
	_ = x[PersistenceErrCharInvalidSlot-7]
	_ = x[PersistenceErrCharNameTaken-8]
	_ = x[PersistenceErrCharNoSlot-9]
	_ = x[PersistenceErrCharInvalidName-10]
	_ = x[PersistenceErrCharInvalidClassGender-11]
	_ = x[PersistenceErrCharNotExists-12]
	_ = x[ZoneItemEquipFailed-13]
	_ = x[ZoneItemMissingData-14]
	_ = x[ZoneInventorySlotOccupied-15]
}

const _ShineErrorCode_name = "PersistenceErrDBPersistenceErrItemInvalidAmountPersistenceErrItemInvalidShnIdPersistenceErrItemInvalidCharacterIdPersistenceErrItemDistinctShnIDPersistenceErrItemSlotUpdatePersistenceErrInventoryFullPersistenceErrCharInvalidSlotPersistenceErrCharNameTakenPersistenceErrCharNoSlotPersistenceErrCharInvalidNamePersistenceErrCharInvalidClassGenderPersistenceErrCharNotExistsZoneItemEquipFailedZoneItemMissingDataZoneInventorySlotOccupied"

var _ShineErrorCode_index = [...]uint16{0, 16, 47, 77, 113, 144, 172, 199, 228, 255, 279, 308, 344, 371, 390, 409, 434}

func (i ShineErrorCode) String() string {
	if i < 0 || i >= ShineErrorCode(len(_ShineErrorCode_index)-1) {
		return "ShineErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ShineErrorCode_name[_ShineErrorCode_index[i]:_ShineErrorCode_index[i+1]]
}
