ALTER TABLE knowledge_spaces
    ADD COLUMN IF NOT EXISTS slug TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS knowledge_spaces_slug_key
    ON knowledge_spaces (slug);
