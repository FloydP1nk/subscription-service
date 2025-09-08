CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
    );

-- Индексы
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id
    ON subscriptions(user_id);

CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date
    ON subscriptions(start_date);

CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date
    ON subscriptions(end_date);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_period
    ON subscriptions(user_id, start_date, end_date);