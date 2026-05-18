UPDATE events
SET date = date + INTERVAL '7 hours'
WHERE title = 'Занятие CourtSword'
  AND discipline = 'CourtSword'
  AND date > NOW()
  AND to_char(date AT TIME ZONE 'UTC', 'HH24:MI') = '12:30';
