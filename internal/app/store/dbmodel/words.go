package dbmodel

type Word struct {
	Data string `gorm:"primaryKey;type:varchar(255);unique"`
}

type Words = []Word
