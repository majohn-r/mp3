package commands

import (
	"bytes"
	"flag"
	"fmt"
	"mp3/internal"
	"mp3/internal/files"
	"os"
	"sort"
	"strings"
	"testing"
)

func Test_ls_validateTrackSorting(t *testing.T) {
	fnName := "ls.validateTrackSorting()"
	tests := []struct {
		name          string
		sortingInput  string
		includeAlbums bool
		wantSorting   string
	}{
		{name: "alpha sorting with albums", sortingInput: "alpha", includeAlbums: true, wantSorting: "alpha"},
		{name: "alpha sorting without albums", sortingInput: "alpha", includeAlbums: false, wantSorting: "alpha"},
		{name: "numeric sorting with albums", sortingInput: "numeric", includeAlbums: true, wantSorting: "numeric"},
		{name: "numeric sorting without albums", sortingInput: "numeric", includeAlbums: false, wantSorting: "alpha"},
		{name: "invalid sorting with albums", sortingInput: "nonsense", includeAlbums: true, wantSorting: "numeric"},
		{name: "invalid sorting without albums", sortingInput: "nonsense", includeAlbums: false, wantSorting: "alpha"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("ls", flag.ContinueOnError)
			lsCommand := newLsSubCommand(internal.EmptyConfiguration(), fs)
			lsCommand.trackSorting = &tt.sortingInput
			lsCommand.includeAlbums = &tt.includeAlbums
			lsCommand.validateTrackSorting()
			if *lsCommand.trackSorting != tt.wantSorting {
				t.Errorf("%s: got %q, want %q", fnName, *lsCommand.trackSorting, tt.wantSorting)
			}
		})
	}
}

type testTrack struct {
	artistName string
	albumName  string
	trackName  string
}

func generateListing(artists, albums, tracks, annotated, sortNumerically bool) string {
	var trackCollection []*testTrack
	for j := 0; j < 10; j++ {
		artist := internal.CreateArtistNameForTesting(j)
		for k := 0; k < 10; k++ {
			album := internal.CreateAlbumNameForTesting(k)
			for m := 0; m < 10; m++ {
				track := internal.CreateTrackNameForTesting(m)
				trackCollection = append(trackCollection, &testTrack{
					artistName: artist,
					albumName:  album,
					trackName:  track,
				})
			}
		}
	}
	var output []string
	switch artists {
	case true:
		tracksByArtist := make(map[string][]*testTrack)
		for _, tt := range trackCollection {
			artistName := tt.artistName
			tracksByArtist[artistName] = append(tracksByArtist[artistName], tt)
		}
		var artistNames []string
		for key := range tracksByArtist {
			artistNames = append(artistNames, key)
		}
		sort.Strings(artistNames)
		for _, artistName := range artistNames {
			output = append(output, fmt.Sprintf("Artist: %s", artistName))
			output = append(output, generateAlbumListings(tracksByArtist[artistName], "  ", artists, albums, tracks, annotated, sortNumerically)...)
		}
	case false:
		output = append(output, generateAlbumListings(trackCollection, "", artists, albums, tracks, annotated, sortNumerically)...)
	}
	if len(output) != 0 {
		output = append(output, "") // force trailing newline
	}
	return strings.Join(output, "\n")
}

type albumType struct {
	artistName string
	albumName  string
}

type albumTypes []albumType

func (a albumTypes) Len() int {
	return len(a)
}

func (a albumTypes) Less(i, j int) bool {
	if a[i].albumName == a[j].albumName {
		return a[i].artistName < a[j].artistName
	}
	return a[i].albumName < a[j].albumName
}

func (a albumTypes) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func generateAlbumListings(testTracks []*testTrack, spacer string, artists, albums, tracks, annotated, sortNumerically bool) []string {
	var output []string
	switch albums {
	case true:
		albumsToList := make(map[albumType][]*testTrack)
		for _, tt := range testTracks {
			albumName := tt.albumName
			var albumTitle string
			if annotated && !artists {
				albumTitle = fmt.Sprintf("%q by %q", albumName, tt.artistName)
			} else {
				albumTitle = albumName
			}
			album := albumType{artistName: tt.artistName, albumName: albumTitle}
			albumsToList[album] = append(albumsToList[album], tt)
		}
		var albumNames albumTypes
		for key := range albumsToList {
			albumNames = append(albumNames, key)
		}

		sort.Sort(albumNames)
		for _, albumTitle := range albumNames {
			output = append(output, fmt.Sprintf("%sAlbum: %s", spacer, albumTitle.albumName))
			output = append(output, generateTrackListings(albumsToList[albumTitle], spacer+"  ", artists, albums, tracks, annotated, sortNumerically)...)
		}
	case false:
		output = append(output, generateTrackListings(testTracks, spacer, artists, albums, tracks, annotated, sortNumerically)...)
	}
	return output
}

