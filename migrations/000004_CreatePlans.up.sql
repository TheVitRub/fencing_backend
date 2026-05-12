-- Учебные планы клуба на определённые периоды.
-- Items хранит JSON-массив строк с пунктами плана,
-- например: ["Рапира XVI в.", "Дага", "Двуручный меч"].
CREATE TABLE IF NOT EXISTS plans (
    id          BIGSERIAL    PRIMARY KEY,
    title       TEXT         NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    period      VARCHAR(100) NOT NULL DEFAULT '',
    items       TEXT         NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  plans             IS 'Учебные планы клуба';
COMMENT ON COLUMN plans.period      IS 'Период плана, например "Осень 2025"';
COMMENT ON COLUMN plans.items       IS 'JSON-массив строк с пунктами плана';
