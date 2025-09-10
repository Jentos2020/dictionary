package pg

import (
	"context"
	"errors"
	"fmt"
	"leetgo/internal/app/store/dbmodel"
	"leetgo/internal/consts"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *PGStore) GetDefinition(ctx context.Context, word, dict string) (string, bool, error) {
	if dict == "" {
		dict = consts.DefaultDict
	}
	fullTable := "definitions.definitions"

	var def dbmodel.Definition
	err := db.Db.WithContext(ctx).Table(fullTable).
		Where("word = ? AND dictionary = ?", word, dict).
		First(&def).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("get definition from %s: %w", fullTable, err)
	}
	return def.Definition, true, nil
}

func (db *PGStore) AddDefinition(ctx context.Context, def dbmodel.Definition, dict string) error {
	if dict == "" {
		dict = consts.DefaultDict
	}
	def.Dictionary = dict
	tx := db.Db.Begin()
	defer tx.Rollback()

	fullTable := "definitions.definitions"
	dbDef := dbmodel.Definition{Word: def.Word, Definition: def.Definition, Dictionary: dict}
	if err := tx.Table(fullTable).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&dbDef).Error; err != nil {
		return fmt.Errorf("add definition to %s: %w", fullTable, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit add definition %s: %w", fullTable, err)
	}
	return nil
}

func (db *PGStore) UpdateDefinition(ctx context.Context, word, newDef, dict string) (bool, error) {
	if dict == "" {
		dict = consts.DefaultDict
	}
	tx := db.Db.Begin()
	defer tx.Rollback()

	fullTable := "definitions.definitions"
	result := tx.Table(fullTable).
		Where("word = ? AND dictionary = ?", word, dict).
		Update("definition", newDef)
	if result.Error != nil {
		return false, fmt.Errorf("update definition in %s: %w", fullTable, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("commit update definition %s: %w", fullTable, err)
	}
	return true, nil
}

func (db *PGStore) RemoveDefinition(ctx context.Context, word, dict string) (bool, error) {
	if dict == "" {
		dict = consts.DefaultDict
	}
	tx := db.Db.Begin()
	defer tx.Rollback()

	fullTable := "definitions.definitions"
	result := tx.Table(fullTable).
		Where("word = ? AND dictionary = ?", word, dict).
		Delete(&dbmodel.Definition{})
	if result.Error != nil {
		return false, fmt.Errorf("delete definition from %s: %w", fullTable, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if err := tx.Commit().Error; err != nil {
		return false, fmt.Errorf("commit delete definition %s: %w", fullTable, err)
	}
	return true, nil
}
