package dbconverters

import (
	"leetgo/internal/app/store/dbmodel"
	"leetgo/internal/entity"
)

func EntityToDBWord(e entity.Word) dbmodel.Word {
	return dbmodel.Word{
		Data: e.Data,
	}
}

func DBToEntityWord(d dbmodel.Word, dictionary string) entity.Word {
	return entity.Word{
		Data: d.Data,
	}
}
