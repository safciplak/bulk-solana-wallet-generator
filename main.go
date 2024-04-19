package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/blocto/solana-go-sdk/types"

)

func main() {
	// Read the number of wallets to generate from user input
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("How many wallets do you want to generate? ")
	scanner.Scan()
	numWalletsInput := scanner.Text()
	numWallets, err := strconv.Atoi(numWalletsInput)
	if err != nil {
		fmt.Println("Error parsing number:", err)
		return
	}

	// Number of goroutines to use
	numWorkers := 16
	walletsPerWorker := (numWallets + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	walletsChannel := make(chan []types.Account, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		start := i * walletsPerWorker
		end := start + walletsPerWorker
		if end > numWallets {
			end = numWallets
		}

		go generateWallets(start, end, walletsChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(walletsChannel)
	}()

	var allWallets []types.Account
	for wallets := range walletsChannel {
		allWallets = append(allWallets, wallets...)
	}

	// Save generated wallets to a file
	file, err := os.Create("generated.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(allWallets, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	file.Write(jsonData)
	fmt.Println("Saved generated wallets to generated.txt")
}

func generateWallets(start, end int, walletsChannel chan<- []types.Account, wg *sync.WaitGroup) {
	defer wg.Done()

	var wallets []types.Account
	for i := start; i < end; i++ {
		account := types.NewAccount()
		wallets = append(wallets, account)
		// fmt.Printf("Generated wallet %d with public key: %s\n", i+1, account.PublicKey.ToBase58())
	}

	walletsChannel <- wallets
}
