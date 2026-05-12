-- Достижения школы: победы на турнирах, награды, рекорды.
CREATE TABLE IF NOT EXISTS achievements (
    id          BIGSERIAL    PRIMARY KEY,
    title       TEXT         NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    year        INT          NOT NULL,
    image_url   TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT achievements_year_check CHECK (year >= 1900 AND year <= 2200)
);

COMMENT ON TABLE  achievements             IS 'Достижения школы';
COMMENT ON COLUMN achievements.year        IS 'Год получения достижения';
COMMENT ON COLUMN achievements.image_url   IS 'URL фотографии или диплома';

CREATE INDEX IF NOT EXISTS achievements_year_idx ON achievements (year DESC);
