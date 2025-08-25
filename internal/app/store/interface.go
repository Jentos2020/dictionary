package store

import (
	"context"
	"leetgo/config"
	"leetgo/internal/app/store/dbmodel"
)

type (
	Repository interface {
		WordReadWriter
		DictsReadWriter
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
