ALTER TABLE keys
ALTER COLUMN description SET DEFAULT '',
ALTER COLUMN canonical_user_id SET DEFAULT NULL;

ALTER TABLE keys ADD COLUMN backends text[] NOT NULL DEFAULT '{}'::text[];
