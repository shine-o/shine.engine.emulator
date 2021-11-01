SET statement_timeout = 60000;
SET lock_timeout = 60000;

--gopg:split
alter table world.character_appearance
    drop constraint character_id_fk;

--gopg:split
alter table world.character_attributes
    drop constraint character_id_fk;

--gopg:split
alter table world.character_equipped_items
    drop constraint character_id_fk;

--gopg:split
alter table world.character_location
    drop constraint character_id_fk;

--gopg:split
alter table world.client_options
    drop constraint character_id_fk;

--gopg:split
alter table world.items
    drop constraint character_id_fk;

--gopg:split
alter table world.item_attributes
    drop constraint item_id_fk;

--gopg:split
alter table world.item_licences
    drop constraint item_id_fk;

--gopg:split
alter table world.item_enchantments
    drop constraint item_id_fk;