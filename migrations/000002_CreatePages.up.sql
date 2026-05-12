-- Страницы сайта с произвольным JSON-контентом.
-- Slug — уникальный идентификатор страницы (например, "home", "events").
-- Content — JSON-объект, структуру которого определяет фронтенд.
--   Это позволяет менять наполнение страниц без изменений схемы БД.
CREATE TABLE IF NOT EXISTS pages (
    id         BIGSERIAL    PRIMARY KEY,
    slug       VARCHAR(100) NOT NULL,
    title      TEXT         NOT NULL,
    content    TEXT         NOT NULL DEFAULT '{}',
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT pages_slug_unique UNIQUE (slug)
);

COMMENT ON TABLE  pages            IS 'Страницы сайта с настраиваемым контентом';
COMMENT ON COLUMN pages.slug       IS 'Уникальный URL-идентификатор страницы';
COMMENT ON COLUMN pages.title      IS 'Заголовок страницы (отображается в <title> и в h1)';
COMMENT ON COLUMN pages.content    IS 'JSON-объект с произвольным контентом, разбирается фронтендом';
COMMENT ON COLUMN pages.updated_at IS 'Время последнего изменения страницы';
