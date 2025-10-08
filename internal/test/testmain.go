package test

import (
	"go-graphql-aggregator/internal/logger"
	"os"
	"testing"
)

// SetupTests inicializa o logger para evitar panic em todos os pacotes de teste
func SetupTests(m *testing.M) {
	os.Setenv("LOG_MODE", "silent")
	logger.Init()
	os.Exit(m.Run())
}
