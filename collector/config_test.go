package collector

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromReader_GenBearerTokenFailed_ReturnFalse(t *testing.T) {
    var buffer bytes.Buffer
    buffer.WriteString("{}")
		var c *Config = &Config{}
    ok := c.LoadConfigFromReader(&buffer)
  	assert.False(t, ok, *c)
}