func generateTrackListings(testTracks []*testTrack, spacer string, artists, albums, tracks, annotated, sortNumerically bool) []string {
	var output []string
	if tracks {
		var tracksToList []string
		for _, tt := range testTracks {
			trackName, trackNumber := files.ParseTrackNameForTesting(tt.trackName)
			key := trackName
			if annotated {
				if !albums {
					key = fmt.Sprintf("%q on %q by %q", trackName, tt.albumName, tt.artistName)
					if !artists {
					} else {
						key = fmt.Sprintf("%q on %q", trackName, tt.albumName)
					}
				}
			}
			if sortNumerically && albums {
				key = fmt.Sprintf("%2d. %s", trackNumber, trackName)
			}
			tracksToList = append(tracksToList, key)
		}
		sort.Strings(tracksToList)
		for _, trackName := range tracksToList {
			output = append(output, fmt.Sprintf("%s%s", spacer, trackName))
		}
	}
	return output
}

func Test_ls_Exec(t *testing.T) {
	fnName := "ls.Exec()"
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
	type args struct {
		args []string
	}
	tests := []struct {
		name  string
		l     *ls
		args  args
		wantW string
	}{
		{
			name: "no output",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=false",
					"-includeTracks=false",
				},
			},
			wantW: generateListing(false, false, false, false, false),
		},
		// tracks only
		{
			name: "unannotated tracks only",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=false",
					"-includeTracks=true",
					"-annotate=false",
				},
			},
			wantW: generateListing(false, false, true, false, false),
		},
		{
			name: "annotated tracks only",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=false",
					"-includeTracks=true",
					"-annotate=true",
				},
			},
			wantW: generateListing(false, false, true, true, false),
		},
		{
			name: "unannotated tracks only with numeric sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=false",
					"-includeTracks=true",
					"-annotate=false",
					"-sort=numeric",
				},
			},
			wantW: generateListing(false, false, true, false, true),
		},
		{
			name: "annotated tracks only with numeric sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=false",
					"-includeTracks=true",
					"-annotate=true",
					"-sort", "numeric",
				},
			},
			wantW: generateListing(false, false, true, true, true),
		},
		// albums only
		{
			name: "unannotated albums only",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=false",
					"-includeTracks=false",
					"-annotate=false",
				},
			},
			wantW: generateListing(false, true, false, false, false),
		},
		{
			name: "annotated albums only",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=false",
					"-includeTracks=false",
					"-annotate=true",
				},
			},
			wantW: generateListing(false, true, false, true, false),
		},
		// artists only
		{
			name: "unannotated artists only",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=true",
					"-includeTracks=false",
					"-annotate=false",
				},
			},
			wantW: generateListing(true, false, false, false, false),
		},
		{
			name: "annotated artists only",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=true",
					"-includeTracks=false",
					"-annotate=true",
				},
			},
			wantW: generateListing(true, false, false, true, false),
		},
		// albums and artists
		{
			name: "unannotated artists and albums only",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=true",
					"-includeTracks=false",
					"-annotate=false",
				},
			},
			wantW: generateListing(true, true, false, false, false),
		},
		{
			name: "annotated artists and albums only",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=true",
					"-includeTracks=false",
					"-annotate=true",
				},
			},
			wantW: generateListing(true, true, false, true, false),
		},
		// albums and tracks
		{
			name: "unannotated albums and tracks with alpha sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=false",
					"-includeTracks=true",
					"-annotate=false",
					"-sort", "alpha",
				},
			},
			wantW: generateListing(false, true, true, false, false),
		},
		{
			name: "annotated albums and tracks with alpha sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=false",
					"-includeTracks=true",
					"-annotate=true",
					"-sort", "alpha",
				},
			},
			wantW: generateListing(false, true, true, true, false),
		},
		{
			name: "unannotated albums and tracks with numeric sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=false",
					"-includeTracks=true",
					"-annotate=false",
					"-sort", "numeric",
				},
			},
			wantW: generateListing(false, true, true, false, true),
		},
		{
			name: "annotated albums and tracks with numeric sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=false",
					"-includeTracks=true",
					"-annotate=true",
					"-sort", "numeric",
				},
			},
			wantW: generateListing(false, true, true, true, true),
		},
		// artists and tracks
		{
			name: "unannotated artists and tracks with alpha sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=true",
					"-includeTracks=true",
					"-annotate=false",
					"-sort", "alpha",
				},
			},
			wantW: generateListing(true, false, true, false, false),
		},
		{
			name: "annotated artists and tracks with alpha sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=true",
					"-includeTracks=true",
					"-annotate=true",
					"-sort", "alpha",
				},
			},
			wantW: generateListing(true, false, true, true, false),
		},
		{
			name: "unannotated artists and tracks with numeric sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=true",
					"-includeTracks=true",
					"-annotate=false",
					"-sort", "numeric",
				},
			},
			wantW: generateListing(true, false, true, false, true),
		},
		{
			name: "annotated artists and tracks with numeric sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=false",
					"-includeArtists=true",
					"-includeTracks=true",
					"-annotate=true",
					"-sort", "numeric",
				},
			},
			wantW: generateListing(true, false, true, true, true),
		},
		// albums, artists, and tracks
		{
			name: "unannotated artists, albums, and tracks with alpha sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=true",
					"-includeTracks=true",
					"-annotate=false",
					"-sort", "alpha",
				},
			},
			wantW: generateListing(true, true, true, false, false),
		},
		{
			name: "annotated artists, albums, and tracks with alpha sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=true",
					"-includeTracks=true",
					"-annotate=true",
					"-sort", "alpha",
				},
			},
			wantW: generateListing(true, true, true, true, false),
		},
		{
			name: "unannotated artists, albums, and tracks with numeric sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=true",
					"-includeTracks=true",
					"-annotate=false",
					"-sort", "numeric",
				},
			},
			wantW: generateListing(true, true, true, false, true),
		},
		{
			name: "annotated artists, albums, and tracks with numeric sorting",
			l:    newLsSubCommand(internal.EmptyConfiguration(), flag.NewFlagSet("ls", flag.ContinueOnError)),
			args: args{
				[]string{
					"-topDir", topDir,
					"-includeAlbums=true",
					"-includeArtists=true",
					"-includeTracks=true",
					"-annotate=true",
					"-sort", "numeric",
				},
			},
			wantW: generateListing(true, true, true, true, true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.l.Exec(w, os.Stderr, tt.args.args)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("%s = %v, want %v", fnName, gotW, tt.wantW)
			}
		})
	}
}

