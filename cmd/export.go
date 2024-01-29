/*
Copyright © 2021 Marc Johnson (marc.johnson27591@gmail.com)
*/
package cmd

import (
	"fmt"
	"path/filepath"

	cmd_toolkit "github.com/majohn-r/cmd-toolkit"
	"github.com/majohn-r/output"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	ExportCommand         = "export"
	ExportFlagDefaults    = "defaults"
	exportDefaultsAsFlag  = "--" + ExportFlagDefaults
	exportDefaultsCure    = "What to do:\nUse either '" + exportDefaultsAsFlag + "' or '" + exportDefaultsAsFlag + "=true' to enable exporting defaults"
	ExportFlagOverwrite   = "overwrite"
	exportOverwriteAsFlag = "--" + ExportFlagOverwrite
	exportOverwriteCure   = "What to do:\nUse either '" + exportOverwriteAsFlag + "' or '" + exportOverwriteAsFlag + "=true' to enable overwriting the existing file"
)

// ExportCmd represents the export command
var (
	ExportCmd = &cobra.Command{
		Use:                   ExportCommand + " [" + exportDefaultsAsFlag + "] [" + exportOverwriteAsFlag + "]",
		DisableFlagsInUseLine: true,
		Short:                 "Exports default program configuration data",
		Long:                  fmt.Sprintf("%q", ExportCommand) + ` exports default program configuration data to %APPDATA%\mp3\defaults.yaml`,
		Example: ExportCommand + " " + exportDefaultsAsFlag + "\n" +
			"  Write default program configuration data\n" +
			ExportCommand + " " + exportOverwriteAsFlag + "\n" +
			"  Overwrite a pre-existing defaults.yaml file",
		Run: ExportRun,
	}
	ExportFlags = SectionFlags{
		SectionName: ExportCommand,
		Flags: map[string]*FlagDetails{
			ExportFlagDefaults: {
				AbbreviatedName: "d",
				Usage:           "write default program configuration data",
				ExpectedType:    BoolType,
				DefaultValue:    false,
			},
			ExportFlagOverwrite: {
				AbbreviatedName: "o",
				Usage:           "overwrite existing file",
				ExpectedType:    BoolType,
				DefaultValue:    false,
			},
		},
	}
	defaultConfigurationSettings = map[string]map[string]any{}
)

func addDefaults(sf SectionFlags) {
	payload := map[string]any{}
	for flag, details := range sf.Flags {
		if bounded, ok := details.DefaultValue.(*cmd_toolkit.IntBounds); ok {
			payload[flag] = bounded.Default()
		} else {
			payload[flag] = details.DefaultValue
		}
	}
	defaultConfigurationSettings[sf.SectionName] = payload
}

type ExportFlagSettings struct {
	DefaultsEnabled  bool
	DefaultsSet      bool
	OverwriteEnabled bool
	OverwriteSet     bool
}

func ExportRun(cmd *cobra.Command, _ []string) {
	o := getBus()
	values, eSlice := ReadFlags(cmd.Flags(), ExportFlags)
	fatalError := true // pessimist!
	if ProcessFlagErrors(o, eSlice) {
		settings, ok := ProcessExportFlags(o, values)
		if ok {
			LogCommandStart(o, ExportCommand, map[string]any{
				exportDefaultsAsFlag:  settings.DefaultsEnabled,
				"defaults-user-set":   settings.DefaultsSet,
				exportOverwriteAsFlag: settings.OverwriteEnabled,
				"overwrite-user-set":  settings.OverwriteSet,
			})
			settings.ExportDefaultConfiguration(o)
			fatalError = false
		}
	}
	if fatalError {
		Exit(1)
	}
}

func ProcessExportFlags(o output.Bus, values map[string]*FlagValue) (*ExportFlagSettings, bool) {
	var err error
	result := &ExportFlagSettings{}
	ok := true // optimistic
	result.DefaultsEnabled, result.DefaultsSet, err = GetBool(o, values, ExportFlagDefaults)
	if err != nil {
		ok = false
	}
	result.OverwriteEnabled, result.OverwriteSet, err = GetBool(o, values, ExportFlagOverwrite)
	if err != nil {
		ok = false
	}
	return result, ok
}

func CreateFile(o output.Bus, f string, content []byte) bool {
	if err := WriteFile(f, content, cmd_toolkit.StdFilePermissions); err != nil {
		cmd_toolkit.ReportFileCreationFailure(o, ExportCommand, f, err)
		return false
	}
	o.WriteCanonicalConsole("File %q has been written", f)
	return true
}

func (efs *ExportFlagSettings) ExportDefaultConfiguration(o output.Bus) {
	if efs.CanWriteDefaults(o) {
		// ignoring error return, as we're not marshalling structs, where mischief
		// can occur
		payload, _ := yaml.Marshal(defaultConfigurationSettings)
		path := ApplicationPath()
		f := filepath.Join(path, cmd_toolkit.DefaultConfigFileName())
		if PlainFileExists(f) {
			efs.OverwriteFile(o, f, payload)
		} else {
			CreateFile(o, f, payload)
		}
	}
}

func (efs *ExportFlagSettings) OverwriteFile(o output.Bus, f string, payload []byte) {
	if efs.CanOverwriteFile(o, f) {
		backup := f + "-backup"
		if err := Rename(f, backup); err != nil {
			o.WriteCanonicalError("The file %q cannot be renamed to %q: %v", f, backup, err)
			o.Log(output.Error, "rename failed", map[string]any{
				"error": err,
				"old":   f,
				"new":   backup,
			})
		} else if CreateFile(o, f, payload) {
			Remove(backup)
		}
	}
}

func (efs *ExportFlagSettings) CanOverwriteFile(o output.Bus, f string) (canOverwrite bool) {
	if !efs.OverwriteEnabled {
		o.WriteCanonicalError("The file %q exists and cannot be overwritten", f)
		o.Log(output.Error, "overwrite is not permitted", map[string]any{
			exportOverwriteAsFlag: false,
			"fileName":            f,
			"user-set":            efs.OverwriteSet,
		})
		if efs.OverwriteSet {
			o.WriteCanonicalError("Why?\nYou explicitly set %s false", exportOverwriteAsFlag)
		} else {
			o.WriteCanonicalError("Why?\nAs currently configured, overwriting the file is disabled")
		}
		o.WriteCanonicalError(exportOverwriteCure)
	} else {
		canOverwrite = true
	}
	return
}

func (efs *ExportFlagSettings) CanWriteDefaults(o output.Bus) (canWrite bool) {
	if !efs.DefaultsEnabled {
		o.WriteCanonicalError("default configuration settings will not be exported")
		o.Log(output.Error, "export defaults disabled", map[string]any{
			exportDefaultsAsFlag: false,
			"user-set":           efs.DefaultsSet,
		})
		if efs.DefaultsSet {
			o.WriteCanonicalError("Why?\nYou explicitly set %s false", exportDefaultsAsFlag)
		} else {
			o.WriteCanonicalError("Why?\nAs currently configured, exporting default configuration settings is disabled")
		}
		o.WriteCanonicalError(exportDefaultsCure)
	} else {
		canWrite = true
	}
	return
}

func init() {
	RootCmd.AddCommand(ExportCmd)
	addDefaults(ExportFlags)
	o := getBus()
	c := getConfiguration()
	AddFlags(o, c, ExportCmd.Flags(), ExportFlags, false)
}