SET statement_timeout = 60000;
SET lock_timeout = 60000;

--gopg:split
create table if not exists world.character_items
(
    id             bigserial not null
        constraint character_items_pkey
            primary key,
    inventory_type bigint    not null,
    slot           integer   not null,
    character_id   bigint    not null,
    shn_id         integer   not null,
    stackable      boolean   not null,
    amount         bigint,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    constraint character_items_inventory_type_slot_character_id_key
        unique (inventory_type, slot, character_id)
);