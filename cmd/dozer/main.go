package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/spoonboy-io/dozer/internal/hook"
	"github.com/spoonboy-io/reprise"

	"github.com/spoonboy-io/dozer/internal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/spoonboy-io/dozer/internal/morpheus"
	"github.com/spoonboy-io/dozer/internal/state"
	"github.com/spoonboy-io/koan"
)

const (
	DB_CONFIG   = "mysql.env"
	HOOK_CONFIG = "webhook.yaml"
)

var (
	version   = "Development build"
	goversion = "Unknown"
)

var logger *koan.Logger
var st *state.State

func init() {
	st = &state.State{}
	logger = &koan.Logger{}

	// read in the db config
	err := godotenv.Load(DB_CONFIG)
	if err != nil {
		logger.FatalError("Failed to read database config file", err)
	}

	// check for state, if exists load it
	if st.HasSavedState() {
		logger.Info("Loading saved state")
		if err = st.ReadAndParse(); err != nil {
			logger.FatalError("Failed to read or parse saved application state", err)
		}
	} else {
		// first run, and we probably don't want to iterate through previous processes
		// so we capture latest process, start the poll from there. On subsequent application
		// launches we will use the state file, if it is gone, we will be back here, and some
		// webhooks will not have been fired

		// TODO
	}

	logger.Info("Loading webhook configuration file")
	err = hook.ReadAndParseConfig(HOOK_CONFIG)
	if err != nil {
		logger.FatalError("Failed to read webhook configuration file", err)
	}
}

// Shutdown runs on SIGINT and panic, we save the database poll state
// which will be loaded upon application restart
func Shutdown(db *sql.DB) {
	fmt.Println("") // break after ^C
	logger.Warn("Application terminated. Closing database connection")
	db.Close()
	logger.Info("Saving application state")
	if err := st.CreateAndWrite(); err != nil {
		logger.Error("Failed to save application state", err)
	}
}

func main() {
	// write a console banner
	reprise.WriteSimple(&reprise.Banner{
		Name:         "Dozer",
		Description:  "Morpheus Processes with Webhooks",
		Version:      version,
		GoVersion:    goversion,
		WebsiteURL:   "https://spoonboy.io",
		VcsURL:       "https://github.com/spoonboy-io/dozer",
		VcsName:      "Github",
		EmailAddress: "hello@spoonboy.io",
	})

	// connect to database
	var db *sql.DB
	var err error
	cnString := fmt.Sprintf("%s:%s@tcp(%s:3306)/morpheus?parseTime=true", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_SERVER"))
	db, err = sql.Open("mysql", cnString)
	if err != nil {
		logger.FatalError("Failed to create database connection", err)
	}

	defer Shutdown(db)

	if err = db.Ping(); err != nil {
		logger.FatalError("Failed to connect to database", err)
	}
	logger.Info("Connected to database")

	logger.Info("Loading process types from database")
	processTypes := map[string]string{}
	if err := morpheus.GetProcessTypes(db, processTypes); err != nil {
		logger.FatalError("Failed to load process types", err)
	}

	go func() {
		pollSecs := internal.POLL_INTERVAL
		if os.Getenv("POLL_INTERVAL_SECONDS") != "" {
			if pollSecs, err = strconv.Atoi(os.Getenv("POLL_INTERVAL_SECONDS")); err != nil {
				logger.Warn("Could not use POLL_INTERVAL_SECONDS, continuing with default")
			}
			logger.Info("Using POLL_INTERVAL_SECCONDS environment variable")
		}
		pollInterval := time.NewTicker(time.Duration(pollSecs) * time.Second)
		for _ = range pollInterval.C {
			/* temp monitor */
			fmt.Printf("lastProcessId: %d\n", st.LastPollProcessId)
			fmt.Printf("ExecutingProcesses: %v\n", st.ExecutingProcesses)

			if err = morpheus.CheckExecuting(db, st); err != nil {
				logger.Error("Error handling executing processes", err)
			}

			if err := morpheus.GetProcesses(db, st); err != nil {
				logger.Error("Database poll error", err)
			}

			lastPollMsg := fmt.Sprintf("Last datasbase poll performed at %s", st.LastPollTimestamp)
			logger.Info(lastPollMsg)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
