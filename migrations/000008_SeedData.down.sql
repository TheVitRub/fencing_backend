-- Удаляем только seed-данные, не трогая таблицы.
DELETE FROM pages  WHERE slug IN ('home','events','plans','honor','achievements','founder');
DELETE FROM admins WHERE login = 'admin';
