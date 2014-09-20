-- Copyright (C) 2014 Miquel Sabaté Solà
-- This file is licensed under the MIT license.
-- See the LICENSE file.


create table users (
    id uuid primary key,
    name varchar(255) unique not null check (name <> ''),
    password_hash text,
    created_at timestamp
);

