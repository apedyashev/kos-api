CREATE TABLE IF NOT EXISTS knowledge_spaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE nodes
    ADD COLUMN IF NOT EXISTS owner_id UUID,
    ADD COLUMN IF NOT EXISTS space_id UUID;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'nodes_space_id_fkey'
          AND conrelid = 'nodes'::regclass
    ) THEN
        ALTER TABLE nodes
            ADD CONSTRAINT nodes_space_id_fkey
            FOREIGN KEY (space_id) REFERENCES knowledge_spaces(id);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS nodes_space_id_updated_at_idx
    ON nodes (space_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS links (
    source_id UUID NOT NULL REFERENCES nodes(id),
    target_id UUID NOT NULL REFERENCES nodes(id),
    PRIMARY KEY (source_id, target_id)
);

CREATE INDEX IF NOT EXISTS links_target_id_idx ON links (target_id);