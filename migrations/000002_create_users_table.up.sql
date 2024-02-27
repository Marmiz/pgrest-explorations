CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.users (
    email citext UNIQUE PRIMARY KEY NOT NULL CHECK (email ~* '^.+@.+\..+$'),
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    password_hash bytea NOT NULL check (length(password_hash) < 512),
    active bool NOT NULL DEFAULT false,
    role name not null check (length(role) < 512)
);

CREATE OR REPLACE function
auth.check_role_exists() returns trigger as $$
begin
  if not exists (select 1 from pg_roles as r where r.rolname = new.role) then
    raise foreign_key_violation using message =
      'unknown database role: ' || new.role;
    return null;
  end if;
  return new;
end
$$ language plpgsql;

DROP TRIGGER IF EXISTS ensure_user_role_exists on auth.users;

CREATE CONSTRAINT trigger ensure_user_role_exists
  after insert or update on auth.users
  for each row
  execute procedure auth.check_role_exists();

-- seed a user for development
-- password is hashed and generated prior.
-- update with your own password hash
INSERT INTO auth.users (email, password_hash, active, role) VALUES
('local_dev@localhost.com', '$2a$12$rpj0BrvVdma7upkr3pRuZ.TxsX26JzL2JLB1sry0p6bu1eOouw0E.', true, 'dev')