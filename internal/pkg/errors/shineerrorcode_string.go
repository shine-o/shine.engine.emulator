// Code generated by "stringer -type=ShineErrorCode"; DO NOT EDIT.

package errors

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PersistenceErrDB-0]
	_ = x[PersistenceItemInvalidAmount-1]
	_ = x[PersistenceItemInvalidShnId-2]
	_ = x[PersistenceItemInvalidCharacterId-3]
	_ = x[PersistenceItemDistinctShnID-4]
	_ = x[PersistenceItemSlotUpdate-5]
	_ = x[PersistenceInventoryFull-6]
	_ = x[PersistenceCharInvalidSlot-7]
	_ = x[PersistenceCharNameTaken-8]
	_ = x[PersistenceCharNoSlot-9]
	_ = x[PersistenceCharInvalidName-10]
	_ = x[PersistenceCharInvalidClassGender-11]
	_ = x[PersistenceCharNotExists-12]
	_ = x[PersistenceUnknownInventory-13]
	_ = x[PersistenceOutOfRangeSlot-14]
	_ = x[ZoneItemEquipFailed-15]
	_ = x[ZoneItemMissingData-16]
	_ = x[ZoneInventorySlotOccupied-17]
	_ = x[ZoneMapNotFound-18]
	_ = x[ZoneUnknownItemClass-19]
	_ = x[ZoneItemSlotChangeNoItem-20]
	_ = x[ZoneItemSlotEquipNoItem-21]
	_ = x[ZoneItemSlotIsBound-22]
	_ = x[ZoneItemSlotInUse-23]
	_ = x[ZoneItemSlotGTS-24]
	_ = x[ZoneItemUnknownInventoryType-25]
	_ = x[ZoneItemDeleteNoItem-26]
	_ = x[ZoneItemNoItemInSlot-27]
	_ = x[ZoneItemSlotIsOccupied-28]
	_ = x[ZoneItemEquipBadType-29]
	_ = x[ZoneItemSlotChangeConstraint-30]
	_ = x[ZoneMissingPlayer-31]
	_ = x[ZoneUnexpectedEvent-32]
	_ = x[ZoneMapCollisionDetected-33]
	_ = x[ZoneUnknownNpcRole-34]
	_ = x[ZoneMissingMapData-35]
	_ = x[ZoneMissingNpcData-36]
	_ = x[UnitTestError-37]
}

const _ShineErrorCode_name = "PersistenceErrDBPersistenceItemInvalidAmountPersistenceItemInvalidShnIdPersistenceItemInvalidCharacterIdPersistenceItemDistinctShnIDPersistenceItemSlotUpdatePersistenceInventoryFullPersistenceCharInvalidSlotPersistenceCharNameTakenPersistenceCharNoSlotPersistenceCharInvalidNamePersistenceCharInvalidClassGenderPersistenceCharNotExistsPersistenceUnknownInventoryPersistenceOutOfRangeSlotZoneItemEquipFailedZoneItemMissingDataZoneInventorySlotOccupiedZoneMapNotFoundZoneUnknownItemClassZoneItemSlotChangeNoItemZoneItemSlotEquipNoItemZoneItemSlotIsBoundZoneItemSlotInUseZoneItemSlotGTSZoneItemUnknownInventoryTypeZoneItemDeleteNoItemZoneItemNoItemInSlotZoneItemSlotIsOccupiedZoneItemEquipBadTypeZoneItemSlotChangeConstraintZoneMissingPlayerZoneUnexpectedEventZoneMapCollisionDetectedZoneUnknownNpcRoleZoneMissingMapDataZoneMissingNpcDataUnitTestError"

var _ShineErrorCode_index = [...]uint16{0, 16, 44, 71, 104, 132, 157, 181, 207, 231, 252, 278, 311, 335, 362, 387, 406, 425, 450, 465, 485, 509, 532, 551, 568, 583, 611, 631, 651, 673, 693, 721, 738, 757, 781, 799, 817, 835, 848}

func (i ShineErrorCode) String() string {
	if i < 0 || i >= ShineErrorCode(len(_ShineErrorCode_index)-1) {
		return "ShineErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ShineErrorCode_name[_ShineErrorCode_index[i]:_ShineErrorCode_index[i+1]]
}
