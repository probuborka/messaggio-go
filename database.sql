CREATE TABLE IF NOT EXISTS public.messages
(
    id character varying(36) COLLATE pg_catalog."default" NOT NULL,
    message text COLLATE pg_catalog."default",
    processed boolean NOT NULL DEFAULT 'false',
    date_create timestamp without time zone NOT NULL DEFAULT now(),
    date_processed_start timestamp without time zone DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    date_processed_end timestamp without time zone DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    CONSTRAINT messages_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.messages
    OWNER to postgres;
