-- pgREST will read from this schema alone
CREATE schema api;

create role web_anon nologin;
grant usage on schema api to web_anon;
-- web_anon can only read from tables in the api schema
grant SELECT ON ALL TABLES IN SCHEMA api TO web_anon;
-- alter the api schema to allow web_anon reads also on future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA api GRANT SELECT ON TABLES TO web_anon;

create role olm_dev nologin;
grant usage on schema api to olm_dev;
-- allow olm_dev to read and write to tables in the api schema
grant SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA api TO olm_dev;
-- alter the api schema to allow olm_dev rights also on future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA api GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO olm_dev;

-- authenticator role is used by postgrest to authenticate users
-- allow authenticator to switch into web_anon and olm_dev
create role authenticator noinherit;
grant web_anon to authenticator;
grant olm_dev to authenticator;


--set up citext extension
CREATE EXTENSION IF NOT EXISTS citext;