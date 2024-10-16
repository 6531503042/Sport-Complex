package models

type (
	KafkaOffset struct {
		Offset int64 `json:"offset" bson:"offset"`
	}
)