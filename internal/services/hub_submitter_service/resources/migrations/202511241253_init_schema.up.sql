CREATE TABLE IF NOT EXISTS inbox
(
    message_id   uuid        NOT NULL,
    correlation_id   uuid,
    message_key uuid NOT NULL,
    topic        varchar     NOT NULL,
    received_at  timestamptz NOT NULL,
    processed_at timestamptz,
    payload      JSONB       NOT NULL,
    unique (message_id, topic),
    unique(message_key)
);