
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table matches add column open boolean not null default true;
alter table matches add column complete boolean not null default false;
alter table matches add column replay_id integer;
alter table matches add foreign key (replay_id) references replays (replay_id) on update no action on delete no action;



-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table matches drop column open;
alter table matches drop column complete;
alter table matches drop column replay_id;
