package commands

import (
	"flag"
	"fmt"
	"mp3/internal"
	"sort"
	"strings"
)

const (
	checkCommand      = "check"
	fkCommandName     = "command"
	fkCount           = "count"
	lsCommand         = "ls"
	postRepairCommand = "postRepair"
	repairCommand     = "repair"
)

type CommandProcessor interface {
	name() string
	Exec(internal.OutputBus, []string) bool
}

type commandInitializer struct {
	name           string
	defaultCommand bool
	initializer    func(*internal.Configuration, *flag.FlagSet) CommandProcessor
}

// ProcessCommand selects which command to be run and returns the relevant
// CommandProcessor, command line arguments and ok status
func ProcessCommand(o internal.OutputBus, args []string) (cmd CommandProcessor, cmdArgs []string, ok bool) {
	var c *internal.Configuration
	if c, ok = internal.ReadConfigurationFile(o); !ok {
		return nil, nil, false
	}
	var defaultSettings map[string]bool
	if defaultSettings, ok = getDefaultSettings(o, c.SubConfiguration("command")); !ok {
		return nil, nil, false
	}
	var initializers []commandInitializer
	lsCommand := commandInitializer{
		name:           lsCommand,
		defaultCommand: defaultSettings[lsCommand],
		initializer:    newLs,
	}
	checkCommand := commandInitializer{
		name:           checkCommand,
		defaultCommand: defaultSettings[checkCommand],
		initializer:    newCheck,
	}
	repairCommand := commandInitializer{
		name:           repairCommand,
		defaultCommand: defaultSettings[repairCommand],
		initializer:    newRepair,
	}
	postRepairCommand := commandInitializer{
		name:           postRepairCommand,
		defaultCommand: defaultSettings[postRepairCommand],
		initializer:    newPostRepair,
	}
	initializers = append(initializers, lsCommand, checkCommand, repairCommand, postRepairCommand)
	cmd, cmdArgs, ok = selectCommand(o, c, initializers, args)
	return
}

func getDefaultSettings(o internal.OutputBus, c *internal.Configuration) (m map[string]bool, ok bool) {
	defaultCommand, ok := c.StringValue("default")
	if !ok { // no definition
		m = map[string]bool{
			checkCommand:      false,
			lsCommand:         true,
			postRepairCommand: false,
			repairCommand:     false,
		}
		ok = true
		return
	}
	m = map[string]bool{
		checkCommand:      defaultCommand == checkCommand,
		lsCommand:         defaultCommand == lsCommand,
		postRepairCommand: defaultCommand == postRepairCommand,
		repairCommand:     defaultCommand == repairCommand,
	}
	found := false
	for _, value := range m {
		if value {
			found = true
			break
		}
	}
	if !found {
		o.Log(internal.WARN, internal.LW_INVALID_DEFAULT_COMMAND, map[string]interface{}{
			fkCommandName: defaultCommand,
		})
		fmt.Fprintf(o.ErrorWriter(), internal.USER_INVALID_DEFAULT_COMMAND, defaultCommand)
		m = nil
		ok = false
		return
	}
	ok = true
	return
}

func selectCommand(o internal.OutputBus, c *internal.Configuration, i []commandInitializer, args []string) (cmd CommandProcessor, callingArgs []string, ok bool) {
	if len(i) == 0 {
		o.Log(internal.ERROR, internal.LE_COMMAND_COUNT, map[string]interface{}{
			fkCount: 0,
		})
		fmt.Fprint(o.ErrorWriter(), internal.USER_NO_COMMANDS_DEFINED)
		return
	}
	var defaultInitializers int
	var defaultInitializerName string
	for _, initializer := range i {
		if initializer.defaultCommand {
			defaultInitializers++
			defaultInitializerName = initializer.name
		}
	}
	if defaultInitializers != 1 {
		o.Log(internal.ERROR, internal.LE_DEFAULT_COMMAND_COUNT, map[string]interface{}{
			fkCount: defaultInitializers,
		})
		fmt.Fprintf(o.ErrorWriter(), internal.USER_INCORRECT_NUMBER_OF_DEFAULT_COMMANDS_DEFINED, defaultInitializers)
		return
	}
	processorMap := make(map[string]CommandProcessor)
	for _, commandInitializer := range i {
		fSet := flag.NewFlagSet(commandInitializer.name, flag.ContinueOnError)
		processorMap[commandInitializer.name] = commandInitializer.initializer(c, fSet)
	}
	if len(args) < 2 {
		cmd = processorMap[defaultInitializerName]
		callingArgs = []string{defaultInitializerName}
		ok = true
		return
	}
	commandName := args[1]
	if strings.HasPrefix(commandName, "-") {
		cmd = processorMap[defaultInitializerName]
		callingArgs = args[1:]
		ok = true
		return
	}
	cmd, found := processorMap[commandName]
	if !found {
		cmd = nil
		callingArgs = nil
		o.Log(internal.WARN, internal.LW_UNRECOGNIZED_COMMAND, map[string]interface{}{
			fkCommandName: commandName,
		})
		var commandNames []string
		for _, initializer := range i {
			commandNames = append(commandNames, initializer.name)
		}
		sort.Strings(commandNames)
		fmt.Fprintf(o.ErrorWriter(), internal.USER_NO_SUCH_COMMAND, commandName, commandNames)
		return
	}
	callingArgs = args[2:]
	ok = true
	return
}
