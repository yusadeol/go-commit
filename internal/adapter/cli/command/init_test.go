package command

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yusadeol/go-commit/internal/adapter/cli/dispatcher"
)

func TestInit(t *testing.T) {
	t.Run("should be able to create the configuration file", func(t *testing.T) {
		tempDir := t.TempDir()
		configurationDirPath := filepath.Join(tempDir, ".config")
		err := os.Mkdir(configurationDirPath, 0755)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		init := NewInit(configurationDirPath)
		result, err := init.Execute(&dispatcher.CommandInput{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "configuration file created successfully"
		if !strings.Contains(result.Message.StripMarkup(), expected) {
			t.Fatalf("expected message to contain %q, got: %q", expected, result.Message.StripMarkup())
		}
	})

	t.Run("should return error when the configuration file already exists", func(t *testing.T) {
		tempDir := t.TempDir()
		configurationDirPath := filepath.Join(tempDir, ".config")
		err := os.Mkdir(configurationDirPath, 0755)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		configurationFilePath := filepath.Join(configurationDirPath, "commit.json")
		file, err := os.Create(configurationFilePath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		err = file.Close()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		init := NewInit(configurationDirPath)
		result, err := init.Execute(&dispatcher.CommandInput{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "configuration file already exists"
		if !strings.Contains(result.Message.StripMarkup(), expected) {
			t.Fatalf("expected message to contain %q, got: %q", expected, result.Message.StripMarkup())
		}
	})
}
