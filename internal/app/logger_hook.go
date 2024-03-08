package app

import (
	"task_service/pkg/logger"
)

func initLoggerApplicationHook(app *Application) error {
	l, err := logger.New(logger.Options{
		Level:      app.GetConfig().LogLevel,
		Outputs:    app.GetConfig().LogFile,
		ErrOutputs: app.GetConfig().ErrorLogFile,
	})

	if err != nil {
		panic(err)
	}

	logger.SetLogger(l)
	app.SetLogger(l.Sugar())

	return nil
}
