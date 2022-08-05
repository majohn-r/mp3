package commands

import (
	"flag"
	"mp3/internal"
	"mp3/internal/files"
	"os"
	"sort"
)

type postrepair struct {
	n  string
	sf *files.SearchFlags
}

func (p *postrepair) name() string {
	return p.n
}

func newPostRepairCommand(o internal.OutputBus, c *internal.Configuration, fSet *flag.FlagSet) (*postrepair, bool) {
	sFlags, sFlagsOk := files.NewSearchFlags(o, c, fSet)
	if sFlagsOk {
		return &postrepair{
			n:  fSet.Name(),
			sf: sFlags,
		}, true
	}
	return nil, false
}

func newPostRepair(o internal.OutputBus, c *internal.Configuration, fSet *flag.FlagSet) (CommandProcessor, bool) {
	return newPostRepairCommand(o, c, fSet)
}

func (p *postrepair) Exec(o internal.OutputBus, args []string) (ok bool) {
	if s, argsOk := p.sf.ProcessArgs(o, args); argsOk {
		ok = p.runCommand(o, s)
	}
	return
}

func (p *postrepair) logFields() map[string]interface{} {
	return map[string]interface{}{fkCommandName: p.name()}
}

func (p *postrepair) runCommand(o internal.OutputBus, s *files.Search) (ok bool) {
	o.LogWriter().Info(internal.LI_EXECUTING_COMMAND, p.logFields())
	artists, ok := s.LoadData(o)
	if ok {
		backups := make(map[string]*files.Album)
		var backupDirectories []string
		for _, artist := range artists {
			for _, album := range artist.Albums() {
				backupDirectory := album.BackupDirectory()
				if internal.DirExists(backupDirectory) {
					backupDirectories = append(backupDirectories, backupDirectory)
					backups[backupDirectory] = album
				}
			}
		}
		if len(backupDirectories) == 0 {
			o.WriteConsole(true, "There are no backup directories to delete")
		} else {
			sort.Strings(backupDirectories)
			for _, backupDirectory := range backupDirectories {
				removeBackupDirectory(o, backupDirectory, backups[backupDirectory])
			}
		}
	}
	return
}

func removeBackupDirectory(o internal.OutputBus, d string, a *files.Album) {
	if err := os.RemoveAll(d); err != nil {
		o.LogWriter().Error(internal.LE_CANNOT_DELETE_DIRECTORY, map[string]interface{}{
			internal.FK_DIRECTORY: d,
			internal.FK_ERROR:     err,
		})
		o.WriteError(internal.USER_CANNOT_DELETE_DIRECTORY, d, err)
	} else {
		o.WriteConsole(true, "The backup directory for artist %q album %q has been deleted\n", a.RecordingArtistName(), a.Name())
	}
}
