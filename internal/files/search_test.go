package files

import (
	"flag"
	"io/fs"
	"mp3/internal"
	"reflect"
	"testing"
)

func Test_readDirectory(t *testing.T) {
	fnName := "readDirectory()"
	// generate test data
	topDir := "loadTest"
	if err := internal.Mkdir(topDir); err != nil {
		t.Errorf("%s error creating %s: %v", fnName, topDir, err)
	}
	defer func() {
		internal.DestroyDirectoryForTesting(fnName, topDir)
	}()
	type args struct {
		dir string
	}
	tests := []struct {
		name      string
		args      args
		wantFiles []fs.FileInfo
		wantOk    bool
		internal.WantedOutput
	}{
		{name: "default", args: args{topDir}, wantFiles: []fs.FileInfo{}, wantOk: true},
		{
			name: "non-existent dir",
			args: args{"non-existent directory"},
			WantedOutput: internal.WantedOutput{
				WantErrorOutput: "The directory \"non-existent directory\" cannot be read: open non-existent directory: The system cannot find the file specified.\n",
				WantLogOutput:   "level='warn' directory='non-existent directory' error='open non-existent directory: The system cannot find the file specified.' msg='cannot read directory'\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := internal.NewOutputDeviceForTesting()
			gotFiles, gotOk := readDirectory(o, tt.args.dir)
			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("readDirectory() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
			}
			if gotOk != tt.wantOk {
				t.Errorf("readDirectory() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if issues, ok := o.CheckOutput(tt.WantedOutput); !ok {
				for _, issue := range issues {
					t.Errorf("readDirectory() %s", issue)
				}
			}
		})
	}
}

func TestSearch_FilterArtists(t *testing.T) {
	fnName := "Search.FilterArtists()"
	// generate test data
	topDir := "loadTest"
	if err := internal.Mkdir(topDir); err != nil {
		t.Errorf("%s error creating %s: %v", fnName, topDir, err)
	}
	defer func() {
		internal.DestroyDirectoryForTesting(fnName, topDir)
	}()
	if err := internal.PopulateTopDirForTesting(topDir); err != nil {
		t.Errorf("%s error populating %s: %v", fnName, topDir, err)
	}
	realFlagSet := flag.NewFlagSet("real", flag.ContinueOnError)
	realS, _ := NewSearchFlags(internal.EmptyConfiguration(), realFlagSet).ProcessArgs(
		internal.NewOutputDeviceForTesting(), []string{"-topDir", topDir})
	realArtists, _ := realS.LoadData(internal.NewOutputDeviceForTesting())
	overFilteredS, _ := NewSearchFlags(internal.EmptyConfiguration(),
		flag.NewFlagSet("overFiltered", flag.ContinueOnError)).ProcessArgs(
		internal.NewOutputDeviceForTesting(), []string{"-topDir", topDir, "-artistFilter", "^Filter all out$"})
	a, _ := realS.LoadUnfilteredData(internal.NewOutputDeviceForTesting())
	type args struct {
		unfilteredArtists []*Artist
	}
	tests := []struct {
		name        string
		s           *Search
		args        args
		wantArtists []*Artist
		wantOk      bool
		internal.WantedOutput
	}{
		{
			name:        "default",
			s:           realS,
			args:        args{unfilteredArtists: a},
			wantArtists: realArtists,
			wantOk:      true,
			WantedOutput: internal.WantedOutput{
				WantLogOutput: "level='info' -albumFilter='.*' -artistFilter='.*' -ext='.mp3' -topDir='loadTest' msg='filtering music files'\n",
			},
		},
		{
			name: "all filtered out",
			s:    overFilteredS,
			args: args{unfilteredArtists: a},
			WantedOutput: internal.WantedOutput{
				WantErrorOutput: "No music files could be found using the specified parameters.\n",
				WantLogOutput: "level='info' -albumFilter='.*' -artistFilter='^Filter all out$' -ext='.mp3' -topDir='loadTest' msg='filtering music files'\n" +
					"level='warn' -albumFilter='.*' -artistFilter='^Filter all out$' -ext='.mp3' -topDir='loadTest' msg='cannot find any artist directories'\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := internal.NewOutputDeviceForTesting()
			gotArtists, gotOk := tt.s.FilterArtists(o, tt.args.unfilteredArtists)
			if !reflect.DeepEqual(gotArtists, tt.wantArtists) {
				t.Errorf("%s = %v, want %v", fnName, gotArtists, tt.wantArtists)
			}
			if gotOk != tt.wantOk {
				t.Errorf("%s ok = %v, want %v", fnName, gotOk, tt.wantOk)
			}
			if issues, ok := o.CheckOutput(tt.WantedOutput); !ok {
				for _, issue := range issues {
					t.Errorf("%s %s", fnName, issue)
				}
			}
		})
	}
}

func TestSearch_LoadData(t *testing.T) {
	fnName := "Search.LoadData()"
	// generate test data
	topDir := "loadTest"
	if err := internal.Mkdir(topDir); err != nil {
		t.Errorf("%s error creating %s: %v", fnName, topDir, err)
	}
	defer func() {
		internal.DestroyDirectoryForTesting(fnName, topDir)
	}()
	if err := internal.PopulateTopDirForTesting(topDir); err != nil {
		t.Errorf("%s error populating %s: %v", fnName, topDir, err)
	}
	tests := []struct {
		name        string
		s           *Search
		wantArtists []*Artist
		wantOk      bool
		internal.WantedOutput
	}{
		{
			name:        "read all",
			s:           CreateFilteredSearchForTesting(topDir, "^.*$", "^.*$"),
			wantArtists: CreateAllArtistsForTesting(topDir, false),
			wantOk:      true,
			WantedOutput: internal.WantedOutput{
				WantLogOutput: "level='info' -albumFilter='^.*$' -artistFilter='^.*$' -ext='.mp3' -topDir='loadTest' msg='reading filtered music files'\n",
			},
		},
		{
			name:        "read with filtering",
			s:           CreateFilteredSearchForTesting(topDir, "^.*[13579]$", "^.*[02468]$"),
			wantArtists: CreateAllOddArtistsWithEvenAlbumsForTesting(topDir),
			wantOk:      true,
			WantedOutput: internal.WantedOutput{
				WantLogOutput: "level='info' -albumFilter='^.*[02468]$' -artistFilter='^.*[13579]$' -ext='.mp3' -topDir='loadTest' msg='reading filtered music files'\n",
			},
		},
		{
			name: "read with all artists filtered out",
			s:    CreateFilteredSearchForTesting(topDir, "^.*X$", "^.*$"),
			WantedOutput: internal.WantedOutput{
				WantErrorOutput: "No music files could be found using the specified parameters.\n",
				WantLogOutput: "level='info' -albumFilter='^.*$' -artistFilter='^.*X$' -ext='.mp3' -topDir='loadTest' msg='reading filtered music files'\n" +
					"level='warn' -albumFilter='^.*$' -artistFilter='^.*X$' -ext='.mp3' -topDir='loadTest' msg='cannot find any artist directories'\n",
			},
		},
		{
			name: "read with all albums filtered out",
			s:    CreateFilteredSearchForTesting(topDir, "^.*$", "^.*X$"),
			WantedOutput: internal.WantedOutput{
				WantErrorOutput: "No music files could be found using the specified parameters.\n",
				WantLogOutput: "level='info' -albumFilter='^.*X$' -artistFilter='^.*$' -ext='.mp3' -topDir='loadTest' msg='reading filtered music files'\n" +
					"level='warn' -albumFilter='^.*X$' -artistFilter='^.*$' -ext='.mp3' -topDir='loadTest' msg='cannot find any artist directories'\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := internal.NewOutputDeviceForTesting()
			gotArtists, gotOk := tt.s.LoadData(o)
			if !reflect.DeepEqual(gotArtists, tt.wantArtists) {
				t.Errorf("Search.LoadData() = %v, want %v", gotArtists, tt.wantArtists)
			}
			if issues, ok := o.CheckOutput(tt.WantedOutput); !ok {
				for _, issue := range issues {
					t.Errorf("Search.LoadData() %s", issue)
				}
			}
			if gotOk != tt.wantOk {
				t.Errorf("Search.LoadData() ok = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestSearch_LoadUnfilteredData(t *testing.T) {
	fnName := "Search.LoadUnfilteredData()"
	// generate test data
	topDir := "loadTest"
	emptyDir := "empty directory"
	defer func() {
		internal.DestroyDirectoryForTesting(fnName, topDir)
		internal.DestroyDirectoryForTesting(fnName, emptyDir)
	}()
	if err := internal.Mkdir(topDir); err != nil {
		t.Errorf("%s error creating %s: %v", fnName, topDir, err)
	}
	if err := internal.PopulateTopDirForTesting(topDir); err != nil {
		t.Errorf("%s error populating %s: %v", fnName, topDir, err)
	}
	if err := internal.Mkdir(emptyDir); err != nil {
		t.Errorf("%s error creating %s: %v", fnName, emptyDir, err)
	}
	tests := []struct {
		name        string
		s           *Search
		wantArtists []*Artist
		wantOk      bool
		internal.WantedOutput
	}{
		{
			name:        "read all",
			s:           CreateSearchForTesting(topDir),
			wantArtists: CreateAllArtistsForTesting(topDir, true),
			wantOk:      true,
			WantedOutput: internal.WantedOutput{
				WantLogOutput: "level='info' -ext='.mp3' -topDir='loadTest' msg='reading unfiltered music files'\n",
			},
		},
		{
			name: "empty dir",
			s:    CreateSearchForTesting(emptyDir),
			WantedOutput: internal.WantedOutput{
				WantErrorOutput: "No music files could be found using the specified parameters.\n",
				WantLogOutput: "level='info' -ext='.mp3' -topDir='empty directory' msg='reading unfiltered music files'\n" +
					"level='warn' -ext='.mp3' -topDir='empty directory' msg='cannot find any artist directories'\n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArtists []*Artist
			var gotOk bool
			o := internal.NewOutputDeviceForTesting()
			if gotArtists, gotOk = tt.s.LoadUnfilteredData(o); !reflect.DeepEqual(gotArtists, tt.wantArtists) {
				t.Errorf("Search.LoadUnfilteredData() = %v, want %v", gotArtists, tt.wantArtists)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Search.LoadUnfilteredData() ok = %v, want %v", gotOk, tt.wantOk)
			}
			if issues, ok := o.CheckOutput(tt.WantedOutput); !ok {
				for _, issue := range issues {
					t.Errorf("Search.LoadUnfilteredData() %s", issue)
				}
			}
		})
	}
}
