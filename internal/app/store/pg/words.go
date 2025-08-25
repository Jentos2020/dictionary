package pg

import (
	"context"
	"fmt"
	"leetgo/internal/app/store/dbmodel"
	"leetgo/internal/consts"

	"gorm.io/gorm/clause"
)

func (db *PGStore) UpdateWord(ctx context.Context, oldWord dbmodel.Word, newWord dbmodel.Word, dict string) (bool, error) {
	tx := db.Db.Begin()
	defer tx.Rollback()

	if dict == "" {
		dict = consts.DefaultDict
	}
	fullTable := fmt.Sprintf("%s.%s", db.Schema, dict)

	result := tx.Table(fullTable).Where("data = ?", oldWord.Data).Delete(&dbmodel.Word{})
	if result.Error != nil {
		return false, fmt.Errorf("delete old in %s: %w", fullTable, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	dbNew := dbmodel.Word{Data: newWord.Data}
	if err := tx.Table(fullTable).Clauses(clause.OnConflict{DoNothing: true}).Create(&dbNew).Error; err != nil {
		return false, fmt.Errorf("insert new in %s: %w", fullTable, err)
	}

	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("commit update %s: %w", fullTable, err)
	}
	return true, nil
}

func (db *PGStore) RemoveWord(ctx context.Context, word dbmodel.Word, dict string) (bool, error) {
	tx := db.Db.Begin()
	defer tx.Rollback()

	if dict == "" {
		dict = consts.DefaultDict
	}
	fullTable := fmt.Sprintf("%s.%s", db.Schema, dict)

	result := tx.Table(fullTable).Where("data = ?", word.Data).Delete(&dbmodel.Word{})
	if result.Error != nil {
		return false, fmt.Errorf("delete from %s: %w", fullTable, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("commit delete %s: %w", fullTable, err)
	}
	return true, nil
}

func (db *PGStore) AddWord(ctx context.Context, word dbmodel.Word, dict string) error {
	tx := db.Db.Begin()
	defer tx.Rollback()

	if dict == "" {
		dict = consts.DefaultDict
	}
	fullTable := fmt.Sprintf("%s.%s", db.Schema, dict)

	dbWord := dbmodel.Word{Data: word.Data}
	if err := tx.Table(fullTable).Clauses(clause.OnConflict{DoNothing: true}).Create(&dbWord).Error; err != nil {
		return fmt.Errorf("add word to table %s: %w", fullTable, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Transaction commit error %s: %w", fullTable, err)
	}
	return nil
}

func (db *PGStore) GetWords(ctx context.Context, tableName string) (trie dbmodel.Words, err error) {
	var allWords dbmodel.Words
	var tables []string

	if tableName == "" {
		tables, err = db.GetDictionaryTables(ctx, db.Schema)
		if err != nil {
			return dbmodel.Words{}, err
		}
	} else {
		tables = []string{tableName}
		var exists bool
		err := db.Db.WithContext(ctx).Raw(
			`SELECT EXISTS (
				SELECT 1 FROM information_schema.tables 
				WHERE table_schema = ? AND table_name = ?
			)`, db.Schema, tableName,
		).Scan(&exists).Error
		if err != nil {
			return dbmodel.Words{}, err
		}
		if !exists {
			return dbmodel.Words{}, fmt.Errorf("table %s not found", tableName)
		}
	}

	for _, table := range tables {
		fullTableName := fmt.Sprintf("%s.%s", db.Schema, table)
		if err := db.Db.WithContext(ctx).Table(fullTableName).Find(&allWords).Error; err != nil {
			return dbmodel.Words{}, err
		}
	}
	return allWords, nil
}
