CREATE TABLE snippets(id serial not null unique, title character varying(100), content text,
created timestamp without time zone, expires timestamp without time zone );
INSERT INTO snippets(title,content,created,expires)
VALUES ('Docker Example','There should be some content',CURRENT_DATE, CURRENT_DATE + INTERVAL '365 day')