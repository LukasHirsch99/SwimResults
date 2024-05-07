CREATE TABLE test (
  id serial PRIMARY KEY,
  text varchar
);

CREATE ROLE anonymous NOLOGIN;
GRANT anonymous TO admin;
GRANT select ON ALL TABLES IN SCHEMA public TO anonymous;
