package args

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigFileNoFileFoundError(t *testing.T) {
	filePath := "dummy_file.json"

	_, err := ParseFileConfigArgs(filePath, &CLI)

	assert.EqualError(t, err, "open dummy_file.json: no such file or directory")
}

func TestParseConfigFileMarshalError(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "")
	defer os.Remove(tempFile.Name())

	_, err := ParseFileConfigArgs(tempFile.Name(), &CLI)

	assert.EqualError(t, err, "Error during Unmarshal(): unexpected end of JSON input")
}
