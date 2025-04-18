package main

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// Define spinner configurations with both sequence and delay
type spinnerConfig struct {
	sequence []string
	delay    time.Duration
}

// Create a map of named spinner configurations
var spinnerConfigs = map[string]spinnerConfig{
	"dots": {
		sequence: []string{".", "..", "...", "...."},
		delay:    time.Millisecond * 200,
	},
	"clock": {
		sequence: []string{"-", "\\", "|", "/"},
		delay:    time.Millisecond * 100,
	},
	"brackets": {
		sequence: []string{"[    ]", "[=   ]", "[==  ]", "[=== ]", "[====]", "[ ===]", "[  ==]", "[   =]"},
		delay:    time.Millisecond * 80,
	},
	"arrows": {
		sequence: []string{"v", "<", "^", ">"},
		delay:    time.Millisecond * 120,
	},
}

// Function to get a configured spinner
func getSpinner(configName string) *pterm.SpinnerPrinter {
	config, exists := spinnerConfigs[configName]
	if !exists {
		// Default to clock if config name doesn't exist
		config = spinnerConfigs["clock"]
	}

	return pterm.DefaultSpinner.
		WithSequence(config.sequence...).
		WithDelay(config.delay).
		WithRemoveWhenDone(false)
}

func main() {
	// Show usage when no arguments are given
	if len(os.Args) < 2 {
		pterm.Error.Println("Usage: spinner <command> [args...]")
		os.Exit(1)
	}

	// Build the command to run, using variadic slice
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	// Input command as string variable for print statements
	inputcmd := strings.Join(os.Args[1:], " ")

	// Comment stdout/stderr to suppress output during the spinning
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// Only pass stdin so the command can accept input if needed
	cmd.Stdin = os.Stdin

	// Start the spinner
	spinner, _ := getSpinner("clock").WithText("Running " + inputcmd).Start()

	// Start the command
	err := cmd.Start()
	if err != nil {
		spinner.Fail("Failed to start command: " + err.Error())
		os.Exit(1)
	}

	// Non blocking goroutine waiting until command is finished
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	// Spinner loop while the command is running
	for {
		select {
		case err := <-done:
			if err != nil {
				spinner.Fail(inputcmd, " failed: "+err.Error())
				os.Exit(1)
			} else {
				spinner.Success(inputcmd, " completed")
				os.Exit(0)
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
