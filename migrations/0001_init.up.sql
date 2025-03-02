CREATE TABLE files
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename   TEXT   NOT NULL,
    path       TEXT   NOT NULL,
    size       BIGINT NOT NULL,
    mime_type  TEXT   NOT NULL,
    created_at TIMESTAMP        DEFAULT now(),
    updated_at TIMESTAMP        DEFAULT now()
);

CREATE INDEX idx_files_filename ON files (filename);
