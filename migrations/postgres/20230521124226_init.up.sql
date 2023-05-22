CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS links
(
    id uuid NOT NULL default uuid_generate_v4(),
    user_id uuid NOT NULL,
    sst_email character varying(45) NOT NULL,
    sst_password character varying(45) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS user_id_idx
    ON links USING btree
        (user_id ASC NULLS LAST);

CREATE TABLE IF NOT EXISTS logs (
    id uuid NOT NULL default uuid_generate_v4(),
    link_id uuid NOT NULL,
    time timestamp without time zone NOT NULL DEFAULT now(),
    level log_level NOT NULL,
    message text NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (link_id) REFERENCES links(id) ON UPDATE CASCADE ON DELETE CASCADE
);
