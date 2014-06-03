-- +goose Up
-- sql in section 'up' is executed when this migration is applied
create table replays
(
  replay_id serial not null,
  data json not null,
  constraint replays_pkey primary key (replay_id)
);

create index replays_id_index
  on replays (replay_id);

-- +goose Down
-- sql section 'down' is executed when this migration is rolled back
drop table replays;
