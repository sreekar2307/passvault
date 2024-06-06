package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/resources"
	"sort"
	"time"
)

type BackUpServiceImpl struct {
	s3 resources.S3
}

func NewBackupService(s3 resources.S3) interfaces.BackupService {
	return &BackUpServiceImpl{s3: s3}
}

var (
	ist, _ = time.LoadLocation("Asia/Kolkata")
)

func (b BackUpServiceImpl) BackupDb(ctx context.Context) error {

	var (
		config       = resources.Config()
		today        = time.Now().In(ist).Format("2006-01-02")
		dumpFileName = fmt.Sprintf("dump_%s.tar", today)
		dumpFilePath = fmt.Sprintf("/tmp/%s", dumpFileName)
		cmd          = exec.Command("pg_dump", "-U",
			config.GetString(dtos.ConfigKeys.Database.Username),
			"-h", config.GetString(dtos.ConfigKeys.Database.Host),
			"-p", config.GetString(dtos.ConfigKeys.Database.Port),
			"-d", config.GetString(dtos.ConfigKeys.Database.Name),
			"-F", "tar",
			"-f", dumpFilePath)
		s3Path = fmt.Sprintf("%s/%s", config.GetString(dtos.ConfigKeys.DatabaseBackup.Location), dumpFileName)
	)
	cmd.Env = []string{
		fmt.Sprintf("PGPASSWORD=%s", config.GetString(dtos.ConfigKeys.Database.Password)),
	}
	if err := cmd.Run(); err != nil {
		return err
	}

	file, err := os.Open(dumpFilePath)
	if err != nil {
		return err
	}

	defer func() {
		file.Close()
		os.Remove(dumpFilePath)
	}()

	if err := b.s3.Push(ctx, config.GetString(dtos.ConfigKeys.DatabaseBackup.Bucket), s3Path, file); err != nil {
		return err
	}

	objects, err := b.s3.List(ctx, config.GetString(dtos.ConfigKeys.DatabaseBackup.Bucket),
		config.GetString(dtos.ConfigKeys.DatabaseBackup.Location))
	if err != nil {
		return err
	}

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].LastModified.Before(objects[j].LastModified)
	})

	if len(objects) > 5 {
		paths := make([]string, len(objects)-5)
		for i := 0; i < len(objects)-5; i++ {
			paths[i] = objects[i].Key
		}
		if err := b.s3.DeleteBulk(ctx, config.GetString(dtos.ConfigKeys.DatabaseBackup.Bucket), paths); err != nil {
			return err
		}
	}
	return nil
}
