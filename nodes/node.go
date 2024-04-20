package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yusufwib/blockchain-medical-record/models/dblockchain"
)

type Blockchain struct {
	Chain []dblockchain.Block `json:"chain"`
}

var blockchain Blockchain

const dataDirectory = "./data/"
const nodesFile = "./nodes/nodes.json"

func main() {
	nodeID := flag.String("NODE_ID", "", "Node ID")
	flag.Parse()

	if *nodeID == "" {
		fmt.Println("Please specify a node ID using --NODE_ID flag")
		os.Exit(1)
	}
	// Run the blockchain node
	runNode(*nodeID)
}

func runNode(nodeID string) {
	// Register with discovery service
	err := registerNodeWithDiscovery(nodeID)
	if err != nil {
		log.Fatalf("Error registering with discovery service: %v", err)
	}
	log.Printf("Node %s is running...\n", nodeID)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, err := mineBlock(nodeID)
		if err != nil {
			log.Printf("Error mining block: %v", err)
			continue
		}

	}
}

func mineBlock(nodeID string) (response []dblockchain.Block, err error) {
	resp, err := http.Get("http://localhost:9009/v1/blockchain/mine-all")
	if err != nil {
		return response, fmt.Errorf("failed to fetch data: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response body: %v", err)
	}

	resp.Body.Close()

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("failed to store data to DB: %v", err)
	}

	blocks, err := loadNodeData(nodeID)
	if err != nil {
		return response, fmt.Errorf("failed to load node data: %v", err)
	}

	currentBlockMap := make(map[string]dblockchain.Block, 0)
	for _, v := range blocks {
		currentBlockMap[v.Hash] = v
	}

	for _, v := range response {
		// cBlock := currentBlockMap[v.Hash]
		// if cBlock.Hash == "" {
		blockchain.Chain = append(blockchain.Chain, v)

		log.Println("\n\n\nLogging medical record block storage...")
		storeBlockData(v, nodeID)

		log.Println("Synchronizing block data...")
		err = syncBlock()
		if err != nil {
			log.Printf("Error syncing block: %v", err)
		}
		// }
	}

	return
}

func storeBlockData(block dblockchain.Block, nodeID string) {
	// Filename for storing block data
	filename := dataDirectory + fmt.Sprintf("%s.json", nodeID)

	// Read existing block data from file
	var blocks []dblockchain.Block
	data, err := os.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(data, &blocks)
		if err != nil {
			log.Printf("Error unmarshalling block data: %v\n", err)
			return
		}
	} else if !os.IsNotExist(err) {
		log.Printf("Error reading existing block data: %v\n", err)
	}

	// Remove any block with the same ID
	updatedBlocks := make([]dblockchain.Block, 0)
	for _, b := range blocks {
		if b.Hash != block.Hash {
			updatedBlocks = append(updatedBlocks, b)
		}
	}

	// Append the new block
	updatedBlocks = append(updatedBlocks, block)

	// Marshal the updated blocks
	data, err = json.Marshal(updatedBlocks)
	if err != nil {
		log.Printf("Error marshalling block data: %v\n", err)
		return
	}

	// Write updated data back to the file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		log.Printf("Error writing block data to file: %v\n", err)
	}
}

func syncBlock() error {
	otherNodes, err := loadNodes()
	if err != nil {
		return err
	}

	// map[nodeID]map[hash]encryptedData
	blockHashMap := make(map[string]map[string]string, 0)
	for _, nodeID := range otherNodes {
		blocks, err := loadNodeData(nodeID)
		if err != nil {
			continue
		}

		for _, block := range blocks {
			if blockHashMap[block.Hash] == nil {
				blockHashMap[block.Hash] = make(map[string]string)
			}
			blockHashMap[block.Hash][nodeID] = block.EncryptedData
		}
	}

	for hash, blockHash := range blockHashMap {
		isValid := true
		encryptedDataPivot := blockHash[otherNodes[0]]
		pivotSameNode := []string{}
		pivotNotSameNode := []string{}

		for _, node := range otherNodes {
			if encryptedDataPivot != blockHash[node] {
				isValid = false
				pivotNotSameNode = append(pivotNotSameNode, node)
			} else {
				pivotSameNode = append(pivotSameNode, node)
			}
		}

		if !isValid {
			log.Printf("\nDiscrepancy Alert: Encrypted data does not match hash %s", hash)

			gotHacked := pivotNotSameNode
			if len(pivotNotSameNode) > len(pivotSameNode) {
				gotHacked = pivotSameNode
			}

			for _, v := range gotHacked {
				log.Printf("\nSecurity Breach Detected: Node compromised on node %s", v)
			}

		}
	}

	return nil
}

func loadNodeData(nodeID string) (blocks []dblockchain.Block, err error) {
	filename := dataDirectory + fmt.Sprintf("%s.json", nodeID)

	// Read existing block data from file
	data, err := os.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(data, &blocks)
		if err != nil {
			log.Printf("Error unmarshalling block data from load node %s: %v\n", nodeID, err)
			return
		}
	} else if !os.IsNotExist(err) {
		log.Printf("Error reading existing block data: %v\n", err)
	}

	return
}

func loadNodes() ([]string, error) {
	// Load nodes from nodes.json file
	fileData, err := os.ReadFile(nodesFile)
	if os.IsNotExist(err) {
		// If file not found, create it
		err := createNodesFile()
		if err != nil {
			return nil, err
		}
		// Return an empty array since the file was just created
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}

	var nodes []string
	if err := json.Unmarshal(fileData, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func createNodesFile() error {
	// Create an empty nodes.json file
	emptyNodes := []string{}
	data, err := json.Marshal(emptyNodes)
	if err != nil {
		return err
	}
	err = os.WriteFile(nodesFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func saveNodes(nodes []string) error {
	// Save nodes to nodes.json file
	data, err := json.Marshal(nodes)
	if err != nil {
		return err
	}
	err = os.WriteFile(nodesFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func registerNodeWithDiscovery(nodeID string) error {
	// Load existing nodes
	nodes, err := loadNodes()
	if err != nil {
		return err
	}

	for _, n := range nodes {
		if n == nodeID {
			return nil // Node already exists, no need to register again
		}
	}

	// Add new node
	nodes = append(nodes, nodeID)

	// Save updated nodes
	err = saveNodes(nodes)
	if err != nil {
		return err
	}

	return nil
}
