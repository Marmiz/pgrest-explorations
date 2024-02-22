DROP schema IF EXISTS api CASCADE;

DROP role IF EXISTS web_anon;
DROP role IF EXISTS olm_dev;
DROP role IF EXISTS authenticator;

DROP EXTENSION IF EXISTS citext;

DROP table IF EXISTS api.todos;