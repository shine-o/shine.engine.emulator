SET statement_timeout = 60000;
SET lock_timeout = 60000;

--gopg:split
alter table world.item_attributes
    add column dexterity_base smallint;

alter table world.item_attributes
    add column dexterity_extra smallint;

--gopg:split
alter table world.item_attributes
    add column intelligence_base smallint;

alter table world.item_attributes
    add column intelligence_extra smallint;

--gopg:split
alter table world.item_attributes
    add column endurance_base smallint;

alter table world.item_attributes
    add column endurance_extra smallint;

--gopg:split
alter table world.item_attributes
    add column spirit_base smallint;

alter table world.item_attributes
    add column spirit_extra smallint;

--gopg:split
alter table world.item_attributes
    add column p_attack_base smallint;

alter table world.item_attributes
    add column p_attack_extra smallint;

--gopg:split
alter table world.item_attributes
    add column m_attack_base smallint;

alter table world.item_attributes
    add column m_attack_extra smallint;

--gopg:split
alter table world.item_attributes
    add column m_defense_base smallint;

alter table world.item_attributes
    add column m_defense_extra smallint;

--gopg:split
alter table world.item_attributes
    add column p_defense_base smallint;

alter table world.item_attributes
    add column p_defense_extra smallint;

--gopg:split
alter table world.item_attributes
    add column aim_base smallint;

alter table world.item_attributes
    add column aim_extra smallint;

--gopg:split
alter table world.item_attributes
    add column evasion_base smallint;

alter table world.item_attributes
    add column evasion_extra smallint;

--gopg:split
alter table world.item_attributes
    add column max_hp_base smallint;

alter table world.item_attributes
    add column max_hp_extra smallint;