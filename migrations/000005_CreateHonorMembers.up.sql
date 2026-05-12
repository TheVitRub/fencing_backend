-- Доска почёта: участники, отличившиеся в истории клуба.
-- sort_order управляет порядком отображения на странице.
CREATE TABLE IF NOT EXISTS honor_members (
    id          BIGSERIAL    PRIMARY KEY,
    name        TEXT         NOT NULL,
    title       TEXT         NOT NULL DEFAULT '',
    description TEXT         NOT NULL DEFAULT '',
    photo_url   TEXT         NOT NULL DEFAULT '',
    sort_order  INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE  honor_members             IS 'Доска почёта клуба';
COMMENT ON COLUMN honor_members.name        IS 'Имя и фамилия участника';
COMMENT ON COLUMN honor_members.title       IS 'Звание или титул, например "Мастер клинка"';
COMMENT ON COLUMN honor_members.sort_order  IS 'Порядок сортировки (меньше = выше)';

CREATE INDEX IF NOT EXISTS honor_members_sort_idx ON honor_members (sort_order ASC, id ASC);
