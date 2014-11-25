package main

import (
	"github.com/a-sk/go-nagios-log-analyzer/analyze"
	"github.com/a-sk/go-nagios-log-analyzer/parser"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
)

func newApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "Parse nagios log and analyze uptime"

	var logPath string

	app.Action = func(c *cli.Context) {
		logPath = c.String("data-file")
		hosts := c.StringSlice("hosts")
		criticals := c.StringSlice("criticals")
		file, err := os.Open(logPath)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		parsed := parser.ParseLog(file, hosts)
		down := analyze.FindUptime(parsed, criticals, 3)

		jsonEncodedOutput, err := json.Marshal(down)
		if err != nil {
			fmt.Println("Error: could not encode JSON")
		} else {
			fmt.Println(string(jsonEncodedOutput))
		}
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "data-file",
			Value: "nagios.log",
		},
		cli.StringSliceFlag{
			Name:  "hosts",
			Value: &cli.StringSlice{},
			Usage: "list of regexps to filter log file",
		},
		cli.StringSliceFlag{
			Name:  "criticals",
			Value: &cli.StringSlice{},
			Usage: "list of checks to count downtime",
		},
	}
	return app
}

func main() {
	app := newApp()
	app.Run(os.Args)
}
