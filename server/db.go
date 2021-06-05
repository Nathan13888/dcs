package server

import (
	"dcs/config"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var gormConfig = &gorm.Config{
	// Logger: logger.New(
	// 	getDBLogger(),
	// 	logger.Config{},
	// ),
	DisableAutomaticPing:   false,
	PrepareStmt:            true,
	SkipDefaultTransaction: false,
}

func GetDB() *gorm.DB {
	// db, err := gorm.Open(postgres.New(postgres.Config{
	// 	DSN:                  config.DSN(),
	// 	PreferSimpleProtocol: true,
	// }), gormConfig)
	db, err := gorm.Open(
		sqlite.Open(config.DSN()),
		gormConfig,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to DB")
	}
	return db
}

func InitDB() {
	db := GetDB()
	db.AutoMigrate(DownloadJob{})
}

func DBAddJob(newJob *DownloadJob) {
	res := GetDB().Create(newJob)
	if res.Error != nil {
		logError(res.Error)
	}
}

func DBUpdateJob(job *DownloadJob) {
	GetDB().Save(job)
}

// TODO: check sql injection
func DBGetJob(id string) (DownloadJob, bool) {
	var ret DownloadJob
	res := GetDB().First(ret, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return ret, false
	} else if res.Error != nil {
		logError(res.Error)
		return ret, false
	}
	return ret, true
}

func DBGetJobs() []DownloadJob {
	var ret []DownloadJob
	res := GetDB().Find(&ret)
	if res.Error != nil {
		logError(res.Error)
		return []DownloadJob{}
	}
	return ret
}
