package app

import (
	"fmt"
	"task_service/pkg/database"

	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitDatabaseHook(app *Application) error {
	db, err := database.OpenMysqlDatabase(&app.GetConfig().Database)
	if err != nil {
		return fmt.Errorf("InitDatabaseHook: %s", err)
	}

	app.SetDatabase(db)

	if err := migration(app); err != nil {
		return fmt.Errorf("InitDatabaseHook: %s", err)
	}
	return nil
}

func migration(app *Application) error {
	driver, err := mysql.WithInstance(app.GetDatabase(), &mysql.Config{})
	if err != nil {
		return fmt.Errorf("Migrate: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", app.GetConfig().MigrationFilePath),
		app.GetConfig().Database.DBName,
		driver)
	if err != nil {
		return fmt.Errorf("Migrate: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Migrate: %v", err)
	}

	app.GetLogger().Info("Migrate: Migrate successfully")
	return nil
}
