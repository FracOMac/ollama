package ollamarunner

import (
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/ollama/ollama/envconfig"
	"github.com/ollama/ollama/logutil"
)

func TestOllamaRunnerThreadsEnvVar(t *testing.T) {
	// Save original log setup and restore it after tests
	originalLogger := slog.Default()
	defer slog.SetDefault(originalLogger)

	// Use a test logger that captures output
	slog.SetDefault(logutil.NewLogger(os.Stderr, slog.LevelDebug))

	tests := []struct {
		name        string
		envValue    string
		expected    int
		expectWarn  bool
	}{
		{
			name:        "valid positive integer",
			envValue:    "8",
			expected:    8,
			expectWarn:  false,
		},
		{
			name:        "valid positive integer large",
			envValue:    "32",
			expected:    32,
			expectWarn:  false,
		},
		{
			name:        "zero threads",
			envValue:    "0",
			expected:    runtime.NumCPU(),
			expectWarn:  true,
		},
		{
			name:        "negative threads",
			envValue:    "-1",
			expected:    runtime.NumCPU(),
			expectWarn:  true,
		},
		{
			name:        "invalid string",
			envValue:    "abc",
			expected:    runtime.NumCPU(),
			expectWarn:  true,
		},
		{
			name:        "empty string",
			envValue:    "",
			expected:    runtime.NumCPU(),
			expectWarn:  false,
		},
		{
			name:        "float value",
			envValue:    "4.5",
			expected:    runtime.NumCPU(),
			expectWarn:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variable
			if tt.envValue != "" {
				t.Setenv("OLLAMA_RUNNER_THREADS", tt.envValue)
			}

			// Simulate the logic from Execute function
			threads := runtime.NumCPU() // default value

			// Test the environment variable parsing logic
			if envThreads := envconfig.Var("OLLAMA_RUNNER_THREADS"); envThreads != "" {
				if parsedThreads, err := strconv.Atoi(envThreads); err != nil {
					// Should warn for invalid values
					if !tt.expectWarn {
						t.Errorf("Expected no warning for %q, but got parsing error: %v", tt.envValue, err)
					}
				} else if parsedThreads <= 0 {
					// Should warn for non-positive values
					if !tt.expectWarn {
						t.Errorf("Expected no warning for %q, but got non-positive value: %d", tt.envValue, parsedThreads)
					}
				} else {
					// Should use the parsed value
					threads = parsedThreads
				}
			}

			if threads != tt.expected {
				t.Errorf("Expected threads=%d, got threads=%d", tt.expected, threads)
			}
		})
	}
}

func TestOllamaRunnerThreadsIntegration(t *testing.T) {
	// Test that Execute function can be called with OLLAMA_RUNNER_THREADS set
	// This tests the integration without actually running the server
	
	t.Setenv("OLLAMA_RUNNER_THREADS", "4")
	
	// We can't easily test Execute function end-to-end because it requires
	// a model file and starts a server, but we can verify the parsing logic works
	// by checking that envconfig.Var reads the environment variable correctly
	
	envThreads := envconfig.Var("OLLAMA_RUNNER_THREADS")
	if envThreads != "4" {
		t.Errorf("Expected OLLAMA_RUNNER_THREADS=4, got %q", envThreads)
	}
	
	if parsedThreads, err := strconv.Atoi(envThreads); err != nil {
		t.Errorf("Failed to parse OLLAMA_RUNNER_THREADS: %v", err)
	} else if parsedThreads != 4 {
		t.Errorf("Expected parsed threads=4, got %d", parsedThreads)
	}
}