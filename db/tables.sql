-- Copyright (C) 2014 Miquel Sabaté Solà
-- This file is licensed under the MIT license.
-- See the LICENSE file.


create table if not exists users (
    id uuid primary key,
    name varchar(255) unique not null check (name <> ''),
    password_hash text,
    created_at timestamp
);

create table if not exists topics (
    id uuid primary key,
    name varchar(255) unique not null check (name <> ''),
    user_id uuid references users(id) on delete cascade,
    created_at timestamp
);

