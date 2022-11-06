DROP DATABASE IF EXISTS config_controller;
CREATE DATABASE config_controller;
\c config_controller;

CREATE TABLE config_controller.public.configs (
    id      SERIAL PRIMARY KEY,
    service varchar(15) NOT NULL
);

CREATE TABLE config_controller.public.data_configs (
    config_id   integer REFERENCES config_controller.public.configs (id),
    version     integer,
    data        JSON NOT NULL
);