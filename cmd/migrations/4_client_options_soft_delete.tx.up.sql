SET statement_timeout = 60000;
SET lock_timeout = 60000;

--gopg:split
alter table world.client_options
    add column deleted_at     timestamp with time zone;
