INSERT INTO events (title, description, date, location, type, status, discipline, image_url, images, created_at)
SELECT
  'Занятие CourtSword',
  'Плановое занятие по CourtSword: стойка, мера, линия атаки, парирование и ответ. Можно записаться прямо из календаря.',
  date_trunc('day', NOW()) + INTERVAL '3 days' + TIME '12:30',
  'Зал школы, уточнение в группе VK',
  'training',
  'scheduled',
  'CourtSword',
  '',
  '[]',
  NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM events
  WHERE title = 'Занятие CourtSword'
    AND date > NOW()
);

INSERT INTO events (title, description, date, location, type, status, discipline, image_url, images, created_at)
SELECT
  'Открытая тренировка CourtSword',
  'Открытое занятие для зарегистрированных гостей и учеников. Хороший формат, чтобы попробовать дисциплину и задать вопросы инструкторам.',
  date_trunc('day', NOW()) + INTERVAL '10 days' + TIME '12:00',
  'Зал школы, подробности в VK',
  'open',
  'scheduled',
  'CourtSword',
  '',
  '[]',
  NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM events
  WHERE title = 'Открытая тренировка CourtSword'
    AND date > NOW()
);
