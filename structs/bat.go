package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_BAT_ABSTATERESET_CMD
//{
//  unsigned __int16 handle;
//  ABSTATEINDEX abstate;
//};
type NcBatAbstateResetCmd struct {
	Handle       uint16
	AbstateIndex uint32
}

func (nc *NcBatAbstateResetCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatAbstateResetCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_ABSTATERESET_CMD
	{
	  unsigned __int16 handle;
	  ABSTATEINDEX abstate;
	};
`
}

//struct PROTO_NC_BAT_SPCHANGE_CMD
//{
//  unsigned int sp;
//};
type NcBatSpChangeCmd struct {
	SP uint32
}

func (nc *NcBatSpChangeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatSpChangeCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_SPCHANGE_CMD
	{
	  unsigned int sp;
	};
`
}

//struct PROTO_NC_BAT_LPCHANGE_CMD
//{
//  unsigned int nLP;
//};
type NcBatLpChangeCmd struct {
	LP uint32
}

func (nc *NcBatLpChangeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatLpChangeCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_LPCHANGE_CMD
	{
	  unsigned int nLP;
	};
`
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

func (nc *NcBatAbstateInformNoEffectCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatAbstateInformNoEffectCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_ABSTATEINFORM_NOEFFECT_CMD
	{
	  ABSTATEINDEX abstate;
	  unsigned int keeptime_millisec;
	};
`
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

func (nc *NcBatHpChangeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatHpChangeCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_HPCHANGE_CMD
	{
	  unsigned int hp;
	  unsigned __int16 hpchangeorder;
	};
`
}

//struct PROTO_NC_BAT_CEASE_FIRE_CMD
//{
//  unsigned __int16 handle;
//};
type NcBatCeaseFireCmd struct {
	Handle uint16
}

func (nc *NcBatCeaseFireCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatCeaseFireCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_CEASE_FIRE_CMD
	{
	  unsigned __int16 handle;
	};
`
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

func (nc *NcBatAbstateInformCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatAbstateInformCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_ABSTATEINFORM_CMD
	{
	  ABSTATEINDEX abstate;
	  unsigned int keeptime_millisec;
	};
`
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

func (nc *NcBatSkillBashObjCastReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatSkillBashObjCastReq) PdbType() string {
	return `
	struct PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ
	{
	  unsigned __int16 skill;
	  unsigned __int16 target;
	};
`
}

//struct PROTO_NC_BAT_ABSTATESET_CMD
//{
//  unsigned __int16 handle;
//  ABSTATEINDEX abstate;
//};
type NcBatAbstateSetCmd struct {
	Handle uint16
	Abstate uint32
}

func (nc * NcBatAbstateSetCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc * NcBatAbstateSetCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_ABSTATESET_CMD
	{
	  unsigned __int16 handle;
	  ABSTATEINDEX abstate;
	};
`
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
	Object uint16
	RestHP uint32
	Damage uint16
	Abstate uint16
	HPChangeOrder uint16
	IsMissDamage byte
}

func (nc * NcBatDotDamageCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc * NcBatDotDamageCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_DOTDAMAGE_CMD
	{
	  unsigned __int16 object;
	  unsigned int resthp;
	  unsigned __int16 damage;
	  unsigned __int16 abstate;
	  unsigned __int16 hpchangeorder;
	  char IsMissDamage;
	};
`
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
	Order byte
	Handle uint16
	TargetHP uint32
	TargetSP uint32
	TargetMaxSP uint32
	TargetLP uint16
	TargetMaxLP uint32
	TargetLevel byte
	HpChangeOrder uint16
}

func (nc * NcBatTargetInfoCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc * NcBatTargetInfoCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_TARGETINFO_CMD
	{
	  char order;
	  unsigned __int16 targethandle;
	  unsigned int targethp;
	  unsigned int targetmaxhp;
	  unsigned int targetsp;
	  unsigned int targetmaxsp;
	  unsigned int targetlp;
	  unsigned int targetmaxlp;
	  char targetlevel;
	  unsigned __int16 hpchangeorder;
	};
`
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
	Index uint16
	Caster uint16
	TargetNum byte
	SkillID uint16
	Target uint16
	Targets []NcBatSkillBashHitDamageCmdSkillDamage `struct:"sizefrom=TargetNum"`
}

func (nc * NcBatSkillBashHitDamageCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc * NcBatSkillBashHitDamageCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD
	{
	  unsigned __int16 index;
	  unsigned __int16 caster;
	  char targetnum;
	  unsigned __int16 kSkillID;
	  unsigned __int16 pTarget;
	  PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage target[];
	};
`
}

//struct PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD
//{
//  unsigned __int16 caster;
//  PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD castinfo;
//};
type NcBatSomeoneSkillBashHitObjStartCmd struct {
	Caster uint16
	CastInfo NcBatSkillBashHitObjStartCmd
}

func (nc *NcBatSomeoneSkillBashHitObjStartCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatSomeoneSkillBashHitObjStartCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD
	{
	  unsigned __int16 caster;
	  PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD castinfo;
	};
`
}

//struct PROTO_NC_BAT_SKILLBASH_HIT_OBJ_START_CMD
//{
//  unsigned __int16 skill;
//  unsigned __int16 targetobj;
//  unsigned __int16 index;
//};
type NcBatSkillBashHitObjStartCmd struct {
	Skill uint16
	TargetObj uint16
	Index uint16
}

//struct PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD
//{
//  unsigned __int16 index;
//  unsigned __int16 caster;
//  unsigned __int16 nSkillID;
//};
type NcBatSkillBashHitBlastCmd struct {
	Index uint16
	Caster uint16
	SkillID uint16
}

func (nc *NcBatSkillBashHitBlastCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatSkillBashHitBlastCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_SKILLBASH_HIT_BLAST_CMD
	{
	  unsigned __int16 index;
	  unsigned __int16 caster;
	  unsigned __int16 nSkillID;
	};
`
}