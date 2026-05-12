-- Информация об основателе школы.
-- В таблице всегда ровно одна запись с id=1.
-- Используется UpsertFounder (INSERT ON CONFLICT DO UPDATE).
CREATE TABLE IF NOT EXISTS founder (
    id         BIGSERIAL   PRIMARY KEY,
    name       TEXT        NOT NULL,
    bio        TEXT        NOT NULL DEFAULT '',
    photo_url  TEXT        NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  founder            IS 'Информация об основателе школы (одна запись)';
COMMENT ON COLUMN founder.bio        IS 'Биография основателя в свободной форме';
COMMENT ON COLUMN founder.photo_url  IS 'URL портретной фотографии основателя';
