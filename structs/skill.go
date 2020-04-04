package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_SKILL_ITEMACTIONCOOLTIME_CMD
//{
//  unsigned __int16 num;
//  PROTO_SKILL_ITEMACTIONCOOLTIME_CMD::<unnamed-type-group> group[];
//};
type SkillItemActionCoolTimeCmd struct {
	Num    uint16
	Groups []SkillItemActionCoolTimeCmdGroup `struct:"sizefrom=Num"`
}

func (nc *SkillItemActionCoolTimeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *SkillItemActionCoolTimeCmd) PdbType() string {
	return `
	struct PROTO_SKILL_ITEMACTIONCOOLTIME_CMD
	{
	  unsigned __int16 num;
	  PROTO_SKILL_ITEMACTIONCOOLTIME_CMD::<unnamed-type-group> group[];
	};
`
}
