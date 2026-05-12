-- Начальные данные: страницы-заглушки и администратор по умолчанию.

-- Страницы с базовым JSON-скелетом.
-- Фронт разбирает content и рендерит нужные секции.
INSERT INTO pages (slug, title, content) VALUES
  ('home',
   'Главная',
   '{"hero":{"title":"Школа фехтования Ferrum et Gloria","subtitle":"XVI–XVII век · Исторический бой на клинках"},"intro":"Мы изучаем боевые искусства европейских мастеров эпохи Возрождения."}'),
  ('events',
   'События',
   '{"subtitle":"Турниры, показательные выступления, открытые тренировки"}'),
  ('plans',
   'Учебные планы',
   '{"subtitle":"Программа обучения на каждый период"}'),
  ('honor',
   'Доска почёта',
   '{"subtitle":"Те, кто прославил наш клуб"}'),
  ('achievements',
   'Достижения школы',
   '{"subtitle":"Победы, которыми мы гордимся"}'),
  ('founder',
   'Об основателе',
   '{}')
ON CONFLICT (slug) DO NOTHING;

-- Администратор по умолчанию: логин=admin, пароль=admin123
-- Хэш получен командой: htpasswd -bnBC 10 "" admin123 | tr -d ':\n'
-- ОБЯЗАТЕЛЬНО смените пароль через панель управления после первого входа.
INSERT INTO admins (login, password_hash) VALUES
  ('admin', '$2a$10$e04ncTORCzcZ3.uZkOPEheFkB4MctZWgkxV3JbZtOVz9lwCetad8q')
ON CONFLICT (login) DO NOTHING;
