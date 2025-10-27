package main

import (
	conf "github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/cmd/server"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/store"
)

func main() {
	// Base logger
	appLog := logg.NewLogg("appLog")

	// Config
	config := conf.NewArgsENV(appLog)

	// Extra loggers
	businessLog := logg.NewLogg(config.GetBusinessLog())
	infraLog := logg.NewLogg(config.GetInfraLog())

	// DataBase
	storeDB := store.NewStoreDB(config, infraLog)

	defer businessLog.Close()
	defer infraLog.Close()
	defer storeDB.Close()
	defer appLog.Close()

	server.Start(config, appLog, infraLog, businessLog, storeDB)
}
