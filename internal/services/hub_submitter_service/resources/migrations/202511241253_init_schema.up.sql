CREATE TABLE IF NOT EXISTS inbox
(
    message_id     uuid        NOT NULL,
    correlation_id uuid,
    message_key    varchar        NOT NULL,
    routing_key    varchar     NOT NULL,
    received_at    timestamptz NOT NULL,
    processed_at   timestamptz,
    payload        JSONB       NOT NULL,
    raw_message     JSONB       NOT NULL,
    unique (message_id, routing_key),
    unique (message_key)
);