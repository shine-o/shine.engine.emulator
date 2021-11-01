SET statement_timeout = 60000;
SET lock_timeout = 60000;

--gopg:split
alter table world.character_appearance
    add constraint character_id_fk
        foreign key (character_id)
            references world.characters (id);

--gopg:split
alter table world.character_attributes
    add constraint character_id_fk
        foreign key (character_id)
            references world.characters (id);

--gopg:split
alter table world.character_equipped_items
    add constraint character_id_fk
        foreign key (character_id)
            references world.characters (id);

--gopg:split
alter table world.character_location
    add constraint character_id_fk
        foreign key (character_id)
            references world.characters (id);

--gopg:split
alter table world.client_options
    add constraint character_id_fk
        foreign key (character_id)
            references world.characters (id);

--gopg:split
alter table world.items
    add constraint character_id_fk
        foreign key (character_id)
            references world.characters (id);

--gopg:split
alter table world.item_attributes
    add constraint item_id_fk
        foreign key (item_id)
            references world.items (id);

--gopg:split
alter table world.item_licences
    add constraint item_id_fk
        foreign key (item_id)
            references world.items (id);

--gopg:split
alter table world.item_enchantments
    add constraint item_id_fk
        foreign key (item_id)
            references world.items (id);