package server

type Response struct {
	Message string         `json:"message"`
	Result  map[string]any `json:"result"`
}

type TablePayload struct {
	Name string `json:"name"`
}

type EntryPayload struct {
	TableName string   `json:"table"`
	Key       []string `json:"key"`
	Value     any      `json:"value"`
}

type QueryPayload struct {
	Key       []string `json:"key"`
	TableName string   `json:"table"`
}

type IData interface {
	Add(payload []byte) error
	Del(payload []byte) error
	Query(payload []byte) error
}
