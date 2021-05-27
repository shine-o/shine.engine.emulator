package zone

import "testing"

func TestNewPlayerDefaultSkills(t *testing.T) {
	t.Fail()
}

func TestNewPlayerMageDefaultSkills(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillOk(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillOkNC(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillBadID(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillBadIndex(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillWrongClass(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillWrongClassNC(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillWrongLevel(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillWrongLevelNC(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillMissingRequirement(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillMissingRequirementNC(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillAlreadyLearned(t *testing.T) {
	t.Fail()
}

func TestPlayerLearnSkillAlreadyLearnedNC(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillOk(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillBadID(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillBadIndex(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillWrongClass(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillWrongWeapon(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillWrongWeaponNC(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillWrongLevel(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillWrongLevelNC(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillNoTarget(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillWeaponEquipped(t *testing.T) {
	t.Fail()
	// some skills work without a weapon
}

func TestPlayerUseSkillWeaponNotEquipped(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillNotLearned(t *testing.T) {
	t.Fail()
}

func TestPlayerUseSkillOnCoolDown(t *testing.T) {
	t.Fail()

	// cast skill
	// try to cast the same skill immediately, should not work
}

func TestPlayerUseSkillMovingOk(t *testing.T) {
	t.Fail()

	// test a skill that can be used while the player has movements pending
}

func TestPlayerUseSkillMovingNotOk(t *testing.T) {
	t.Fail()

	// test a skill that can NOT be used while the player has movements pending
	// pending movements should be canceled, skill casts anyway
}

func TestPlayerUseSkillBadTargetNPC(t *testing.T) {
	t.Fail()

	// test a skill that can NOT be used while the player has movements pending
	// pending movements should be canceled, skill casts anyway
}

func TestPlayerUseSkillUnableToUseSkills(t *testing.T) {
	t.Fail()

	// abstate that prevents user from casting skill
}
