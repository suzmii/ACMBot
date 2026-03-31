package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func WriteArtifact(t testing.TB, paths []string, data []byte) {
	t.Helper()

	outputDir := filepath.Join(paths[:len(paths)-1]...)
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			t.Fatalf("create artifact dir %q: %v", outputDir, err)
		}
	}

	outputPath := filepath.Join(paths...)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		t.Fatalf("write artifact %q: %v", outputPath, err)
	}
}
