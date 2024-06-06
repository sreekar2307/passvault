package backup

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"passVault/dtos"
	"passVault/resources"
	"time"
)

var (
	backupFlagSet *flag.FlagSet
	ist           *time.Location
	istErr        error
)

func init() {
	backupFlagSet = flag.NewFlagSet("backup", flag.ExitOnError)
	ist, istErr = time.LoadLocation("Asia/Kolkata")
	if istErr != nil {
		panic(istErr)
	}
}

func Backup(ctx context.Context, args ...string) {
	if err := backupFlagSet.Parse(args); err != nil {
		panic(err)
	}
	var (
		config   = resources.Config()
		today    = time.Now().In(ist).Format("2006-01-02")
		dumpFile = fmt.Sprintf("dump_%s.sql", today)
		cmd      = exec.Command("pg_dump", "-U",
			config.GetString(dtos.ConfigKeys.Database.Username),
			"-h", config.GetString(dtos.ConfigKeys.Database.Host),
			"-p", config.GetString(dtos.ConfigKeys.Database.Port),
			"-d", config.GetString(dtos.ConfigKeys.Database.Name),
			"-f", dumpFile)
	)
	cmd.Env = []string{
		fmt.Sprintf("PGPASSWORD=%s", config.GetString(dtos.ConfigKeys.Database.Password)),
	}
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
