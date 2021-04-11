package structs

//struct PROTO_SKILL_ITEMACTIONCOOLTIME_CMD
type SkillItemActionCoolTimeCmd struct {
	Num    uint16
	Groups []SkillItemActionCoolTimeCmdGroup `struct:"sizefrom=Num"`
}
