package main

import (
	_ "embed"
	"github.com/choria-io/fisk"
	"log"
	"os"
	"regexp"
	"time"
)

const appName = "paws"

//go:embed .version
var fullVersion string

//go:embed examples.md
var cheatExamples string

//go:embed services.md
var cheatServices string

//go:embed examples/docker-mysql.sh
var cheatDocker string

var (
	ax *fisk.Application

	quiet       bool
	dots        bool
	timeout     time.Duration
	wait        time.Duration
	delay       time.Duration
	servicesRaw []string
	host        string

	idleTime      time.Duration
	failOnPattern []*regexp.Regexp
	okOnPattern   []*regexp.Regexp
	file          string
	fileGiven     bool
	passThrough   bool
	printMatches  bool
)

func main() {
	log.SetOutput(os.Stderr)
	ax = fisk.New(appName, "Port Availability Waiting System: Before you leap, let PAWS take a peep."+
		"\n\nUse --help-long for more info and --cheats for examples"+"\n\n"+
		"> Anonymous kitten: 'Curiosity checked the port!'").
		Version(fullVersion).
		Author("Copyright 2023 - METATEXX GmbH authors <kontakt@metatexx.de>")
	portCheckCMD := ax.Command("portcheck", "checks the given port list").Alias("ports").Default()
	portCheckCMD.Flag("host", "Default host if none is given in the ports").Short('h').
		Default("localhost").StringVar(&host)
	portCheckCMD.Flag("quiet", "Don't output anything (only return code)").Short('q').
		UnNegatableBoolVar(&quiet)
	portCheckCMD.Flag("progress", "Output dots while scanning (one dot for each port check)").
		Short('p').UnNegatableBoolVar(&dots)
	portCheckCMD.Flag("wait", "Waiting time to wait for the ports to appear").Short('w').
		DurationVar(&wait)
	portCheckCMD.Flag("timeout", "Timeout for all connect and read timeouts").Short('t').
		Default("250ms").DurationVar(&timeout)
	portCheckCMD.Flag("delay", "Minimal time between checks when waiting").Short('d').
		Default("250ms").DurationVar(&delay)
	portCheckCMD.Arg("services", "the services (ports) to check").
		PlaceHolder("(service(-udp|-tcp):)((host):)port or service(-udp|-tcp)://(user(:pass))@(host):port").
		Help("Services to scan (see 'paws cheat examples' for more information))").
		Required().StringsVar(&servicesRaw)

	fileCheckCMD := ax.Command("filecheck", "checks stdin or a file to be idle or meeting other"+
		" conditions (see flags)").
		Alias("log")
	fileCheckCMD.Flag("quiet", "Don't output status information (only return code)").
		Short('q').UnNegatableBoolVar(&quiet)
	fileCheckCMD.Flag("pass-through", "Passes the data to stdout while reading").
		Short('p').UnNegatableBoolVar(&passThrough)
	fileCheckCMD.Flag("print-matches", "Print lines that match the patterns and highlights them."+
		" Will only show the first match per pattern.").
		Short('P').UnNegatableBoolVar(&printMatches)
	fileCheckCMD.Flag("file", "The file to read (stdin if no file is given)").
		IsSetByUser(&fileGiven).ExistingFileVar(&file)
	fileCheckCMD.Flag("timeout", "Timeout after the program returns a failure in any case").
		Short('t').Default("5s").DurationVar(&timeout)
	fileCheckCMD.Flag("idle", "When idle time is given, the program returns without failure if there"+
		" is no data for this amount of time").
		Short('i').DurationVar(&idleTime)
	fileCheckCMD.Flag("failure", "Regular expression(s) to stop and fail if it is detected in the output"+
		" (either one must be found)").
		Short('F').RegexpListVar(&failOnPattern)
	fileCheckCMD.Flag("success", "Regular expression(s) that stops watching and fail if it is detected in"+
		" the output (all must be found)").
		Short('S').RegexpListVar(&okOnPattern)

	ax.Cheat("examples", cheatExamples)
	ax.Cheat("docker", cheatDocker)
	ax.Cheat("services", cheatServices)

	cmd := ax.MustParseWithUsage(os.Args[1:])

	var rc int
	switch cmd {
	case portCheckCMD.FullCommand():
		rc = fileCheck()
	case fileCheckCMD.FullCommand():
		rc = logCheck()
	}
	os.Exit(rc)
}
