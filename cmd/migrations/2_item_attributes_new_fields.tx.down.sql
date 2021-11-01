SET statement_timeout = 60000;
SET lock_timeout = 60000;

--gopg:split
alter table world.item_attributes
    drop column if exists dexterity_base;

alter table world.item_attributes
    drop column if exists dexterity_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists intelligence_base;

alter table world.item_attributes
    drop column if exists intelligence_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists endurance_base;

alter table world.item_attributes
    drop column if exists endurance_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists spirit_base;

alter table world.item_attributes
    drop column if exists spirit_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists pattack_base;

alter table world.item_attributes
    drop column if exists pattack_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists mattack_base;

alter table world.item_attributes
    drop column if exists mattack_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists mdefense_base;

alter table world.item_attributes
    drop column if exists mdefense_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists pdefense_base;

alter table world.item_attributes
    drop column if exists pdefense_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists aim_base;

alter table world.item_attributes
    drop column if exists aim_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists evasion_base;

alter table world.item_attributes
    drop column if exists evasion_extra;

--gopg:split
alter table world.item_attributes
    drop column if exists max_hp_base;

alter table world.item_attributes
    drop column if exists max_hp_extra;