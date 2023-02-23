package util

import (
	"relorer/test_go_microservice/internal/model"
	"sync"

	"github.com/brianvoe/gofakeit/v6"
)

// A code born out of obscurity
func DeleteRandomFieldsFromAuthorsAndComments(docs []*model.Document) []interface{} {
	result := make([]interface{}, len(docs))
	var wg sync.WaitGroup
	wg.Add(len(docs))

	for i, v := range docs {
		go func(doc *model.Document, index int) {
			defer wg.Done()
			result[index] = ToMaps(doc)

			if authors, ok := result[index].(map[string]interface{})["authors"]; ok {
				for _, author := range authors.([]interface{}) {
					authorMap := author.(map[string]interface{})
					deleteRandomKeys(authorMap)
					if comments, ok := authorMap["comments"]; ok {
						for _, comment := range comments.([]interface{}) {
							deleteRandomKeys(comment.(map[string]interface{}))
						}
					}
				}
			}

		}(v, i)
	}

	wg.Wait()
	return result
}

func deleteRandomKeys(m map[string]interface{}) {
	numToDelete := gofakeit.IntRange(0, len(m))

	for i := 0; i < numToDelete; i++ {
		index := gofakeit.IntRange(0, len(m)-1)
		for k := range m {
			if index == 0 {
				delete(m, k)
				break
			}
			index--
		}
	}
}
