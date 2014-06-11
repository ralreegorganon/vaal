
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table matches
(
  match_id serial not null,
  match uuid not null,
  constraint matches_pkey primary key (match_id)
);

create index match_id_index
  on matches (match_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table match_id;
