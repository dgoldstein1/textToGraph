package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"strconv"
)

// checks environment for required env vars
var logFatalf = log.Fatalf
var logErrorf = log.Errorf
var logMsg = log.Infof

func parseEnv(fileLocation string) {
	if fileLocation == "" {
		logFatalf("no file specified")
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	requiredEnvs := []string{
		"GRAPH_DB_ENDPOINT",
		"TWO_WAY_KV_ENDPOINT",
	}
	for _, v := range requiredEnvs {
		if os.Getenv(v) == "" {
			logFatalf("'%s' was not set", v)
		} else {
			// print out config
			logMsg("%s=%s", v, os.Getenv(v))
		}
	}

	numberVars := []string{"MAX_APPROX_NODES"}
	for _, e := range numberVars {
		i, err := strconv.Atoi(os.Getenv(e))
		if err != nil {
			logFatalf("Could not parse %s for env variable %s. Reccieve: %v", e, os.Getenv(e), err.Error())
		}
		if i < 1 && i != -1 {
			logFatalf("%s must be greater than 1 but was '%i'", e, i)
		}

	}
}

// runs with metrics
func run(fileLocation string) {
	// assert environment
	parseEnv(fileLocation)
	// start parser
	Parse(fileLocation)
}

func main() {
	app := cli.NewApp()
	app.Name = "textToGraph"
	app.Usage = " a customizable indexer for .txt documents into a biggraph database"
	app.Description = "Parses each word in large files delimited by an empty string, and indexes them into a large scale graph database"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "parse",
			Aliases: []string{"p"},
			Usage:   "parse all lines in text file and add each word to graph",
			Action: func(c *cli.Context) error {
				run(c.Args().Get(0))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
