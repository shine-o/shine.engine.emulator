package zone

import "testing"

func TestNewPlayer_DefaultSkills(t *testing.T) {

}

func TestNewPlayer_Mage_DefaultSkills(t *testing.T) {

}

func TestPlayer_LearnSkill_Ok(t *testing.T) {

}

func TestPlayer_LearnSkill_Ok_NC(t *testing.T) {

}

func TestPlayer_LearnSkill_BadID(t *testing.T) {

}

func TestPlayer_LearnSkill_BadIndex(t *testing.T) {

}

func TestPlayer_LearnSkill_WrongClass(t *testing.T) {

}

func TestPlayer_LearnSkill_WrongClass_NC(t *testing.T) {

}

func TestPlayer_LearnSkill_WrongLevel(t *testing.T) {

}

func TestPlayer_LearnSkill_WrongLevel_NC(t *testing.T) {

}

func TestPlayer_LearnSkill_MissingRequirement(t *testing.T) {

}

func TestPlayer_LearnSkill_MissingRequirement_NC(t *testing.T) {

}

func TestPlayer_LearnSkill_AlreadyLearned(t *testing.T) {

}

func TestPlayer_LearnSkill_AlreadyLearned_NC(t *testing.T) {

}

func TestPlayer_UseSkill_Ok(t *testing.T) {

}

func TestPlayer_UseSkill_BadID(t *testing.T) {

}

func TestPlayer_UseSkill_BadIndex(t *testing.T) {

}

func TestPlayer_UseSkill_WrongClass(t *testing.T) {

}

func TestPlayer_UseSkill_WrongWeapon(t *testing.T) {

}

func TestPlayer_UseSkill_WrongWeapon_NC(t *testing.T) {

}

func TestPlayer_UseSkill_WrongLevel(t *testing.T) {

}

func TestPlayer_UseSkill_WrongLevel_NC(t *testing.T) {

}

func TestPlayer_UseSkill_NoTarget(t *testing.T) {

}

func TestPlayer_UseSkill_WeaponEquipped(t *testing.T) {
	// some skills work without a weapon
}

func TestPlayer_UseSkill_WeaponNotEquipped(t *testing.T) {

}

func TestPlayer_UseSkill_NotLearned(t *testing.T) {

}

func TestPlayer_UseSkill_OnCoolDown(t *testing.T) {
	// cast skill
	// try to cast the same skill immediately, should not work
}

func TestPlayer_UseSkill_Moving_Ok(t *testing.T) {
	// test a skill that can be used while the player has movements pending
}

func TestPlayer_UseSkill_Moving_Not_Ok(t *testing.T) {
	// test a skill that can NOT be used while the player has movements pending
	// pending movements should be canceled, skill casts anyway
}

func TestPlayer_UseSkill_BadTarget_NPC(t *testing.T) {
	// test a skill that can NOT be used while the player has movements pending
	// pending movements should be canceled, skill casts anyway
}

func TestPlayer_UseSkill_UnableToUseSkills(t *testing.T) {
	// abstate that prevents user from casting skill
}
