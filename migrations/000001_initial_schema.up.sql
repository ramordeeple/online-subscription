CREATE TABLE subscriptions
(
    id            UUID PRIMARY KEY,
    service_name  TEXT        NOT NULL,
    monthly_price INT         NOT NULL,
    user_id       UUID        NOT NULL,
    start_date    DATE        NOT NULL,
    end_date      DATE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (monthly_price > 0)
);

CREATE INDEX idx_subscriptions_user_id
    ON subscriptions (user_id);

CREATE INDEX idx_subscriptions_service_name
    ON subscriptions (service_name);

CREATE INDEX idx_subscriptions_start_date
    ON subscriptions (start_date);

CREATE INDEX idx_subscriptions_end_date
    ON subscriptions (end_date);
