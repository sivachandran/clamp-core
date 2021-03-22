package services

import (
	"clamp-core/config"
	"clamp-core/repository"
	"fmt"
	"os"
	"testing"
)

var mockDB repository.MockDB

func TestMain(m *testing.M) {
	err := config.Load()
	if err != nil {
		fmt.Printf("Loading config failed: %s\n", err)
	}

	repository.SetDB(&mockDB)

	err = InitServiceRequestWorkers()
	if err != nil {
		fmt.Printf("Initializinng service request workers failed: %s", err)
	}

	err = InitResumeWorkers()
	if err != nil {
		fmt.Printf("Initializinng resume workers failed: %s", err)
	}

	os.Exit(m.Run())
}
