package structs

//struct PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage
//{
//  unsigned __int16 handle;
//  PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage::<unnamed-type-flag> flag;
//  unsigned int hpchange;
//  unsigned int resthp;
//  unsigned __int16 hpchangeorder;
//};
type NcBatSkillBashHitDamageCmdSkillDamage struct {
	Handle uint16
	Flag NcBatSkillBashHitDamageCmdSkillDamageFlag
	HpChange uint32
	RestHP uint32
	HpChangeOrder uint16
}

//struct PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage::<unnamed-type-flag>
//{
//  _BYTE gap0[1];
//  char _bf1;
//};
type NcBatSkillBashHitDamageCmdSkillDamageFlag struct {
	Gap byte
	BF1 byte
}