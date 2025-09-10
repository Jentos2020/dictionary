package dbmodel

type Definition struct {
	Word       string `gorm:"primaryKey;type:varchar(255);uniqueIndex:idx_dict_word;not null"`
	Definition string `gorm:"type:text;not null"`
	Dictionary string `gorm:"type:varchar(50);default:russian;uniqueIndex:idx_dict_word;not null"`
}
