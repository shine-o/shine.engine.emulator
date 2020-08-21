package structs

//struct PROTO_SKILL_ITEMACTIONCOOLTIME_CMD::<unnamed-type-group>
//{
//  unsigned __int16 ItemActionID;
//  unsigned __int16 ItemActionGroupID;
//  unsigned int SecondCoolTime;
//};
type SkillItemActionCoolTimeCmdGroup struct {
	ItemActionID      uint16
	ItemActionGroupID uint16
	SecondCoolTime    uint32
}
