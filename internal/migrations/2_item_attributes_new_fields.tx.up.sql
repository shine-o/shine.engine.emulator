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
    add column pattack_base smallint;

alter table world.item_attributes
    add column pattack_extra smallint;

--gopg:split
alter table world.item_attributes
    add column mattack_base smallint;

alter table world.item_attributes
    add column mattack_extra smallint;

--gopg:split
alter table world.item_attributes
    add column mdefense_base smallint;

alter table world.item_attributes
    add column mdefense_extra smallint;

--gopg:split
alter table world.item_attributes
    add column pdefense_base smallint;

alter table world.item_attributes
    add column pdefense_extra smallint;

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