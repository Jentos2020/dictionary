package pg

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"leetgo/config"
	"leetgo/internal/app/store/dbmodel"
	"leetgo/internal/consts"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm/clause"
)

func (db *PGStore) GetDictionaryTables(ctx context.Context, schema string) ([]string, error) {
	var tables []string
	err := db.Db.WithContext(ctx).Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = ? 
		AND table_type = 'BASE TABLE' 
		AND table_name != ?`, schema, consts.MigrationsTable).Scan(&tables).Error
	return tables, err
}

func (db *PGStore) WriteDictsToDb(ctx context.Context, cfg config.Config, fileName string) error {
	var files []os.DirEntry

	if fileName != "" {
		target := filepath.Join(cfg.Dicts, fileName+".txt")
		if _, err := os.Stat(target); err != nil {
			return fmt.Errorf("file %s not found: %w", target, err)
		}

		files = []os.DirEntry{fileEntry{target, fileName + ".txt"}}
	} else {
		allFiles, err := os.ReadDir(cfg.Dicts)
		if err != nil {
			return err
		}
		files = allFiles
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			tableName := strings.TrimSuffix(file.Name(), ".txt")

			var exists bool
			err := db.Db.WithContext(ctx).Raw(
				`SELECT EXISTS (
					SELECT 1 FROM information_schema.tables 
					WHERE table_schema = ? AND table_name = ?
				)`, db.Schema, tableName,
			).Scan(&exists).Error
			if err != nil {
				return err
			}

			fullTableName := fmt.Sprintf("%s.%s", db.Schema, tableName)

			if !exists {
				createQuery := fmt.Sprintf(`
					CREATE TABLE IF NOT EXISTS %s (
						id SERIAL PRIMARY KEY,
						data VARCHAR(255) UNIQUE NOT NULL
					);
					CREATE INDEX IF NOT EXISTS idx_%s_data ON %s (data text_pattern_ops);
				`, fullTableName, tableName, fullTableName)
				if err := db.Db.WithContext(ctx).Exec(createQuery).Error; err != nil {
					return err
				}
				db.log.Info(fmt.Sprintf("Created table %s", fullTableName))
			} else {
				db.log.Info(fmt.Sprintf("Table %s already exists, skipping creation", fullTableName))
			}

			f, err := os.Open(filepath.Join(cfg.Dicts, file.Name()))
			if err != nil {
				return err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			batch := []dbmodel.Word{}

			for scanner.Scan() {
				word := strings.TrimSpace(scanner.Text())
				if word != "" {
					batch = append(batch, dbmodel.Word{Data: word})
				}

				if len(batch) >= cfg.BatchSize {
					if err := db.Db.WithContext(ctx).Table(fullTableName).
						Clauses(clause.OnConflict{DoNothing: true}).Create(&batch).Error; err != nil {
						return err
					}
					batch = []dbmodel.Word{}
				}
			}

			if len(batch) > 0 {
				if err := db.Db.WithContext(ctx).Table(fullTableName).
					Clauses(clause.OnConflict{DoNothing: true}).Create(&batch).Error; err != nil {
					return err
				}
			}

			db.log.Info(fmt.Sprintf("Seeded table %s", fullTableName))
		}
	}

	return nil
}

func (db *PGStore) SeedDefinitions(ctx context.Context, cfg config.Config, fileName string) error {
	var files []os.DirEntry
	defaultFileName := "definitions.txt"

	if fileName != "" {
		target := filepath.Join(cfg.Dicts, fileName)
		if _, err := os.Stat(target); err != nil {
			return fmt.Errorf("definition file %s not found: %w", target, err)
		}
		files = []os.DirEntry{fileEntry{target, fileName}}
	} else {
		defaultPath := filepath.Join(cfg.Dicts, defaultFileName)
		if _, err := os.Stat(defaultPath); err == nil {
			files = []os.DirEntry{fileEntry{defaultPath, defaultFileName}}
		} else {
			allFiles, err := os.ReadDir(cfg.Dicts)
			if err != nil {
				return fmt.Errorf("read dir %s: %w", cfg.Dicts, err)
			}
			for _, f := range allFiles {
				if !f.IsDir() && strings.HasPrefix(f.Name(), "def_") && strings.HasSuffix(f.Name(), ".txt") {
					files = append(files, f)
				}
			}
		}
	}

	if len(files) == 0 {
		db.log.Info("No definition files found to seed")
		return nil
	}

	for _, file := range files {
		f, err := os.Open(filepath.Join(cfg.Dicts, file.Name()))
		if err != nil {
			return fmt.Errorf("open file %s: %w", file.Name(), err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		batch := []dbmodel.Definition{}
		fullTableName := "definitions.definitions"

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, ":", 3)
			if len(parts) < 2 {
				db.log.Warn(fmt.Sprintf("Invalid format in %s: %s, skipping", file.Name(), line))
				continue
			}

			word := strings.TrimSpace(parts[0])
			definition := strings.TrimSpace(parts[1])
			dictionary := consts.DefaultDict
			if len(parts) == 3 {
				dictionary = strings.TrimSpace(parts[2])
			}

			if word == "" || definition == "" {
				db.log.Warn(fmt.Sprintf("Empty word or definition in %s: %s, skipping", file.Name(), line))
				continue
			}

			batch = append(batch, dbmodel.Definition{
				Word:       word,
				Definition: definition,
				Dictionary: dictionary,
			})

			if len(batch) >= cfg.BatchSize {
				if err := db.Db.WithContext(ctx).Table(fullTableName).
					Clauses(clause.OnConflict{DoNothing: true}).Create(&batch).Error; err != nil {
					return fmt.Errorf("seed definitions to %s: %w", fullTableName, err)
				}
				db.log.Info(fmt.Sprintf("Seeded %d definitions to %s", len(batch), fullTableName))
				batch = []dbmodel.Definition{}
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scan file %s: %w", file.Name(), err)
		}

		if len(batch) > 0 {
			if err := db.Db.WithContext(ctx).Table(fullTableName).
				Clauses(clause.OnConflict{DoNothing: true}).Create(&batch).Error; err != nil {
				return fmt.Errorf("seed definitions to %s: %w", fullTableName, err)
			}
			db.log.Info(fmt.Sprintf("Seeded %d definitions to %s", len(batch), fullTableName))
		}
	}

	return nil
}

type fileEntry struct {
	path string
	name string
}

func (f fileEntry) Name() string      { return f.name }
func (f fileEntry) IsDir() bool       { return false }
func (f fileEntry) Type() fs.FileMode { return 0 }
func (f fileEntry) Info() (fs.FileInfo, error) {
	return os.Stat(f.path)
}
