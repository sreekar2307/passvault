package backup

import (
	"context"
	"flag"
	"passVault/dependency"
)

var (
	backupFlagSet *flag.FlagSet
)

func init() {
	backupFlagSet = flag.NewFlagSet("backup", flag.ExitOnError)
}

func Backup(ctx context.Context, args ...string) {
	if err := backupFlagSet.Parse(args); err != nil {
		panic(err)
	}
	dependencies := dependency.Dependencies()
	if err := dependencies.BackupService.BackupDb(ctx); err != nil {
		panic(err)
	}
}
