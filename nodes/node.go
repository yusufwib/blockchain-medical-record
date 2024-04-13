package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yusufwib/blockchain-medical-record/models/dblockchain"
)

// tbd rewrite
var dataMining []dblockchain.Block

func fetchDataAndStoreToDB() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var response []dblockchain.Block
		resp, err := http.Get("http://localhost:9009/v1/blockchain/mine-all")
		if err != nil {
			fmt.Println("Failed to fetch data:", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body:", err)
			continue
		}

		resp.Body.Close()

		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Failed to store data to Badger DB:", err)
		}

		// checking if the data is the same like another nodes.. if ok then save
		dataMining = response
	}
}

func main() {
	fetchDataAndStoreToDB()
}
