CREATE SCHEMA IF NOT EXISTS words;

CREATE TABLE IF NOT EXISTS words.russian (
    data VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_words_data ON words.russian (data text_pattern_ops);