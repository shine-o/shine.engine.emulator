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
	_ = x[LoginMismatchedEventType-15]
	_ = x[WorldMismatchedEventType-16]
	_ = x[WorldBadSessionType-17]
	_ = x[WorldPacketUnpackFailed-18]
	_ = x[ZoneItemEquipFailed-19]
	_ = x[ZoneItemMissingData-20]
	_ = x[ZoneInventorySlotOccupied-21]
	_ = x[ZoneMapNotFound-22]
	_ = x[ZoneUnknownItemClass-23]
	_ = x[ZoneItemSlotChangeNoItem-24]
	_ = x[ZoneItemSlotEquipNoItem-25]
	_ = x[ZoneItemSlotIsBound-26]
	_ = x[ZoneItemSlotInUse-27]
	_ = x[ZoneItemSlotGTS-28]
	_ = x[ZoneItemUnknownInventoryType-29]
	_ = x[ZoneItemDeleteNoItem-30]
	_ = x[ZoneItemNoItemInSlot-31]
	_ = x[ZoneItemSlotIsOccupied-32]
	_ = x[ZoneItemEquipBadType-33]
	_ = x[ZoneItemSlotChangeConstraint-34]
	_ = x[ZoneMissingPlayer-35]
	_ = x[ZoneUnexpectedEventType-36]
	_ = x[ZoneMapCollisionDetected-37]
	_ = x[ZoneUnknownNpcRole-38]
	_ = x[ZoneMissingMapData-39]
	_ = x[ZoneMissingNpcData-40]
	_ = x[ZoneHandlerMaxReached-41]
	_ = x[ZoneHandlerCapacityWarning-42]
	_ = x[ZoneHandlerMaxAttemptsReached-43]
	_ = x[ZoneBadEntityType-44]
	_ = x[ZoneNilPlayerFields-45]
	_ = x[ZonePlayerSelectedUnknownEntity-46]
	_ = x[ZoneNoSessionAvailable-47]
	_ = x[ZoneEntityNotFound-48]
	_ = x[UnitTestError-49]
	_ = x[PacketSnifferNotEnoughData-50]
}

const _ShineErrorCode_name = "PersistenceErrDBPersistenceItemInvalidAmountPersistenceItemInvalidShnIdPersistenceItemInvalidCharacterIdPersistenceItemDistinctShnIDPersistenceItemSlotUpdatePersistenceInventoryFullPersistenceCharInvalidSlotPersistenceCharNameTakenPersistenceCharNoSlotPersistenceCharInvalidNamePersistenceCharInvalidClassGenderPersistenceCharNotExistsPersistenceUnknownInventoryPersistenceOutOfRangeSlotLoginMismatchedEventTypeWorldMismatchedEventTypeWorldBadSessionTypeWorldPacketUnpackFailedZoneItemEquipFailedZoneItemMissingDataZoneInventorySlotOccupiedZoneMapNotFoundZoneUnknownItemClassZoneItemSlotChangeNoItemZoneItemSlotEquipNoItemZoneItemSlotIsBoundZoneItemSlotInUseZoneItemSlotGTSZoneItemUnknownInventoryTypeZoneItemDeleteNoItemZoneItemNoItemInSlotZoneItemSlotIsOccupiedZoneItemEquipBadTypeZoneItemSlotChangeConstraintZoneMissingPlayerZoneUnexpectedEventTypeZoneMapCollisionDetectedZoneUnknownNpcRoleZoneMissingMapDataZoneMissingNpcDataZoneHandlerMaxReachedZoneHandlerCapacityWarningZoneHandlerMaxAttemptsReachedZoneBadEntityTypeZoneNilPlayerFieldsZonePlayerSelectedUnknownEntityZoneNoSessionAvailableZoneEntityNotFoundUnitTestErrorPacketSnifferNotEnoughData"

var _ShineErrorCode_index = [...]uint16{0, 16, 44, 71, 104, 132, 157, 181, 207, 231, 252, 278, 311, 335, 362, 387, 411, 435, 454, 477, 496, 515, 540, 555, 575, 599, 622, 641, 658, 673, 701, 721, 741, 763, 783, 811, 828, 851, 875, 893, 911, 929, 950, 976, 1005, 1022, 1041, 1072, 1094, 1112, 1125, 1151}

func (i ShineErrorCode) String() string {
	if i < 0 || i >= ShineErrorCode(len(_ShineErrorCode_index)-1) {
		return "ShineErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ShineErrorCode_name[_ShineErrorCode_index[i]:_ShineErrorCode_index[i+1]]
}
