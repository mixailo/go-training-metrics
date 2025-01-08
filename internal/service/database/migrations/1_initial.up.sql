CREATE TABLE IF NOT EXISTS counter_storage
(
    id            serial PRIMARY KEY,
    value_name    VARCHAR(64) UNIQUE NOT NULL,
    counter_value BIGINT

);

CREATE TABLE IF NOT EXISTS gauge_storage
(
    id          serial PRIMARY KEY,
    value_name  VARCHAR(64) UNIQUE NOT NULL,
    gauge_value DOUBLE PRECISION

);


CREATE INDEX IF NOT EXISTS counter_value_name
    ON public.counter_storage USING btree
        (value_name COLLATE pg_catalog."default" varchar_ops ASC NULLS LAST)
    INCLUDE (value_name)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;

CREATE INDEX IF NOT EXISTS gauge_value_name
    ON public.gauge_storage USING btree
        (value_name COLLATE pg_catalog."default" varchar_ops ASC NULLS LAST)
    INCLUDE (value_name)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;
