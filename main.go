package main

import (
	"github.com/urfave/cli"
	"os"
	 log "github.com/sirupsen/logrus"
)

// checks environment for required env vars
var logFatalf = log.Fatalf
var logMsg = log.Infof

func parseEnv() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	requiredEnvs := []string{
		"GRAPH_DB_ENDPOINT",
		"MAX_APPROX_NODES",
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
}

// runs with metrics
func run() {
	// assert environment
	parseEnv()
	// start parser
}

func main() {
	app := cli.NewApp()
	app.Name = "crawler"
	app.Usage = " acustomizable web crawler script for different websites"
	app.Description = "web crawl different URLs and add similar urls to a graph database"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "parse",
			Aliases: []string{"p"},
			Usage:   "parse all lines in text file",
			Action: func(c *cli.Context) error {
				run()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
