package env

import (
	"testing"

	"github.com/Excoriate/go-terradagger/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestMergeEnvVars_EmptyVars(t *testing.T) {
	result := MergeEnvVars()

	assert.Empty(t, result, "Result should be empty when no environment vars are provided")
}

func TestMergeEnvVars_SingleVarMap(t *testing.T) {
	envVars := Vars{"KEY1": "\"value1\""}
	result := MergeEnvVars(envVars)

	assert.Equal(t, 1, len(result), "Result should contain one environment variable")
	assert.Equal(t, "value1", result["KEY1"], "Value should match and be stripped of double quotes")
}

func TestMergeEnvVars_MultipleVarMaps(t *testing.T) {
	envVars1 := Vars{"KEY1": "\"value1\""}
	envVars2 := Vars{"KEY2": "value2", "KEY3": ""}

	result := MergeEnvVars(envVars1, envVars2)

	assert.Equal(t, 2, len(result), "Result should contain two environment variables")
	assert.Equal(t, "value1", utils.RemoveDoubleQuotes(result["KEY1"]), "Value of KEY1 should match and be stripped of double quotes")
	assert.Equal(t, "value2", utils.RemoveDoubleQuotes(result["KEY2"]), "Value of KEY2 should match and be stripped of double quotes")
	assert.NotContains(t, result, "KEY3", "Keys with empty values should not be included")
}

func TestMergeEnvVars_OverrideVarMap(t *testing.T) {
	envVars1 := Vars{"KEY1": "\"value1\""}
	envVars2 := Vars{"KEY1": "newValue1"}

	result := MergeEnvVars(envVars1, envVars2)

	assert.Equal(t, 1, len(result), "Result should contain one environment variable")
	assert.Equal(t, "newValue1", utils.RemoveDoubleQuotes(result["KEY1"]), "Value should be overridden by the last value provided")
}