func Test_newLsSubCommand(t *testing.T) {
	savedState := internal.SaveEnvVarForTesting("APPDATA")
	os.Setenv("APPDATA", internal.SecureAbsolutePathForTesting("."))
	defer func() {
		savedState.RestoreForTesting()
	}()
	topDir := "loadTest"
	fnName := "newLsSubCommand()"
	if err := internal.Mkdir(topDir); err != nil {
		t.Errorf("%s error creating %s: %v", fnName, topDir, err)
	}
	if err := internal.PopulateTopDirForTesting(topDir); err != nil {
		t.Errorf("%s error populating %s: %v", fnName, topDir, err)
	}
	if err := internal.CreateDefaultYamlFileForTesting(); err != nil {
		t.Errorf("error creating defaults.yaml: %v", err)
	}
	defer func() {
		internal.DestroyDirectoryForTesting(fnName, topDir)
		internal.DestroyDirectoryForTesting(fnName, "./mp3")
	}()
	defaultConfig, _ := internal.ReadConfigurationFile(os.Stderr)
	type args struct {
		c *internal.Configuration
	}
	tests := []struct {
		name                 string
		args                 args
		wantIncludeAlbums    bool
		wantIncludeArtists   bool
		wantIncludeTracks    bool
		wantTrackSorting     string
		wantAnnotateListings bool
	}{
		{
			name:                 "ordinary defaults",
			args:                 args{c: internal.EmptyConfiguration()},
			wantIncludeAlbums:    true,
			wantIncludeArtists:   true,
			wantIncludeTracks:    false,
			wantTrackSorting:     "numeric",
			wantAnnotateListings: false,
		},
		{
			name:                 "overridden defaults",
			args:                 args{c: defaultConfig},
			wantIncludeAlbums:    false,
			wantIncludeArtists:   false,
			wantIncludeTracks:    true,
			wantTrackSorting:     "alpha",
			wantAnnotateListings: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := newLsSubCommand(tt.args.c, flag.NewFlagSet("ls", flag.ContinueOnError))
			if _, ok := ls.sf.ProcessArgs(os.Stdout, []string{"-topDir", topDir, "-ext", ".mp3"}); ok {
				if *ls.includeAlbums != tt.wantIncludeAlbums {
					t.Errorf("%s %s: got includeAlbums %t want %t", fnName, tt.name, *ls.includeAlbums, tt.wantIncludeAlbums)
				}
				if *ls.includeArtists != tt.wantIncludeArtists {
					t.Errorf("%s %s: got includeArtists %t want %t", fnName, tt.name, *ls.includeArtists, tt.wantIncludeArtists)
				}
				if *ls.includeTracks != tt.wantIncludeTracks {
					t.Errorf("%s %s: got includeTracks %t want %t", fnName, tt.name, *ls.includeTracks, tt.wantIncludeTracks)
				}
				if *ls.annotateListings != tt.wantAnnotateListings {
					t.Errorf("%s %s: got annotateListings %t want %t", fnName, tt.name, *ls.annotateListings, tt.wantAnnotateListings)
				}
				if *ls.trackSorting != tt.wantTrackSorting {
					t.Errorf("%s %s: got trackSorting %q want %q", fnName, tt.name, *ls.trackSorting, tt.wantTrackSorting)
				}
			} else {
				t.Errorf("%s %s: error processing arguments", fnName, tt.name)
			}
		})
	}
}