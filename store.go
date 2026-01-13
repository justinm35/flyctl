package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func FetchStoredData[T any](key string) T {
	home, _ := os.UserHomeDir()
	storeDir := filepath.Join(home, ".config/flyctl/store")

	file, err := os.ReadFile(fmt.Sprintf("%s/%s.json", storeDir, key))
	if err != nil {
		log.Printf("Store Error: %s \n", string(err.Error()))
	}

	var out T
	json.Unmarshal(file, &out)

	return out
}

func StoreData[T any](key string, value T) error {
	home, _ := os.UserHomeDir()
	storeDir := filepath.Join(home, "/.config/flyctl/store")
	storeDirContents, _ := os.ReadDir(storeDir)

	if len(storeDirContents) == 0 {
		log.Printf("Store Dir did not exist, creating a new one")
		createStoreErr := os.MkdirAll(storeDir, 0o755)
		if createStoreErr != nil {
			log.Printf("Store Error: %s \n", string(createStoreErr.Error()))
		}
	}

	marhsalledJson, err := json.Marshal(value)
	if err != nil {
		log.Printf("Store Error: %s \n", string(err.Error()))
	}

	writeErr := os.WriteFile(fmt.Sprintf("%s/%s.json", storeDir, key), marhsalledJson, 0o755)
	if writeErr != nil {
		log.Printf("Store Error: %s \n", string(writeErr.Error()))
	}

	return nil
}
