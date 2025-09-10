CREATE SCHEMA IF NOT EXISTS words;

CREATE TABLE IF NOT EXISTS words.definitions (
    word VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    definition TEXT NOT NULL,
    dictionary VARCHAR(50) DEFAULT 'russian' NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_definitions_word ON words.definitions (word);
CREATE INDEX IF NOT EXISTS idx_definitions_dict_word ON words.definitions (dictionary, word);