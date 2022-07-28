package main

import (
	"mp3/internal"
	"mp3/internal/commands"
	"os"
	"time"
)

// these variables' values are injected by the mage build
var (
	// semantic version; read by the mage build from version.txt
	version string = "unknown version!"
	// build timestamp in RFC3339 format (2006-01-02T15:04:05Z07:00)
	creation string
)

func main() {
	os.Exit(exec(internal.InitLogging, os.Args))
}

const (
	fkCommandLineArguments = "args"
	fkDuration             = "duration"
	fkExitCode             = "exitCode"
	fkTimeStamp            = "timeStamp"
	fkVersion              = "version"
	statusFormat           = "%q version %s, created at %s, failed"
)

func exec(logInit func(internal.OutputBus) bool, cmdLine []string) (returnValue int) {
	returnValue = 1
	o := internal.NewOutputDevice()
	if logInit(o) {
		returnValue = run(o, cmdLine)
	}
	report(o, returnValue)
	return
}

func report(o internal.OutputBus, returnValue int) {
	if returnValue != 0 {
		o.WriteError(statusFormat, internal.AppName, version, creation)
	}
}

func run(o internal.OutputBus, cmdlineArgs []string) (returnValue int) {
	returnValue = 1
	startTime := time.Now()
	o.LogWriter().Info(internal.LI_BEGIN_EXECUTION, map[string]interface{}{
		fkVersion:              version,
		fkTimeStamp:            creation,
		fkCommandLineArguments: cmdlineArgs,
	})
	if cmd, args, ok := commands.ProcessCommand(o, cmdlineArgs); ok {
		if cmd.Exec(o, args) {
			returnValue = 0
		}
	}
	o.LogWriter().Info(internal.LI_END_EXECUTION, map[string]interface{}{
		fkDuration: time.Since(startTime),
		fkExitCode: returnValue,
	})
	return
}
