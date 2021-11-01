package structs

// struct PROTO_NC_BAT_ABSTATERESET_CMD
type NcBatAbstateResetCmd struct {
	Handle       uint16
	AbstateIndex uint32
}

// struct PROTO_NC_BAT_SPCHANGE_CMD
type NcBatSpChangeCmd struct {
	SP uint32
}

// struct PROTO_NC_BAT_LPCHANGE_CMD
type NcBatLpChangeCmd struct {
	LP uint32
}

// struct PROTO_NC_BAT_ABSTATEINFORM_NOEFFECT_CMD
type NcBatAbstateInformNoEffectCmd struct {
	Abstate          uint32
	KeepTimeMillisec uint32
}

// struct PROTO_NC_BAT_HPCHANGE_CMD
type NcBatHpChangeCmd struct {
	HP            uint32
	HpChangeOrder uint16
}

// struct PROTO_NC_BAT_CEASE_FIRE_CMD
type NcBatCeaseFireCmd struct {
	Handle uint16
}

// struct PROTO_NC_BAT_ABSTATEINFORM_CMD
type NcBatAbstateInformCmd struct {
	Abstate          uint32
	KeepTimeMillisec uint32
}

// struct PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ
type NcBatSkillBashObjCastReq struct {
	Skill  uint16
	Target uint16
}

// struct PROTO_NC_BAT_ABSTATESET_CMD
type NcBatAbstateSetCmd struct {
	Handle  uint16
	Abstate uint32
}

// struct PROTO_NC_BAT_DOTDAMAGE_CMD
type NcBatDotDamageCmd struct {
	Object        uint16
	RestHP        uint32
	Damage        uint16
	Abstate       uint16
	HPChangeOrder uint16
	IsMissDamage  byte
}

type NcBatTargetInfoReq struct {
	TargetHandle uint16
}

type NcBatUnTargetReq struct{}

// struct PROTO_NC_BAT_TARGETINFO_CMD
type NcBatTargetInfoCmd struct {
	Order         byte
	Handle        uint16
	TargetHP      uint32
	TargetMaxHP   uint32
	TargetSP      uint32
	TargetMaxSP   uint32
	TargetLP      uint32
	TargetMaxLP   uint32
	TargetLevel   byte
	HpChangeOrder uint16
}

// struct PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD
type NcBatSkillBashHitDamageCmd struct {
	Index     uint16
	Caster    uint16
	TargetNum byte
	SkillID   uint16
	Target    uint16
	Targets   []NcBatSkillBashHitDamageCmdSkillDamage `struct:"sizefrom=TargetNum"`
}

// struct PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD
type NcBatSomeoneSkillBashHitObjStartCmd struct {
	Caster   uint16
	CastInfo NcBatSkillBashHitObjStartCmd
}

// struct PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD
type NcBatSkillBashHitObjStartCmd struct {
	Skill     uint16
	TargetObj uint16
	Index     uint16
	Unk       uint32
}

// struct PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD
type NcBatSkillBashHitBlastCmd struct {
	Index   uint16
	Caster  uint16
	SkillID uint16
}

// struct PROTO_NC_BAT_SWING_START_CMD
type NcBatSwingStartCmd struct {
	Attacker       uint16
	Defender       uint16
	ActionCode     byte
	AttackSpeed    uint16
	DamageIndex    byte
	AttackSequence byte
}

// struct PROTO_NC_BAT_SWING_DAMAGE_CMD::<unnamed-type-flag>
type NcBatSwingDamageCmdFlag struct {
	Gap [1]byte
	BF1 byte
}

// struct PROTO_NC_BAT_SWING_DAMAGE_CMD
type NcBatSwingDamageCmd struct {
	Attacker       uint16
	Defender       uint16
	Flag           NcBatSwingDamageCmdFlag
	Damage         uint16
	RestHP         uint32
	HpChangeOrder  uint16
	DamageIndex    byte
	AttackSequence byte
}
