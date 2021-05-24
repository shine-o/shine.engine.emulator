package structs

// struct PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage
type NcBatSkillBashHitDamageCmdSkillDamage struct {
	Handle        uint16
	Flag          NcBatSkillBashHitDamageCmdSkillDamageFlag
	HpChange      uint32
	RestHP        uint32
	HpChangeOrder uint16
}

// struct PROTO_NC_BAT_SKILLBASH_HIT_DAMAGE_CMD::SkillDamage::<unnamed-type-flag>
type NcBatSkillBashHitDamageCmdSkillDamageFlag struct {
	Gap byte
	BF1 byte
}
