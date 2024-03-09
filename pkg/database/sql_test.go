package database

import (
	"task_service/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConnectionString(t *testing.T) {
	tests := []struct {
		opt      config.DatabaseOption
		isErr    bool
		expected string
	}{
		{
			config.DatabaseOption{
				Driver:   "mysql",
				Host:     "127.0.0.1",
				Port:     8088,
				Username: "user",
				Password: "pass",
				DBName:   "test",
			},
			false,
			"user:pass@tcp(127.0.0.1:8088)/test?loc=Local&multiStatements=true&parseTime=true",
		},
		{
			config.DatabaseOption{
				Driver:   "mysql",
				Host:     "127.0.0.1",
				Port:     8088,
				Username: "user",
				Password: "pass",
				DBName:   "test",
				Timezone: "test",
			},
			true,
			"",
		},
	}

	for _, testItem := range tests {
		connectString, err := GetConnectionString(&testItem.opt)
		if testItem.isErr {
			assert.NotNil(t, err)
			continue
		}
		assert.Nil(t, err)
		assert.Equal(t, testItem.expected, connectString)
	}

}
