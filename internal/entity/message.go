package entity

type Message struct {
	Key   int    `db:"id" json:"id"`
	Value string `db:"value" json:"value"`
}
