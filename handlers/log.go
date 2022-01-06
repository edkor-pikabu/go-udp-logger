package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type handler struct {
	db *sql.DB
}

type Message struct {
	Name string
	Group string
	Data string
}


func New(database *sql.DB) *handler {
	return &handler{
		db: database,
	}
}

func (h *handler) Handle(messages []string) {
	if len(messages) == 0 {
		return
	}

	valueStrings := make([]string, 0, len(messages))
	valueArgs := make([]interface{}, 0, len(messages) * 3)

	var message Message
	for _, value := range messages {
		err := json.Unmarshal([]byte(value), &message)
		if err != nil {
			fmt.Println("Parsing error", err)
			continue
		}

		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, message.Name)
		valueArgs = append(valueArgs, message.Group)
		valueArgs = append(valueArgs, message.Data)
		fmt.Println(message)
	}

	ctx := context.Background()
	SQL := fmt.Sprintf(
		"INSERT INTO logs (record_name, group_name, data) VALUES %s",
		strings.Join(valueStrings, ","),
	)
	_, err := h.db.ExecContext(ctx, SQL, valueArgs...)
	if err != nil {
		fmt.Println(err)
	}
}
