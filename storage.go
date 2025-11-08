package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Storage struct {
	FileName string
}

func NewStorage(fileName string) *Storage {
	return &Storage{
		FileName: fileName,
	}
}

func (storage *Storage) Save(data Todos) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal Todos data: %w", err)
	}

	if err := os.WriteFile(storage.FileName, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (storage *Storage) Load(todos *Todos) error {
	jsonData, err := os.ReadFile(storage.FileName)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	if err := json.Unmarshal(jsonData, todos); err != nil {
		return fmt.Errorf("failed to unmarshal Todos data: %w", err)
	}
	
	return nil
}
