package args

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigFileNoFileFoundError(t *testing.T) {
	cli := CLI{}
	cli.ConfigFile = "dummy_file.json"

	_, err := ParseFileConfigArgs(&cli)

	assert.EqualError(t, err, "open dummy_file.json: no such file or directory")
}

func TestParseConfigFileMarshalError(t *testing.T) {
	cli := CLI{}
	tempFile, _ := os.CreateTemp("", "")
	cli.ConfigFile = tempFile.Name()
	defer os.Remove(cli.ConfigFile)

	_, err := ParseFileConfigArgs(&cli)

	assert.EqualError(t, err, "Error during Unmarshal(): unexpected end of JSON input")
}
