package projects

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RegistryFile is the name of the projects registry file
const RegistryFile = "projects.json"

// ProjectRegistry represents the mapping of project names to paths
type ProjectRegistry map[string]string

// LoadRegistry loads the projects registry from the global steria directory
func LoadRegistry() (ProjectRegistry, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".steria", RegistryFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return make(ProjectRegistry), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var registry ProjectRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, err
	}

	return registry, nil
}

// SaveRegistry saves the projects registry to the global steria directory
func SaveRegistry(registry ProjectRegistry) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, ".steria")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, RegistryFile), data, 0644)
}

// AddProject adds a project to the registry
func AddProject(name, path string) error {
	registry, err := LoadRegistry()
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	registry[name] = absPath
	return SaveRegistry(registry)
}

// RemoveProject removes a project from the registry
func RemoveProject(name string) error {
	registry, err := LoadRegistry()
	if err != nil {
		return err
	}

	delete(registry, name)
	return SaveRegistry(registry)
}
