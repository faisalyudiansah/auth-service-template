drop index if exists idx_users_role;
drop index if exists idx_user_details_birth_date_not_deleted;
drop index if exists idx_user_details_user_id_not_deleted;
drop index if exists idx_token_verification_users_user_id_not_deleted;
drop index if exists idx_token_verification_users_token_not_deleted;
drop index if exists idx_token_reset_users_user_id_not_deleted;
drop index if exists idx_token_reset_users_token_not_deleted;
drop index if exists uq_user_details_phone_number;
drop index if exists uq_user_details_user_id;
drop index if exists uq_users_email;

drop table if exists token_reset_users cascade;
drop table if exists token_verification_users cascade;
drop table if exists user_details cascade;
drop table if exists users cascade;

drop extension if exists pgcrypto;
