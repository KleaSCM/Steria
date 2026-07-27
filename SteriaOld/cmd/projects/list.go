package projects

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var jsonOutput bool

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all projects",
		Long:  "List all projects registered in Steria",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList()
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runList() error {
	registry, err := LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if jsonOutput {
		data, err := json.MarshalIndent(registry, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal registry: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s Registered Projects:\n", cyan("📂"))

	if len(registry) == 0 {
		fmt.Println("  No projects found.")
		return nil
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		path := registry[name]
		fmt.Printf("  %s %s\n", green(name), yellow("("+path+")"))
	}

	return nil
}
