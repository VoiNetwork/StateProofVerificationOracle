package main

import (
	"fmt"
	"log"
	"os"
	"state-proof-relayer/internal/querier"
	"state-proof-relayer/internal/servicestate"
	"state-proof-relayer/internal/utilities"
	"state-proof-relayer/internal/writer"
	"strings"
	"time"
)

type ServiceConfiguration struct {
	LogPath      string
	GenesisRound uint64
	StatePath    string
	NodePath     string
	BackoffTime  time.Duration
}

func fetchStateProof(state *servicestate.ServiceState, nodeQuerier querier.Querier, contractWriter writer.Writer) error {
	err := state.Load()
	if err != nil {
		return err
	}

	log.Printf("Fetching state proof")

	nextStateProofData, err := nodeQuerier.QueryNextStateProofData(state)

	if err != nil {
		return err
	}

	log.Printf("FETCHED proof for round %d", nextStateProofData.Message.Lastattestedround)

	err = contractWriter.UploadStateProof(state, nextStateProofData)
	if err != nil {
		return err
	}

	log.Printf("UPLOADED proof for round %d", nextStateProofData.Message.Lastattestedround)
	err = state.Save()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Hello World!")

	var config ServiceConfiguration
	err := utilities.DecodeFromFile(&config, "config.json")
	if err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	logFile, err := os.OpenFile(config.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file log file: %s", err)
	}

	log.SetOutput(logFile)
	state, err := servicestate.InitializeState(config.StatePath, config.GenesisRound)
	if err != nil {
		log.Printf("Could not initialize state: %s", err)
		return
	}

	nodeQuerier, err := querier.InitializeQuerier(config.NodePath)
	if err != nil {
		log.Printf("Could not initialize querier: %s", err)
		return
	}

	log.Printf("Setup complete")

	contractWriter := writer.InitializeWriter()

	log.Printf("Writer initilized - fetching state proofs")

	for {
		err = fetchStateProof(state, *nodeQuerier, *contractWriter)
		if err == nil {
			continue
		}

		if strings.Contains(err.Error(), "HTTP 404") || strings.Contains(err.Error(), "given round is greater than the latest round") {
			log.Printf("Sleeping: ", config.BackoffTime*time.Millisecond)
			log.Print(err.Error())
			time.Sleep(config.BackoffTime * time.Millisecond)
			continue
		}

		log.Printf("Error while fetching state proof: %s", err)
		break
	}
}
