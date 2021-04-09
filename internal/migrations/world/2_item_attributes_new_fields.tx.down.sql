SET statement_timeout = 60000;
SET lock_timeout = 60000;

--gopg:split
alter table world.item_attributes
    drop column dexterity_base;

alter table world.item_attributes
    drop column dexterity_extra;

--gopg:split
alter table world.item_attributes
    drop column intelligence_base;

alter table world.item_attributes
    drop column intelligence_extra;

--gopg:split
alter table world.item_attributes
    drop column endurance_base;

alter table world.item_attributes
    drop column endurance_extra;

--gopg:split
alter table world.item_attributes
    drop column spirit_base;

alter table world.item_attributes
    drop column spirit_extra;

--gopg:split
alter table world.item_attributes
    drop column pattack_base;

alter table world.item_attributes
    drop column pattack_extra;

--gopg:split
alter table world.item_attributes
    drop column mattack_base;

alter table world.item_attributes
    drop column mattack_extra;

--gopg:split
alter table world.item_attributes
    drop column mdefense_base;

alter table world.item_attributes
    drop column mdefense_extra;

--gopg:split
alter table world.item_attributes
    drop column pdefense_base;

alter table world.item_attributes
    drop column pdefense_extra;

--gopg:split
alter table world.item_attributes
    drop column aim_base;

alter table world.item_attributes
    drop column aim_extra;

--gopg:split
alter table world.item_attributes
    drop column evasion_base;

alter table world.item_attributes
    drop column evasion_extra;

--gopg:split
alter table world.item_attributes
    drop column max_hp_base;

alter table world.item_attributes
    drop column max_hp_extra;