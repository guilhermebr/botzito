package common

import "go.mongodb.org/mongo-driver/bson"

// Atodoc converts any struct to bson Document
func Atodoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
