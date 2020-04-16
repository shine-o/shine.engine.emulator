package structs

//struct PROTO_NC_BAT_ABSTATERESET_CMD
//{
//  unsigned __int16 handle;
//  ABSTATEINDEX abstate;
//};
type NcBatAbstateResetCmd struct {
	Handle       uint16
	AbstateIndex uint32
}

//struct PROTO_NC_BAT_SPCHANGE_CMD
//{
//  unsigned int sp;
//};
type NcBatSpChangeCmd struct {
	SP uint32
}

//struct PROTO_NC_BAT_LPCHANGE_CMD
//{
//  unsigned int nLP;
//};
type NcBatLpChangeCmd struct {
	LP uint32
}

//struct PROTO_NC_BAT_ABSTATEINFORM_NOEFFECT_CMD
//{
//  ABSTATEINDEX abstate;
//  unsigned int keeptime_millisec;
//};
type NcBatAbstateInformNoEffectCmd struct {
	Abstate          uint32
	KeepTimeMillisec uint32
}

//struct PROTO_NC_BAT_HPCHANGE_CMD
//{
//  unsigned int hp;
//  unsigned __int16 hpchangeorder;
//};
type NcBatHpChangeCmd struct {
	HP            uint32
	HpChangeOrder uint16
}

//struct PROTO_NC_BAT_CEASE_FIRE_CMD
//{
//  unsigned __int16 handle;
//};
type NcBatCeaseFireCmd struct {
	Handle uint16
}

//struct PROTO_NC_BAT_ABSTATEINFORM_CMD
//{
//  ABSTATEINDEX abstate;
//  unsigned int keeptime_millisec;
//};
type NcBatAbstateInformCmd struct {
	Abstate          uint32
	KeepTimeMillisec uint32
}

//struct PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ
//{
//  unsigned __int16 skill;
//  unsigned __int16 target;
//};
type NcBatSkillBashObjCastReq struct {
	Skill  uint16
	Target uint16
}

//struct PROTO_NC_BAT_ABSTATESET_CMD
//{
//  unsigned __int16 handle;
//  ABSTATEINDEX abstate;
//};
type NcBatAbstateSetCmd struct {
	Handle  uint16
	Abstate uint32
}

//struct PROTO_NC_BAT_DOTDAMAGE_CMD
//{
//  unsigned __int16 object;
//  unsigned int resthp;
//  unsigned __int16 damage;
//  unsigned __int16 abstate;
//  unsigned __int16 hpchangeorder;
//  char IsMissDamage;
//};
type NcBatDotDamageCmd struct {
	Object        uint16
	RestHP        uint32
	Damage        uint16
	Abstate       uint16
	HPChangeOrder uint16
	IsMissDamage  byte
}

//struct PROTO_NC_BAT_TARGETINFO_CMD
//{
//  char order;
//  unsigned __int16 targethandle;
//  unsigned int targethp;
//  unsigned int targetmaxhp;
//  unsigned int targetsp;
//  unsigned int targetmaxsp;
//  unsigned int targetlp;
//  unsigned int targetmaxlp;
//  char targetlevel;
//  unsigned __int16 hpchangeorder;
//};
type NcBatTargetInfoCmd struct {
	Order         byte
	Handle        uint16
	TargetHP      uint32
	TargetSP      uint32
	TargetMaxSP   uint32
	TargetLP      uint16
	TargetMaxLP   uint32
	TargetLevel   byte
	HpChangeOrder uint16
}

//struct PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD
//{
//  unsigned __int16 index;
//  unsigned __int16 caster;
//  char targetnum;
//  unsigned __int16 kSkillID;
//  unsigned __int16 pTarget;
//  PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage target[];
//};
type NcBatSkillBashHitDamageCmd struct {
	Index     uint16
	Caster    uint16
	TargetNum byte
	SkillID   uint16
	Target    uint16
	Targets   []NcBatSkillBashHitDamageCmdSkillDamage `struct:"sizefrom=TargetNum"`
}

//struct PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD
//{
//  unsigned __int16 caster;
//  PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD castinfo;
//};
type NcBatSomeoneSkillBashHitObjStartCmd struct {
	Caster   uint16
	CastInfo NcBatSkillBashHitObjStartCmd
}

//struct PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD
//{
//  unsigned __int16 skill;
//  unsigned __int16 targetobj;
//  unsigned __int16 index;
//};
type NcBatSkillBashHitObjStartCmd struct {
	Skill     uint16
	TargetObj uint16
	Index     uint16
}

//struct PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD
//{
//  unsigned __int16 index;
//  unsigned __int16 caster;
//  unsigned __int16 nSkillID;
//};
type NcBatSkillBashHitBlastCmd struct {
	Index   uint16
	Caster  uint16
	SkillID uint16
}
