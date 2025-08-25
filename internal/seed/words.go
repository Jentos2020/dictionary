package seed

// import (
// 	"bufio"
// 	"context"
// 	"fmt"
// 	"io/fs"
// 	"log/slog"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"leetgo/config"
// 	"leetgo/internal/app/store/dbmodel"
// 	"leetgo/internal/app/store/pg"

// 	"gorm.io/gorm/clause"
// )

// func DictionaryToDB(ctx context.Context, db *pg.PGStore, log *slog.Logger, cfg config.Config, fileName string) error {
// 	var files []os.DirEntry

// 	if fileName != "" {
// 		target := filepath.Join(cfg.Dicts, fileName+".txt")
// 		if _, err := os.Stat(target); err != nil {
// 			return fmt.Errorf("file %s not found: %w", target, err)
// 		}

// 		files = []os.DirEntry{fileEntry{target, fileName + ".txt"}}
// 	} else {
// 		allFiles, err := os.ReadDir(cfg.Dicts)
// 		if err != nil {
// 			return err
// 		}
// 		files = allFiles
// 	}

// 	for _, file := range files {
// 		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
// 			tableName := strings.TrimSuffix(file.Name(), ".txt")

// 			var exists bool
// 			err := db.Db.WithContext(ctx).Raw(
// 				`SELECT EXISTS (
// 					SELECT 1 FROM information_schema.tables
// 					WHERE table_schema = ? AND table_name = ?
// 				)`, db.Schema, tableName,
// 			).Scan(&exists).Error
// 			if err != nil {
// 				return err
// 			}

// 			fullTableName := fmt.Sprintf("%s.%s", db.Schema, tableName)

// 			if !exists {
// 				createQuery := fmt.Sprintf(`
// 					CREATE TABLE IF NOT EXISTS %s (
// 						id SERIAL PRIMARY KEY,
// 						data VARCHAR(255) UNIQUE NOT NULL
// 					);
// 					CREATE INDEX IF NOT EXISTS idx_%s_data ON %s (data text_pattern_ops);
// 				`, fullTableName, tableName, fullTableName)
// 				if err := db.Db.WithContext(ctx).Exec(createQuery).Error; err != nil {
// 					return err
// 				}
// 				log.Info(fmt.Sprintf("Created table %s", fullTableName))
// 			} else {
// 				log.Info(fmt.Sprintf("Table %s already exists, skipping creation", fullTableName))
// 			}

// 			f, err := os.Open(filepath.Join(cfg.Dicts, file.Name()))
// 			if err != nil {
// 				return err
// 			}
// 			defer f.Close()

// 			scanner := bufio.NewScanner(f)
// 			batch := []dbmodel.Word{}

// 			for scanner.Scan() {
// 				word := strings.TrimSpace(scanner.Text())
// 				if word != "" {
// 					batch = append(batch, dbmodel.Word{Data: word})
// 				}

// 				if len(batch) >= cfg.BatchSize {
// 					if err := db.Db.WithContext(ctx).Table(fullTableName).
// 						Clauses(clause.OnConflict{DoNothing: true}).Create(&batch).Error; err != nil {
// 						return err
// 					}
// 					batch = []dbmodel.Word{}
// 				}
// 			}

// 			if len(batch) > 0 {
// 				if err := db.Db.WithContext(ctx).Table(fullTableName).
// 					Clauses(clause.OnConflict{DoNothing: true}).Create(&batch).Error; err != nil {
// 					return err
// 				}
// 			}

// 			log.Info(fmt.Sprintf("Seeded table %s", fullTableName))
// 		}
// 	}

// 	return nil
// }

// type fileEntry struct {
// 	path string
// 	name string
// }

// func (f fileEntry) Name() string      { return f.name }
// func (f fileEntry) IsDir() bool       { return false }
// func (f fileEntry) Type() fs.FileMode { return 0 }
// func (f fileEntry) Info() (fs.FileInfo, error) {
// 	return os.Stat(f.path)
// }
