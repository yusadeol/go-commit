package command

import (
	"strings"
	"testing"

	"github.com/yusadeol/go-commit/internal/adapter/cli/dispatcher"

	"github.com/yusadeol/go-commit/internal/domain/vo"
	"github.com/yusadeol/go-commit/internal/infra/service/ai"
)

const mockDiff = `
	diff --git a/example.go b/example.go
	index 6711592..6108c39 100644
	--- a/example.go
	+++ b/example.go
	@@ -2,9 +2,9 @@
	
	package main
	
	-import "fmt"
	+import "fmt"
	
	-func helloWorld() string {
	-    return "Hello world"
	+func helloYsoCode() string {
	+    return "Hello YSO Code"
	}
	
	func main() {
	-    fmt.Println(helloWorld())
	+    fmt.Println(helloYsoCode())
	}
`

type MockProvider struct{}

func (m *MockProvider) Ask(input *ai.ProviderInput) (*ai.ProviderOutput, error) {
	return &ai.ProviderOutput{
		Status: "success",
		Text:   "feat: rename function and update greeting message",
	}, nil
}

type MockDefaultProviderFactory struct{}

func (m *MockDefaultProviderFactory) Create(id string, apiKey string) (ai.Provider, error) {
	return &MockProvider{}, nil
}

func TestGenerate(t *testing.T) {
	t.Run("should be able to generate a commit", func(t *testing.T) {
		mockConfiguration := vo.Configuration{
			AIProviders: map[string]vo.AIProvider{
				"mock": {
					ID:           "mock",
					APIKey:       "fake-api-key",
					DefaultModel: "mock-model",
				},
			},
			Languages: map[string]vo.Language{
				"en_US": {ID: "en_US", DisplayName: "English (US)"},
			},
		}
		generate := NewGenerate(&mockConfiguration, &MockDefaultProviderFactory{})
		result, err := generate.Execute(&dispatcher.CommandInput{
			Arguments: map[string]dispatcher.ArgumentInput{
				"diff": {Value: mockDiff, Meta: dispatcher.Argument{Name: "diff", Description: "Git diff", Required: false}},
			},
			Options: map[string]dispatcher.OptionInput{
				"provider": {Value: "mock", Meta: dispatcher.Option{Name: "provider", Flag: "p", Description: "AI Provider", Default: "mock"}},
				"language": {Value: "en_US", Meta: dispatcher.Option{Name: "language", Flag: "l", Description: "Language", Default: "en_US"}},
				"commit":   {Value: "false", Meta: dispatcher.Option{Name: "commit", Flag: "c", Description: "Commit", Default: "true"}},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ExitCode != vo.ExitCodeSuccess {
			t.Fatalf("unexpected exit code: %v", result.ExitCode)
		}
		expected := "feat: rename function and update greeting message"
		if !strings.Contains(result.Message.StripMarkup(), expected) {
			t.Fatalf("expected message to contain %q, got: %q", expected, result.Message.StripMarkup())
		}
	})
}
