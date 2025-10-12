package main

import (
	"github.com/boginskiy/Gophermart/cmd/config"
	"github.com/boginskiy/Gophermart/cmd/server"
	"github.com/boginskiy/Gophermart/internal/logg"
	"github.com/boginskiy/Gophermart/internal/store"
)

func main() {
	// Base logger
	appLog := logg.NewLogg("appLog")

	// Args
	args := config.NewArgs(appLog)

	// Extra loggers
	businessLog := logg.NewLogg(args.GetBusinessLog())
	infraLog := logg.NewLogg(args.GetInfraLog())

	// DataBase
	storeDB := store.NewStoreDB(args, infraLog)

	defer businessLog.Clouse()
	defer infraLog.Clouse()
	defer storeDB.Clouse()
	defer appLog.Clouse()

	server.Start(args, appLog, infraLog, businessLog, storeDB)
}
