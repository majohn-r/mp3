/*
Copyright © 2021 Marc Johnson (marc.johnson27591@gmail.com)
*/
package cmd

import (
	"fmt"
	"mp3/internal/files"
	"sort"

	"github.com/majohn-r/output"
	"github.com/spf13/cobra"
)

const postRepairCommandName = "postRepair"

var (
	// PostRepairCmd represents the postRepair command
	PostRepairCmd = &cobra.Command{
		Use:                   postRepairCommandName + " " + searchUsage,
		DisableFlagsInUseLine: true,
		Short: "Deletes the backup directories, and their contents, created" +
			" by the " + repairCommandName + " command",
		Long: fmt.Sprintf(
			"%q deletes the backup directories (and their contents) created by the %q command",
			postRepairCommandName, repairCommandName),
		RunE: PostRepairRun,
	}
)

func PostRepairRun(cmd *cobra.Command, _ []string) error {
	exitError := NewExitProgrammingError(postRepairCommandName)
	o := getBus()
	producer := cmd.Flags()
	ss, searchFlagsOk := EvaluateSearchFlags(o, producer)
	if searchFlagsOk {
		// do some work here!
		LogCommandStart(o, postRepairCommandName, ss.Values())
		allArtists, loaded := ss.Load(o)
		exitError = PostRepairWork(o, ss, allArtists, loaded)
	}
	return ToErrorInterface(exitError)
}

func PostRepairWork(o output.Bus, ss *SearchSettings, allArtists []*files.Artist,
	loaded bool) (e *ExitError) {
	e = NewExitUserError(postRepairCommandName)
	if loaded {
		if filteredArtists, filtered := ss.Filter(o, allArtists); filtered {
			e = nil
			dirCount := 0
			for _, artist := range filteredArtists {
				dirCount += len(artist.Albums())
			}
			dirs := make([]string, 0, dirCount)
			for _, artist := range filteredArtists {
				for _, album := range artist.Albums() {
					dir := album.BackupDirectory()
					if DirExists(dir) {
						dirs = append(dirs, dir)
					}
				}
			}
			o.WriteCanonicalConsole("Backup directories to delete: %d", len(dirs))
			if len(dirs) > 0 {
				sort.Strings(dirs)
				dirsDeleted := 0
				for _, dir := range dirs {
					if RemoveBackupDirectory(o, dir) {
						dirsDeleted++
					} else {
						e = NewExitSystemError(postRepairCommandName)
					}
				}
				o.WriteCanonicalConsole("Backup directories deleted: %d", dirsDeleted)
			}
		}
	}
	return
}

func RemoveBackupDirectory(o output.Bus, dir string) bool {
	if err := RemoveAll(dir); err != nil {
		o.Log(output.Error, "cannot delete directory", map[string]any{
			"directory": dir,
			"error":     err,
		})
		return false
	}
	o.Log(output.Info, "directory deleted", map[string]any{"directory": dir})
	return true
}

func init() {
	RootCmd.AddCommand(PostRepairCmd)
	bus := getBus()
	c := getConfiguration()
	AddFlags(bus, c, PostRepairCmd.Flags(), SearchFlags)
}
