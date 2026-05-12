-- События клуба: турниры, тренировки, показательные выступления.
CREATE TABLE IF NOT EXISTS events (
    id          BIGSERIAL    PRIMARY KEY,
    title       TEXT         NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    date        TIMESTAMPTZ  NOT NULL,
    location    TEXT         NOT NULL DEFAULT '',
    image_url   TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  events            IS 'События клуба (турниры, тренировки, выступления)';
COMMENT ON COLUMN events.title      IS 'Название события';
COMMENT ON COLUMN events.description IS 'Описание события';
COMMENT ON COLUMN events.date       IS 'Дата и время проведения события';
COMMENT ON COLUMN events.location   IS 'Место проведения события';
COMMENT ON COLUMN events.image_url  IS 'URL обложки события';

-- Индекс для сортировки по дате (самый частый запрос).
CREATE INDEX IF NOT EXISTS events_date_idx ON events (date DESC);
