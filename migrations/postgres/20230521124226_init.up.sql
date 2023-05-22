CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS routers
(
    id uuid NOT NULL default uuid_generate_v4(),
    user_id uuid NOT NULL,
    name character varying(45) NOT NULL,
    address character varying(45) NOT NULL,
    username character varying(45) NOT NULL,
    password character varying(45) NOT NULL,
    lease_period_check bigint NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS user_id_idx
    ON routers USING btree
        (user_id ASC NULLS LAST);

CREATE TABLE IF NOT EXISTS logs (
    id uuid NOT NULL default uuid_generate_v4(),
    router_id uuid NOT NULL,
    time timestamp without time zone NOT NULL DEFAULT now(),
    level log_level NOT NULL,
    message text NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (router_id) REFERENCES routers(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS hosts (
    id uuid NOT NULL default uuid_generate_v4(),
    router_id uuid NOT NULL,
    name character varying(45) NOT NULL,
    address character varying(45),
    mac_address character varying(45),
    host_name character varying(45),
    last_online timestamp without time zone NOT NULL DEFAULT now(),
    is_online bool NOT NULL DEFAULT false,
    online_timeout bigint NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    FOREIGN KEY (router_id) REFERENCES routers(id) ON UPDATE CASCADE ON DELETE CASCADE
);