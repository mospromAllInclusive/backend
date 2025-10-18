create table if not exists app.changelog
(
    change_id  bigserial primary key,
    target     text                     not null,
    user_id    integer                  not null,
    table_id   text,
    column_id  text,
    row_id     bigint,
    change     jsonb                    not null default '{}',
    changed_at timestamp with time zone not null
);

create index if not exists changelog_target_table_column_row_idx on app.changelog (target, table_id, column_id, row_id);
create index if not exists changelog_changed_at_idx on app.changelog (changed_at);
