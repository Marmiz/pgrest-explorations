-- pgREST will read from this schema alone
CREATE schema api;

create role web_anon nologin;
grant usage on schema api to web_anon;
-- web_anon can only read from tables in the api schema
grant SELECT ON ALL TABLES IN SCHEMA api TO web_anon;
-- alter the api schema to allow web_anon reads also on future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA api GRANT SELECT ON TABLES TO web_anon;

create role dev nologin;
grant usage on schema api to dev;
-- allow olm_dev to read and write to tables in the api schema
grant SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA api TO dev;
-- alter the api schema to allow dev rights also on future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA api GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO dev;

-- authenticator role is used by postgrest to authenticate users
-- allow authenticator to switch into web_anon and dev
create role authenticator noinherit;
grant web_anon to authenticator;
grant dev to authenticator;


--set up citext extension
CREATE EXTENSION IF NOT EXISTS citext;