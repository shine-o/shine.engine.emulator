package structs

//struct PROTO_SKILL_ITEMACTIONCOOLTIME_CMD
//{
//  unsigned __int16 num;
//  PROTO_SKILL_ITEMACTIONCOOLTIME_CMD::<unnamed-type-group> group[];
//};
type SkillItemActionCoolTimeCmd struct {
	Num    uint16
	Groups []SkillItemActionCoolTimeCmdGroup `struct:"sizefrom=Num"`
}
