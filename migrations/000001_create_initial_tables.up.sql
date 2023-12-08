create table if not exists user_role
(
    id        int primary key,
    role_name text unique
);

create table if not exists "user"
(
    id       bigserial primary key,
    login    text unique,
    password text,
    role_id  int references user_role (id)
);

create table if not exists resource
(
    id                bigserial primary key,
    method            text,
    resource_name     text,
    allowed_roles_ids jsonb
);

insert into user_role (id, role_name)
values (0, 'all_message_sender'),
       (1, 'text_sender'),
       (2, 'file_sender');

insert into resource (method, resource_name, allowed_roles_ids)
values ('POST', '/v1/messages', '[
  0,
  1
]'),
       ('POST', '/v1/files', '[
         0,
         2
       ]');

insert into "user" (login, password, role_id)
values ('ronaldo', 'password1', 0),
       ('messi', 'password1', 1),
       ('pele', 'password1', 2);

create or replace procedure check_user_role(_resource_name text, _resource_method text, _role_id int)
    language plpgsql
as
$$
begin
    if not exists(select from resource r where r.resource_name = _resource_name and r.method = _resource_method) then
        raise exception 'resource not found';
    end if;

    if not exists(select
                  from resource r
                  where r.resource_name = _resource_name
                    and r.method = _resource_method
                    and _role_id in (select jsonb_array_elements(r.allowed_roles_ids)::int)) then
        raise exception 'not allowed role';
    end if;
end;
$$;