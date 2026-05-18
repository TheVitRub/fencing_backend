# API And Data Context

## Public API

Все публичные маршруты находятся под `/api`.

- `POST /api/auth/login` - вход администратора.
- `GET /api/pages/:slug` - получить страницу.
- `GET /api/events` - список событий.
- `GET /api/plans` - список учебных планов.
- `GET /api/honor` - доска почёта.
- `GET /api/achievements` - достижения.
- `GET /api/founder` - информация об основателе.

## Admin API

Маршруты ниже требуют JWT.

- `GET /api/admin/me` - текущий пользователь.
- `GET /api/admin/users` - список пользователей для админа/основателя.
- `PUT /api/admin/users/:id/role` - назначить роль.
- `POST /api/admin/comments` - оставить комментарий авторизованным пользователем.
- `PUT /api/admin/comments/:id/status` - скрыть/вернуть/удалить комментарий.
- `PUT /api/admin/events/:id/attendance` - отметить "Я приду" или отменить участие.
- `GET /api/admin/notifications` - уведомления пользователя.
- `PUT /api/admin/notifications/:id/read` - отметить уведомление прочитанным.
- `POST /api/admin/upload` - загрузить файл.
- `PUT /api/admin/pages/:slug` - создать или обновить страницу.
- `POST /api/admin/events` - создать событие.
- `PUT /api/admin/events/:id` - обновить событие.
- `DELETE /api/admin/events/:id` - удалить событие.
- `POST /api/admin/plans` - создать учебный план.
- `PUT /api/admin/plans/:id` - обновить учебный план.
- `DELETE /api/admin/plans/:id` - удалить учебный план.
- `POST /api/admin/honor` - создать участника доски почёта.
- `PUT /api/admin/honor/:id` - обновить участника.
- `DELETE /api/admin/honor/:id` - удалить участника.
- `POST /api/admin/achievements` - создать достижение.
- `PUT /api/admin/achievements/:id` - обновить достижение.
- `DELETE /api/admin/achievements/:id` - удалить достижение.
- `PUT /api/admin/founder` - обновить информацию об основателе.
- `PUT /api/admin/instructors` - создать/обновить профиль инструктора.
- `POST|PUT|DELETE /api/admin/knowledge` - материалы раздела "Ученикам".
- `POST|PUT|DELETE /api/admin/glossary` - глоссарий.
- `GET|PUT /api/admin/progress` - учебный прогресс.

## Миграции

Миграции лежат в `migrations`.

Текущие группы:

- admins;
- pages;
- events;
- plans;
- honor members;
- achievements;
- founder;
- seed data;
- image arrays.
- club platform: users, roles, comments, attendance, notifications, instructor profiles, knowledge, glossary, progress.

Seed-данные находятся в `000008_SeedData.up.sql`.

## Данные страниц

Таблица `pages` хранит:

- `slug`;
- `title`;
- `content`.

`content` - JSON-строка, которую фронтенд разбирает сам. Это удобно для простых блоков вроде subtitle, hero и intro.

## События и изображения

События поддерживают одиночное изображение и массив изображений. Фронтенд берёт первое изображение как cover, а остальные показывает в галерее.

## Комментарии, роли и посещаемость

Комментарии хранятся в таблице `comments` и привязываются к `event`, `achievement` или `plan`.

Пользовательские роли:

- `registered`
- `student`
- `instructor`
- `admin`
- `founder`

Отметки "Я приду" лежат в `event_attendees`. При изменении даты, места или статуса события сервис создаёт внутренние уведомления для записавшихся пользователей.
