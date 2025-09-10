package store

import (
	"context"
	"leetgo/config"
	"leetgo/internal/app/store/dbmodel"
)

type (
	DefinitionReadWriter interface {
		GetDefinition(ctx context.Context, word, dict string) (string, bool, error)
		AddDefinition(ctx context.Context, def dbmodel.Definition, dict string) error
		UpdateDefinition(ctx context.Context, word, newDef, dict string) (bool, error)
		RemoveDefinition(ctx context.Context, word, dict string) (bool, error)
	}

	Repository interface {
		WordReadWriter
		DictsReadWriter
		DefinitionReadWriter
		Ping
	}

	WordReadWriter interface {
		UpdateWord(ctx context.Context, oldWord dbmodel.Word, newWord dbmodel.Word, dict string) (bool, error)
		RemoveWord(ctx context.Context, word dbmodel.Word, dict string) (bool, error)
		AddWord(ctx context.Context, words dbmodel.Word, dict string) error
	}

	DictsReadWriter interface {
		GetDictionaryTables(ctx context.Context, schema string) ([]string, error)
		WriteDictsToDb(ctx context.Context, cfg config.Config, fileName string) error
		GetWords(ctx context.Context, tableName string) (trie dbmodel.Words, err error)
	}

	Ping interface {
		Ping(ctx context.Context) error
	}
)
