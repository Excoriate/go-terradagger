package config

import (
	"regexp"
	"testing"
)

// Test that version strings follow semantic versioning.
func TestVersionPatterns(t *testing.T) {
	semanticVersionPattern := `^\d+\.\d+\.\d+$`
	versionTests := []struct {
		name    string
		version string
	}{
		{"TerraformDefaultVersion", TerraformDefaultVersion},
		{"TerragruntDefaultVersion", TerragruntDefaultVersion},
	}

	for _, tt := range versionTests {
		matched, err := regexp.MatchString(semanticVersionPattern, tt.version)
		if err != nil {
			t.Fatalf("Regex failed: %v", err)
		}
		if !matched {
			t.Errorf("%s = %v, want pattern %v", tt.name, tt.version, semanticVersionPattern)
		}
	}
}

// Test that Docker image names are in a valid format.
func TestDockerImageNames(t *testing.T) {
	imageTests := []struct {
		name  string
		image string
	}{
		{"TerraformDefaultImage", TerraformDefaultImage},
		{"TerragruntDefaultImage", TerragruntDefaultImage},
	}

	for _, tt := range imageTests {
		if !isValidDockerImageName(tt.image) {
			t.Errorf("%s = %v, want a valid Docker image name", tt.name, tt.image)
		}
	}
}

// Placeholder for a real validation function for Docker image names.
func isValidDockerImageName(image string) bool {
	// For the sake of this example, a simple check for '/' will suffice.
	// A real implementation should more thoroughly validate against Docker's image naming conventions.
	return regexp.MustCompile(`/`).MatchString(image)
}
