package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func SendToGoogleScript(url string, payload any) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status %s: %s", resp.Status, string(body))
	}

	return body, nil
}

type HistoryResponse struct {
	Success bool          `json:"success"`
	Data    []HistoryItem `json:"data"`
}

type HistoryItem struct {
	Time  string `json:"Время"`
	Sheet string `json:"Таблица"`
	Row   int    `json:"Строка"`
}

type UniqueRecord struct {
	Sheet string
	Row   int
}

func extractUnique(data []byte) ([]UniqueRecord, error) {
	var res HistoryResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	seen := make(map[string]struct{})
	var unique []UniqueRecord

	for _, item := range res.Data {
		key := fmt.Sprintf("%s:%d", item.Sheet, item.Row)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			unique = append(unique, UniqueRecord{
				Sheet: item.Sheet,
				Row:   item.Row,
			})
		}
	}

	return unique, nil
}

func HandleChanges(resp []byte, logger *slog.Logger) error {
	uniqueItems, err := extractUnique(resp)
	if err != nil {
		return fmt.Errorf("failed to extract unique items %w", err)
	}

	for _, item := range uniqueItems {
		logger.Info("Unique row", "sheet", item.Sheet, "row", item.Row)
	}
	return nil
}
