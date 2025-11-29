CREATE TABLE subscriptions
(
    id           UUID PRIMARY KEY,
    service_name TEXT NOT NULL,
    price        INT  NOT NULL,
    user_id      UUID NOT NULL,
    start_month  INT  NOT NULL,
    start_year   INT  NOT NULL,
    end_month    INT,
    end_year     INT
);
