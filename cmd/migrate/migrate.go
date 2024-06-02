package migrate

import (
	"context"
	"flag"
	"passVault/models"
	"passVault/resources"
)

var migrateFlagSet *flag.FlagSet

func init() {
	migrateFlagSet = flag.NewFlagSet("migrate", flag.ExitOnError)
}

func Migrate(ctx context.Context, args ...string) {
	if err := migrateFlagSet.Parse(args); err != nil {
		panic(err)
	}

	var (
		db = resources.Database()
	)

	if err := db.WithContext(ctx).AutoMigrate(
		&models.User{},
		&models.UserSalt{},
		&models.Password{},
		&models.PasswordGenerationHistory{},
		&models.PasswordVersion{},
	); err != nil {
		panic(err.Error())
	}
}
