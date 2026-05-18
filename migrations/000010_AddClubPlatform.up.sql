-- Users, roles, comments, attendance, learning content and progress.

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    login         VARCHAR(100) NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL DEFAULT '',
    password_hash TEXT NOT NULL DEFAULT '',
    display_name  TEXT NOT NULL DEFAULT '',
    role          VARCHAR(32) NOT NULL DEFAULT 'registered',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_role_check CHECK (role IN ('registered', 'student', 'instructor', 'admin', 'founder'))
);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique_not_empty ON users (email) WHERE email <> '';

CREATE TABLE IF NOT EXISTS user_identities (
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider          VARCHAR(64) NOT NULL,
    provider_user_id  TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_identities_provider_unique UNIQUE (provider, provider_user_id)
);

INSERT INTO users (login, email, password_hash, display_name, role)
SELECT login, login || '@local.fencing', password_hash, login, 'admin'
FROM admins
ON CONFLICT (login) DO NOTHING;

CREATE TABLE IF NOT EXISTS comments (
    id          BIGSERIAL PRIMARY KEY,
    target_type VARCHAR(32) NOT NULL,
    target_id   BIGINT NOT NULL,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body        TEXT NOT NULL,
    status      VARCHAR(16) NOT NULL DEFAULT 'visible',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT comments_target_type_check CHECK (target_type IN ('event', 'achievement', 'plan')),
    CONSTRAINT comments_status_check CHECK (status IN ('visible', 'hidden', 'deleted'))
);
CREATE INDEX IF NOT EXISTS comments_target_idx ON comments (target_type, target_id, created_at);

ALTER TABLE events ADD COLUMN IF NOT EXISTS type VARCHAR(32) NOT NULL DEFAULT 'event';
ALTER TABLE events ADD COLUMN IF NOT EXISTS status VARCHAR(32) NOT NULL DEFAULT 'scheduled';

CREATE TABLE IF NOT EXISTS event_attendees (
    event_id   BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status     VARCHAR(16) NOT NULL DEFAULT 'going',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (event_id, user_id),
    CONSTRAINT event_attendees_status_check CHECK (status IN ('going', 'cancelled'))
);

CREATE TABLE IF NOT EXISTS notifications (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id   BIGINT REFERENCES events(id) ON DELETE CASCADE,
    title      TEXT NOT NULL,
    body       TEXT NOT NULL DEFAULT '',
    is_read    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS notifications_user_idx ON notifications (user_id, is_read, created_at DESC);

CREATE TABLE IF NOT EXISTS instructor_profiles (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    name           TEXT NOT NULL DEFAULT '',
    photo_url      TEXT NOT NULL DEFAULT '',
    specialization TEXT NOT NULL DEFAULT '',
    weapons        TEXT NOT NULL DEFAULT '',
    experience     TEXT NOT NULL DEFAULT '',
    quote          TEXT NOT NULL DEFAULT '',
    bio            TEXT NOT NULL DEFAULT '',
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS knowledge_articles (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    category   TEXT NOT NULL DEFAULT '',
    body       TEXT NOT NULL DEFAULT '',
    visibility VARCHAR(32) NOT NULL DEFAULT 'public',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT knowledge_visibility_check CHECK (visibility IN ('public', 'registered', 'student'))
);

CREATE TABLE IF NOT EXISTS glossary_terms (
    id          BIGSERIAL PRIMARY KEY,
    term        TEXT NOT NULL,
    category    TEXT NOT NULL DEFAULT '',
    definition  TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS student_progress (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    discipline    TEXT NOT NULL,
    level         TEXT NOT NULL DEFAULT '',
    passed_checks TEXT NOT NULL DEFAULT '[]',
    instructor_note TEXT NOT NULL DEFAULT '',
    updated_by    BIGINT REFERENCES users(id) ON DELETE SET NULL,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT student_progress_unique UNIQUE (user_id, discipline)
);
