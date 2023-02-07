package commands

import (
	"mp3/internal"
	"os"
	"path/filepath"
	"testing"

	tools "github.com/majohn-r/cmd-toolkit"
	"github.com/majohn-r/output"
)

func Test_markDirty(t *testing.T) {
	const fnName = "markDirty()"
	emptyDir := "empty"
	if err := tools.Mkdir(emptyDir); err != nil {
		t.Errorf("%s error creating %q: %v", fnName, emptyDir, err)
	}
	filledDir := "filled"
	if err := tools.Mkdir(filledDir); err != nil {
		t.Errorf("%s error creating %q: %v", fnName, filledDir, err)
	}
	if err := internal.CreateFileForTesting(filledDir, dirtyFileName); err != nil {
		t.Errorf("%s error creating %q: %v", fnName, filepath.Join(filledDir, dirtyFileName), err)
	}
	defer func() {
		internal.DestroyDirectoryForTesting(fnName, emptyDir)
		internal.DestroyDirectoryForTesting(fnName, filledDir)
	}()
	type args struct {
		cmd string
	}
	tests := map[string]struct {
		appPath string
		args
		output.WantedRecording
	}{
		"typical first use": {
			appPath: emptyDir,
			args:    args{cmd: "calling command"},
			WantedRecording: output.WantedRecording{
				Log: "level='info' fileName='empty\\metadata.dirty' msg='metadata dirty file written'\n",
			},
		},
		"typical second use": {
			appPath: filledDir,
			args:    args{cmd: "calling command"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tools.SetApplicationPath(tt.appPath)
			markDirty(o, tt.args.cmd)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("%s %s", fnName, issue)
				}
			}
		})
	}
}

func Test_dirty(t *testing.T) {
	const fnName = "dirty()"
	testDir := "dirty"
	if err := tools.Mkdir(testDir); err != nil {
		t.Errorf("%s error creating %q: %v", fnName, testDir, err)
	}
	oldAppPath := tools.SetApplicationPath(testDir)
	defer func() {
		internal.DestroyDirectoryForTesting(fnName, testDir)
		tools.SetApplicationPath(oldAppPath)
	}()
	tests := map[string]struct {
		want bool
	}{
		"file definitively exists":       {want: true},
		"file definitely does not exist": {want: false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.want {
				if err := internal.CreateFileForTestingWithContent(testDir, dirtyFileName, []byte("dirty")); err != nil {
					t.Errorf("%s error creating %q: %v", fnName, dirtyFileName, err)
				}
			} else {
				os.Remove(filepath.Join(testDir, dirtyFileName))
			}
			if gotDirty := dirty(); gotDirty != tt.want {
				t.Errorf("%s = %t, want %t", fnName, gotDirty, tt.want)
			}
		})
	}
}

func Test_clearDirty(t *testing.T) {
	const fnName = "clearDirty()"
	testDir := "clearDirty"
	if err := tools.Mkdir(testDir); err != nil {
		t.Errorf("%s error creating %q: %v", fnName, testDir, err)
	}
	oldAppPath := tools.ApplicationPath()
	if err := internal.CreateFileForTestingWithContent(testDir, dirtyFileName, []byte("dirty")); err != nil {
		t.Errorf("%s error creating %q: %v", fnName, dirtyFileName, err)
	}
	// create another file structure with a dirty file that is open for reading
	uncleanable := "clearDirty2"
	if err := tools.Mkdir(uncleanable); err != nil {
		t.Errorf("%s error creating %q: %v", fnName, uncleanable, err)
	}
	if err := internal.CreateFileForTestingWithContent(uncleanable, dirtyFileName, []byte("dirty")); err != nil {
		t.Errorf("%s error creating second %q: %v", fnName, dirtyFileName, err)
	}
	f, err := os.Open(filepath.Join(uncleanable, dirtyFileName))
	if err != nil {
		t.Errorf("%s error opening second %q: %v", fnName, dirtyFileName, err)
	}
	defer func() {
		tools.SetApplicationPath(oldAppPath)
		internal.DestroyDirectoryForTesting(fnName, testDir)
		f.Close()
		internal.DestroyDirectoryForTesting(fnName, uncleanable)
	}()
	tests := map[string]struct {
		initialDirtyFolder string
		output.WantedRecording
	}{
		"successful removal": {
			initialDirtyFolder: testDir,
			WantedRecording: output.WantedRecording{
				Log: "level='info' fileName='clearDirty\\metadata.dirty' msg='metadata dirty file deleted'\n",
			},
		},
		"nothing to remove": {
			initialDirtyFolder: ".",
		},
		"unremovable file": {
			initialDirtyFolder: uncleanable,
			WantedRecording: output.WantedRecording{
				Error: "The file \"clearDirty2\\\\metadata.dirty\" cannot be deleted: remove clearDirty2\\metadata.dirty: The process cannot access the file because it is being used by another process.\n",
				Log:   "level='error' error='remove clearDirty2\\metadata.dirty: The process cannot access the file because it is being used by another process.' fileName='clearDirty2\\metadata.dirty' msg='cannot delete file'\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tools.SetApplicationPath(tt.initialDirtyFolder)
			o := output.NewRecorder()
			clearDirty(o)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("%s %s", fnName, issue)
				}
			}
		})
	}
}
