CREATE SCHEMA IF NOT EXISTS words;

CREATE TABLE IF NOT EXISTS words.definitions (
    word VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    definition TEXT NOT NULL,
    dictionary VARCHAR(50) DEFAULT 'russian' NOT NULL
);

INSERT INTO words.definitions (word, definition, dictionary)
SELECT word, definition, dictionary
FROM definitions.definitions
WHERE EXISTS (
    SELECT 1 FROM information_schema.tables 
    WHERE table_schema = 'definitions' AND table_name = 'definitions'
);

CREATE INDEX IF NOT EXISTS idx_definitions_word ON words.definitions (word);
CREATE INDEX IF NOT EXISTS idx_definitions_dict_word ON words.definitions (dictionary, word);

DROP TABLE IF EXISTS definitions.definitions;