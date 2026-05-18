DELETE FROM glossary_terms
WHERE term IN ('CourtSword', 'Мера', 'Темп', 'Линия', 'Парирование', 'Ответ', 'Выпад', 'Свободная игра');

DELETE FROM knowledge_articles
WHERE title IN ('Первая тренировка', 'Безопасность в зале', 'CourtSword: что это за занятие', 'Экипировка ученика');

DELETE FROM instructor_profiles
WHERE name IN ('Основатель школы', 'Инструктор CourtSword');

ALTER TABLE events DROP COLUMN IF EXISTS discipline;
