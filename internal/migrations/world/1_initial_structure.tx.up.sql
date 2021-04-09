SET statement_timeout = 60000;
SET lock_timeout = 60000;

--gopg:split
create schema if not exists world

--gopg:split
create table if not exists world.item_enchantments
(
    id      bigserial not null
        constraint item_enchantments_pkey
            primary key,
    item_id bigint,
    item    jsonb
);

--gopg:split
create table if not exists world.item_licences
(
    id      bigserial not null
        constraint item_licences_pkey
            primary key,
    item_id bigint,
    item    jsonb,
    shn_id  integer
);

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

--gopg:split
create table if not exists world.characters
(
    id          bigserial not null
        constraint characters_pkey
            primary key,
    user_id     bigint    not null,
    name        text      not null
        constraint characters_name_key
            unique,
    admin_level smallint  not null,
    slot        smallint  not null,
    is_deleted  boolean,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
);

--gopg:split
create table if not exists world.character_appearance
(
    id           bigserial not null
        constraint character_appearance_pkey
            primary key,
    character_id bigint,
    class        smallint  not null,
    gender       smallint  not null,
    hair_type    smallint  not null,
    hair_color   smallint  not null,
    face_type    smallint  not null,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone
);

--gopg:split
create table if not exists world.character_attributes
(
    id           bigserial not null
        constraint character_attributes_pkey
            primary key,
    character_id bigint,
    level        smallint  not null,
    experience   bigint    not null,
    fame         bigint    not null,
    hp           bigint    not null,
    sp           bigint    not null,
    intelligence smallint  not null,
    strength     smallint  not null,
    dexterity    smallint  not null,
    endurance    smallint  not null,
    spirit       smallint  not null,
    money        bigint    not null,
    kill_points  bigint    not null,
    hp_stones    integer   not null,
    sp_stones    integer   not null,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone
);

--gopg:split
create table if not exists world.character_location
(
    id           bigserial not null
        constraint character_location_pkey
            primary key,
    character_id bigint,
    map_id       bigint    not null,
    map_name     text      not null,
    x            bigint    not null,
    y            bigint    not null,
    d            bigint    not null,
    is_kq        boolean   not null,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone
);

--gopg:split
create table if not exists world.client_options
(
    id           bigserial not null
        constraint client_options_pkey
            primary key,
    character_id bigint,
    game_options bytea     not null,
    keymap       bytea     not null,
    shortcuts    bytea     not null,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone
);

--gopg:split
create table if not exists world.character_equipped_items
(
    id                 bigserial not null
        constraint character_equipped_items_pkey
            primary key,
    character_id       bigint
        constraint character_equipped_items_character_id_fkey
            references world.characters,
    head               integer,
    face               integer,
    body               integer,
    pants              integer,
    boots              integer,
    left_hand          integer,
    right_hand         integer,
    left_mini_pet      integer,
    right_mini_pet     integer,
    apparel_head       integer,
    apparel_face       integer,
    apparel_eye        integer,
    apparel_body       integer,
    apparel_pants      integer,
    apparel_boots      integer,
    apparel_left_hand  integer,
    apparel_right_hand integer,
    apparel_back       integer,
    apparel_tail       integer,
    apparel_aura       integer,
    apparel_shield     integer,
    deleted_at         timestamp with time zone
);

--gopg:split
create table if not exists world.items
(
    id             bigserial not null
        constraint items_pkey
            primary key,
    inventory_type bigint    not null,
    slot           bigint    not null,
    character_id   bigint    not null,
    shn_id         integer   not null,
    shn_inx_name   text      not null,
    stackable      boolean   not null,
    amount         bigint,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    constraint items_inventory_type_slot_character_id_key
        unique (inventory_type, slot, character_id)
);

--gopg:split
create table if not exists world.item_attributes
(
    id             bigserial not null
        constraint item_attributes_pkey
            primary key,
    item_id        bigint    not null,
    strength_base  bigint,
    strength_extra bigint,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    constraint item_attributes_id_item_id_key
        unique (id, item_id)
);

