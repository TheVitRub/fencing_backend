ALTER TABLE glossary_terms
  DROP COLUMN IF EXISTS images,
  DROP COLUMN IF EXISTS image_url;

ALTER TABLE knowledge_articles
  DROP COLUMN IF EXISTS images,
  DROP COLUMN IF EXISTS image_url;

