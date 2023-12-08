CREATE TABLE IF NOT EXISTS user_role
(
    id        int primary key,
    role_name text unique
);

CREATE TABLE IF NOT EXISTS "user"
(
    id       bigserial primary key,
    login    text unique,
    password text,
    role_id  int references user_role (id)
);

CREATE TABLE IF NOT EXISTS resource
(
    id                bigserial primary key,
    method            text,
    resource_name     text,
    allowed_roles_ids jsonb
);

INSERT INTO user_role (id, role_name)
VALUES (0, 'all_message_sender'),
       (1, 'text_sender'),
       (2, 'file_sender');

INSERT INTO resource (method, resource_name, allowed_roles_ids)
VALUES ('POST', '/v1/messages', '[
  0,
  1
]'),
       ('POST', '/v1/files', '[
         0,
         2
       ]');

INSERT INTO "user" (login, password, role_id)
VALUES ('ronaldo', 'password1', 0),
       ('messi', 'password1', 1),
       ('pele', 'password1', 2);

CREATE OR REPLACE PROCEDURE check_user_role(_resource_name text, _resource_method text, _role_id int)
    LANGUAGE plpgsql
AS
$$
BEGIN
    IF NOT EXISTS(SELECT FROM resource r WHERE r.resource_name = _resource_name AND r.method = _resource_method) THEN
        RAISE EXCEPTION 'resource not found';
    END IF;

    IF NOT EXISTS(SELECT
                  FROM resource r
                  WHERE r.resource_name = _resource_name
                    AND r.method = _resource_method
                    AND _role_id IN (SELECT jsonb_array_elements(r.allowed_roles_ids)::int)) THEN
        RAISE EXCEPTION 'not allowed role';
    END IF;
END;
$$;