package logger

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	legal := []string{
		"debug", "DEBUG", "info", "INFO", "warn", "WARN",
		"error", "ERROR", "fatal", "FATAL",
	}
	illegal := []string{
		"illegal string",
	}
	output := "stdout"

	for _, lvl := range legal {
		opt := Options{
			Level:   lvl,
			Outputs: []string{output},
		}
		_, err := New(opt)
		assert.NoError(t, err, "log init fail: level = ", lvl)
	}

	for _, lvl := range illegal {
		opt := Options{
			Level:   lvl,
			Outputs: []string{output},
		}
		_, err := New(opt)
		assert.Error(t, err, "log init fail: level = ", lvl)
	}
}

func TestLogger(t *testing.T) {
	output := "./scaffold.test.log"

	removeLogFile := func() {
		os.Remove(output)
	}

	clean := func() {
		os.Truncate(output, 0)
	}

	setup := func() {
		opt := Options{
			Level:   "DEBUG",
			Outputs: []string{output},
		}

		l, _ := New(opt)
		SetLogger(l)
	}

	assertLogContains := func(expected map[string]interface{}) {
		actual, err := os.ReadFile(output)
		assert.Nil(t, err, "open log file fail：", output)

		t.Log(string(actual))

		var actualJSON = make(map[string]interface{})
		err = json.Unmarshal(actual, &actualJSON)
		assert.Nil(t, err, "JSON decode log fail: ", actual)

		for k, v := range expected {
			av, has := actualJSON[k]
			assert.True(t, has, "key not exist in log：", k)
			assert.Equal(t, v, av, "log content error ：want = ", v, ", got = ", av)
		}
	}

	testLog := func(msg string) {
		defer removeLogFile()

		clean()
		setup()

		Debug(msg)
		assertLogContains(map[string]interface{}{
			"level":   "debug",
			"message": msg,
		})

		clean()
		Info(msg)
		assertLogContains(map[string]interface{}{
			"level":   "info",
			"message": msg,
		})

		clean()
		Warn(msg)
		assertLogContains(map[string]interface{}{
			"level":   "warn",
			"message": msg,
		})

		clean()
		Error(msg)
		assertLogContains(map[string]interface{}{
			"level":   "error",
			"message": msg,
		})

		clean()
		Errorf(msg)
		assertLogContains(map[string]interface{}{
			"level":   "error",
			"message": msg,
		})

		clean()
		assert.Panics(t, func() { Panic(msg) })

		clean()
		With(
			"client_ip", "127.0.0.1",
			"customer_id", "WE6TEST1",
			"host", "test-service.production.svc.cluster.local:9005",
			"method", "POST",
			"is_panic", false,
		).Debug(msg)
		assertLogContains(map[string]interface{}{
			"level":       "debug",
			"message":     msg,
			"client_ip":   "127.0.0.1",
			"customer_id": "WE6TEST1",
			"host":        "test-service.production.svc.cluster.local:9005",
			"method":      "POST",
			"is_panic":    false,
		})

		clean()

		GetLoggerWithKeys(map[string]interface{}{
			"level":   "debug",
			"message": msg,
		}, map[string]interface{}{
			"extra_info": "extra_info",
		}).Debug(msg)

		assertLogContains(map[string]interface{}{
			"level":      "debug",
			"message":    msg,
			"extra_info": "extra_info",
		})

		ext := make(map[string]interface{})
		ext["time"] = time.Now()
		LoadExtra(ext)
	}

	testLog("hello world")
}

func TestGetLogger(t *testing.T) {
	assert.NotNil(t, GetLogger())
}
