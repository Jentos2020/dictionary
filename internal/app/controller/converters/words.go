package converters

import (
	"leetgo/internal/app/store/dbmodel"
	"leetgo/internal/entity"
	"leetgo/internal/gen"
)

func DBWordsToEntity(dbWords dbmodel.Words) entity.Words {
	words := make(entity.Words, len(dbWords))
	for i, w := range dbWords {
		words[i] = entity.Word{Data: w.Data}
	}
	return words
}

func EntityWordsToDB(entWords entity.Words) dbmodel.Words {
	words := make(dbmodel.Words, len(entWords))
	for i, w := range entWords {
		words[i] = dbmodel.Word{Data: w.Data}
	}
	return words
}

func GenToEntityWord(g *gen.Word) entity.Word {
	if g == nil {
		return entity.Word{}
	}
	var data, dict string
	if g.Data != nil {
		data = *g.Data
	}
	if g.Dictionary != nil {
		dict = *g.Dictionary
	}
	return entity.Word{
		Data:       data,
		Dictionary: dict,
	}
}

func EntityToGenWord(e entity.Word) gen.Word {
	return gen.Word{
		Data:       &e.Data,
		Dictionary: &e.Dictionary,
	}
}
