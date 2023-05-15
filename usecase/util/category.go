package util

import "github.com/jseow5177/pockteer-be/data/entity"

func GetCatIDs(categories []*entity.Category) []string {
	ids := make([]string, len(categories))
	for idx, cat := range categories {
		ids[idx] = cat.GetCatID()
	}
	return ids
}

func GetCatIDToCategoryMap(categories []*entity.Category) map[string]*entity.Category {
	_map := make(map[string]*entity.Category)

	for _, cat := range categories {
		_map[cat.GetCatID()] = cat
	}
	return _map
}
