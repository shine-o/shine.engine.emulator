package structs

// struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_CMD
type NcBriefInfoAbstateChangeCmd struct {
	Handle uint16
	Info   AbstateInformation
}

// struct PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD
type NcBriefInfoDeleteCmd struct {
	Handle uint16
}

// struct PROTO_NC_BRIEFINFO_DROPEDITEM_CMD
type NcBriefInfoDroppedItemCmd struct {
	Handle        uint16
	ItemID        uint16
	Location      ShineXYType
	DropMobHandle uint16
	Attr          NcBriefInfoDroppedItemCmdAttr
}

// struct PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD
type NcBriefInfoChangeDecorateCmd struct {
	Handle uint16
	Item   uint16
	Slot   byte
}

// struct PROTO_NC_BRIEFINFO_MOB_CMD
type NcBriefInfoMobCmd struct {
	MobNum byte
	Mobs   []NcBriefInfoRegenMobCmd `struct:"sizefrom=MobNum"`
}

// struct PROTO_NC_BRIEFINFO_UNEQUIP_CMD
type NcBriefInfoUnEquipCmd struct {
	Handle uint16
	Slot   byte
}

// struct PROTO_NC_BRIEFINFO_REGENMOB_CMD
type NcBriefInfoRegenMobCmd struct {
	Handle uint16
	Mode   byte
	MobID  uint16
	Coord  ShineCoordType
	// 0,1 FlagData size depends on this flag
	// if 1, FlagData is 12 bytes
	FlagState      byte
	FlagData       BriefInfoRegenMobCmdFlag
	Animation      [32]byte
	AnimationLevel byte
	KQTeamType     byte
	RegenAni       byte
}

// struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD
type NcBriefInfoAbstateChangeListCmd struct {
	Handle uint16
	Count  byte
	List   []AbstateInformation `struct:"sizefrom=Count"`
}

// struct PROTO_NC_BRIEFINFO_CHARACTER_CMD
type NcBriefInfoCharacterCmd struct {
	Number     byte
	Characters []NcBriefInfoLoginCharacterCmd `struct:"sizefrom=Number"`
}

// struct PROTO_NC_BRIEFINFO_REGENMOVER_CMD
type NcBriefInfoRegenMoverCmd struct {
	Handle      uint16
	ID          uint32
	HP          uint32
	Coordinates ShineCoordType
	AbstateBit  AbstateBit
	Grade       byte
	SlotHandle  [10]uint16
}

// struct PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD
type NcBriefInfoChangeUpgradeCmd struct {
	Handle  uint16
	Item    uint16
	Upgrade byte
	Slot    byte
}

// struct PROTO_NC_BRIEFINFO_MOVER_CMD
type NcBriefInfoMoverCmd struct {
	Count  byte
	Movers []NcBriefInfoRegenMoverCmd `struct:"sizefrom=Count"`
}

// struct PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD
type NcBriefInfoLoginCharacterCmd struct {
	Handle          uint16
	CharID          Name5
	Coordinates     ShineCoordType
	Mode            byte
	Class           byte
	Shape           ProtoAvatarShapeInfo
	ShapeData       NcBriefInfoLoginCharacterCmdShapeData
	Polymorph       uint16
	Emoticon        StopEmoticonDescript
	CharTitle       CharTitleBriefInfo
	AbstateBit      AbstateBit
	MyGuild         uint32
	Type            byte
	IsAcademyMember byte
	IsAutoPick      byte
	Level           byte
	Animation       [32]byte
	MoverHandle     uint16
	MoverSlot       byte
	KQTeamType      byte
	UsingMinipet    byte
	Unk             byte
}

// struct PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD
type NcBriefInfoChangeWeaponCmd struct {
	UpgradeInfo      NcBriefInfoChangeUpgradeCmd
	CurrentMobID     uint16
	CurrentKillLevel byte
}

// struct PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD
type NcBriefInfoDeleteHandleCmd struct {
	Handle uint16
}

// struct PROTO_NC_BRIEFINFO_INFORM_CMD
// nMyHnd is the affected client that received a server command involving hnd that was previously tagged as out of range.
type NcBriefInfoInformCmd struct {
	AffectedHandle  uint16
	ReceivedCommand NetCommand
	ForeignHandle   uint16
}
