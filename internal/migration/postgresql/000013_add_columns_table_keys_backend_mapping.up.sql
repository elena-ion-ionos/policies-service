ALTER TABLE keys_backends_mapping ADD COLUMN operation VARCHAR(255) NOT NULL;
ALTER TABLE keys_backends_mapping ADD COLUMN created_at timestamp NOT NULL;
ALTER TABLE keys_backends_mapping ADD COLUMN updated_at timestamp NOT NULL;