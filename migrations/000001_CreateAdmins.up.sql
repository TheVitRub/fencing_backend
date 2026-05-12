-- Таблица администраторов сайта.
-- Хранит логин и bcrypt-хэш пароля.
-- Количество администраторов не ограничено, но на практике обычно один.
CREATE TABLE IF NOT EXISTS admins (
    id            BIGSERIAL    PRIMARY KEY,
    login         VARCHAR(100) NOT NULL,
    password_hash TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT admins_login_unique UNIQUE (login)
);

COMMENT ON TABLE  admins              IS 'Администраторы сайта клуба';
COMMENT ON COLUMN admins.login        IS 'Уникальный логин для входа в панель управления';
COMMENT ON COLUMN admins.password_hash IS 'bcrypt-хэш пароля (cost=10)';
