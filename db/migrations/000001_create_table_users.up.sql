create extension if not exists pgcrypto;

create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    role int not null default 2,
    email varchar(255) not null,
    hash_password text default '',
    is_verified boolean not null default false,
    is_oauth boolean not null default false,
    is_active boolean not null default true,
    created_at timestamp not null default current_timestamp,
    created_by uuid not null,
    updated_at timestamp default null,
    updated_by uuid default null,
    deleted_at timestamp default null,
    deleted_by varchar(255) default null,
    deleted_reason varchar(255) default null
);

create table if not exists user_details (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    full_name varchar(255) not null,
    sex smallint not null default 0,
    phone_number varchar(255) default null,
    image_url text not null,
    birth_date date not null,
    created_at timestamp not null default current_timestamp,
    created_by uuid not null,
    updated_at timestamp default null,
    updated_by uuid default null,
    deleted_at timestamp default null,
    deleted_by varchar(255) default null,
    deleted_reason varchar(255) default null
);

create table if not exists token_verification_users (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    verification_token uuid not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp default null,
    deleted_at timestamp default null
);

create table if not exists token_reset_users (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    reset_token uuid not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp default null,
    deleted_at timestamp default null
);

comment on column user_details.sex is
'0 = other, 1 = female, 2 = male';

comment on column users.role is
'0 = unknown, 1 = admin, 2 = user';

create index if not exists idx_users_role on users (role);
create index if not exists idx_user_details_user_id_not_deleted on user_details (user_id) where deleted_at is null;
create index if not exists idx_user_details_birth_date_not_deleted on user_details (birth_date) where deleted_at is null;
create index if not exists idx_token_verification_users_token_not_deleted on token_verification_users (verification_token) where deleted_at is null;
create index if not exists idx_token_verification_users_user_id_not_deleted on token_verification_users (user_id) where deleted_at is null;
create index if not exists idx_token_reset_users_token_not_deleted on token_reset_users (reset_token) where deleted_at is null;
create index if not exists idx_token_reset_users_user_id_not_deleted on token_reset_users (user_id) where deleted_at is null;

create unique index if not exists uq_users_email on users (email) where deleted_at is null;
create unique index if not exists uq_user_details_user_id on user_details (user_id);
create unique index if not exists uq_user_details_phone_number on user_details (phone_number) where deleted_at is null and phone_number is not null;
