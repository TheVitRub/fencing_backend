-- Добавление поддержки множественных изображений для постов.
-- Хранится JSON-массив URL: ["/uploads/abc.jpg", "/uploads/def.jpg"].
-- Старые поля image_url / photo_url оставлены для обратной совместимости
-- и используются как «главная обложка», когда нужен один рендер.
ALTER TABLE events        ADD COLUMN IF NOT EXISTS images TEXT NOT NULL DEFAULT '[]';
ALTER TABLE achievements  ADD COLUMN IF NOT EXISTS images TEXT NOT NULL DEFAULT '[]';
ALTER TABLE honor_members ADD COLUMN IF NOT EXISTS images TEXT NOT NULL DEFAULT '[]';

COMMENT ON COLUMN events.images        IS 'JSON-массив URL изображений события';
COMMENT ON COLUMN achievements.images  IS 'JSON-массив URL изображений достижения';
COMMENT ON COLUMN honor_members.images IS 'JSON-массив URL дополнительных фотографий участника';
