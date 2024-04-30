package usersdb

// Column     |            Type             | Collation | Nullable |              Default
// ---------------+-----------------------------+-----------+----------+-----------------------------------
// id            | integer                     |           | not null | nextval('users_id_seq'::regclass)
// username      | character varying(50)       |           | not null |
// email         | character varying(100)      |           | not null |
// password_hash | character varying(100)      |           | not null |
// created_at    | timestamp without time zone |           |          | CURRENT_TIMESTAMP
// Indexes:
// "users_pkey" PRIMARY KEY, btree (id)
