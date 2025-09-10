CREATE SCHEMA IF NOT EXISTS definitions;

CREATE TABLE IF NOT EXISTS definitions.definitions (
    word VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    definition TEXT NOT NULL,
    dictionary VARCHAR(50) DEFAULT 'russian' NOT NULL
);


INSERT INTO definitions.definitions (word, definition, dictionary)
SELECT word, definition, dictionary
FROM words.definitions
WHERE EXISTS (
    SELECT 1 FROM information_schema.tables 
    WHERE table_schema = 'words' AND table_name = 'definitions'
);


CREATE INDEX IF NOT EXISTS idx_definitions_word ON definitions.definitions (word);
CREATE INDEX IF NOT EXISTS idx_definitions_dict_word ON definitions.definitions (dictionary, word);

DROP TABLE IF EXISTS words.definitions;
