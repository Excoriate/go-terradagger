package config

import (
	"os"
	"reflect"
	"testing"
)

// Setup and Teardown helpers for environment variables
func setEnvVars(vars map[string]string) func() {
	originalVars := make(map[string]string)
	for k, v := range vars {
		originalVars[k] = os.Getenv(k)
		_ = os.Setenv(k, v)
	}
	return func() {
		for k, v := range originalVars {
			_ = os.Setenv(k, v)
		}
	}
}

func TestOptions_GetWorkspace(t *testing.T) {
	opts := New("customWorkspace", nil, nil, nil)
	if ws := opts.GetWorkspace(); ws != "customWorkspace" {
		t.Errorf("Expected customWorkspace, got %s", ws)
	}

	opts = New("", nil, nil, nil)
	if ws := opts.GetWorkspace(); ws != defaultWorkspace {
		t.Errorf("Expected default workspace, got %s", ws)
	}
}

func TestOptions_GetExcludedDirs(t *testing.T) {
	customExcludes := []string{"customExclude/**"}
	opts := New("", nil, customExcludes, nil)
	expected := append(excludedDirsDefault, customExcludes...)

	if dirs := opts.GetExcludedDirs(); !reflect.DeepEqual(dirs, expected) {
		t.Errorf("Expected %v, got %v", expected, dirs)
	}
}

func TestOptions_GetTerraformEnvVars(t *testing.T) {
	cleanup := setEnvVars(map[string]string{"TF_VAR_example": "value"})
	defer cleanup()

	opts := New("", map[string]string{"customVar": "customValue"}, nil, nil)
	expected := map[string]string{"TF_VAR_example": "value", "customVar": "customValue"}

	if vars := opts.GetTerraformEnvVars(); !reflect.DeepEqual(vars, expected) {
		t.Errorf("Expected %v, got %v", expected, vars)
	}
}

func TestOptions_GetAWSEnvVars(t *testing.T) {
	cleanup := setEnvVars(map[string]string{"AWS_ACCESS_KEY_ID": "access", "AWS_SECRET_ACCESS_KEY": "secret"})
	defer cleanup()

	opts := New("", nil, nil, nil)
	expected := map[string]string{"AWS_ACCESS_KEY_ID": "access", "AWS_SECRET_ACCESS_KEY": "secret"}

	if vars := opts.GetAWSEnvVars(); !reflect.DeepEqual(vars, expected) {
		t.Errorf("Expected %v, got %v", expected, vars)
	}
}

// Additional tests should be implemented for other methods following similar patterns.
