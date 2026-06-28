package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	log, err := Setup("text", "info", &buf)
	require.NoError(t, err)

	log.Info("hello", slog.String("key", "value"))

	out := buf.String()
	assert.Contains(t, out, "hello")
	assert.Contains(t, out, "key=value")
}

func TestSetup_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	log, err := Setup("json", "info", &buf)
	require.NoError(t, err)

	log.Info("hello", slog.String("key", "value"))

	out := buf.String()
	var entry map[string]any
	require.NoError(t, json.Unmarshal([]byte(strings.TrimSpace(out)), &entry))
	assert.Equal(t, "hello", entry["msg"])
	assert.Equal(t, "value", entry["key"])
}

func TestSetup_DebugLevel(t *testing.T) {
	var buf bytes.Buffer
	log, err := Setup("text", "debug", &buf)
	require.NoError(t, err)

	log.Debug("debug-msg")
	assert.Contains(t, buf.String(), "debug-msg")
}

func TestSetup_InfoLevelFiltersDebug(t *testing.T) {
	var buf bytes.Buffer
	log, err := Setup("text", "info", &buf)
	require.NoError(t, err)

	log.Debug("debug-msg")
	assert.NotContains(t, buf.String(), "debug-msg")
}

func TestSetup_InvalidFormat(t *testing.T) {
	_, err := Setup("xml", "info", nil)
	assert.Error(t, err)
}

func TestSetup_InvalidLevel(t *testing.T) {
	_, err := Setup("text", "trace", nil)
	assert.Error(t, err)
}
