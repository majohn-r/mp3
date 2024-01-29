/*
Copyright © 2021 Marc Johnson (marc.johnson27591@gmail.com)
*/
package cmd_test

import (
	"fmt"
	"mp3/cmd"
	"mp3/internal/files"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	cmd_toolkit "github.com/majohn-r/cmd-toolkit"
	"github.com/majohn-r/output"
	"github.com/spf13/cobra"
)

func TestProcessCheckFlags(t *testing.T) {
	tests := map[string]struct {
		values map[string]*cmd.FlagValue
		want   *cmd.CheckSettings
		want1  bool
		output.WantedRecording
	}{
		"no data": {
			values: map[string]*cmd.FlagValue{},
			want:   &cmd.CheckSettings{},
			want1:  false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"An internal error occurred: flag \"empty\" is not found.\n" +
					"An internal error occurred: flag \"files\" is not found.\n" +
					"An internal error occurred: flag \"numbering\" is not found.\n",
				Log: "" +
					"level='error' error='flag not found' flag='empty' msg='internal error'\n" +
					"level='error' error='flag not found' flag='files' msg='internal error'\n" +
					"level='error' error='flag not found' flag='numbering' msg='internal error'\n",
			},
		},
		"out of the box": {
			values: map[string]*cmd.FlagValue{
				"empty":     {Value: false},
				"files":     {Value: false},
				"numbering": {Value: false},
			},
			want:  &cmd.CheckSettings{},
			want1: true,
		},
		"overridden": {
			values: map[string]*cmd.FlagValue{
				"empty":     {Value: true, ExplicitlySet: true},
				"files":     {Value: true, ExplicitlySet: true},
				"numbering": {Value: true, ExplicitlySet: true},
			},
			want: &cmd.CheckSettings{
				Empty:            true,
				EmptyUserSet:     true,
				Files:            true,
				FilesUserSet:     true,
				Numbering:        true,
				NumberingUserSet: true,
			},
			want1: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			got, got1 := cmd.ProcessCheckFlags(o, tt.values)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessCheckFlags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ProcessCheckFlags() got1 = %v, want %v", got1, tt.want1)
			}
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("ProcessCheckFlags() %s", issue)
				}
			}
		})
	}
}

func TestCheckSettings_HasWorkToDo(t *testing.T) {
	tests := map[string]struct {
		cs   *cmd.CheckSettings
		want bool
		output.WantedRecording
	}{
		"no work, as configured": {
			cs:   &cmd.CheckSettings{},
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"The flags --empty, --files, and --numbering are all configured false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
			},
		},
		"no work, empty configured that way": {
			cs:   &cmd.CheckSettings{EmptyUserSet: true},
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"In addition to --files and --numbering configured false, you explicitly set --empty false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
			},
		},
		"no work, files configured that way": {
			cs:   &cmd.CheckSettings{FilesUserSet: true},
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"In addition to --empty and --numbering configured false, you explicitly set --files false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
			},
		},
		"no work, numbering configured that way": {
			cs:   &cmd.CheckSettings{NumberingUserSet: true},
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"In addition to --empty and --files configured false, you explicitly set --numbering false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
			},
		},
		"no work, empty and files configured that way": {
			cs:   &cmd.CheckSettings{EmptyUserSet: true, FilesUserSet: true},
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"In addition to --numbering configured false, you explicitly set --empty and --files false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
			},
		},
		"no work, empty and numbering configured that way": {
			cs:   &cmd.CheckSettings{EmptyUserSet: true, NumberingUserSet: true},
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"In addition to --files configured false, you explicitly set --empty and --numbering false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
			},
		},
		"no work, numbering and files configured that way": {
			cs:   &cmd.CheckSettings{NumberingUserSet: true, FilesUserSet: true},
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"In addition to --empty configured false, you explicitly set --files and --numbering false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
			},
		},
		"no work, all flags configured that way": {
			cs:   &cmd.CheckSettings{NumberingUserSet: true, FilesUserSet: true, EmptyUserSet: true},
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"You explicitly set --empty, --files, and --numbering false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
			},
		},
		"check empty":               {cs: &cmd.CheckSettings{Empty: true}, want: true},
		"check files":               {cs: &cmd.CheckSettings{Files: true}, want: true},
		"check numbering":           {cs: &cmd.CheckSettings{Numbering: true}, want: true},
		"check empty and files":     {cs: &cmd.CheckSettings{Empty: true, Files: true}, want: true},
		"check empty and numbering": {cs: &cmd.CheckSettings{Empty: true, Numbering: true}, want: true},
		"check numbering and files": {cs: &cmd.CheckSettings{Numbering: true, Files: true}, want: true},
		"check everything":          {cs: &cmd.CheckSettings{Empty: true, Files: true, Numbering: true}, want: true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			if got := tt.cs.HasWorkToDo(o); got != tt.want {
				t.Errorf("CheckSettings.HasWorkToDo() = %v, want %v", got, tt.want)
			}
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("CheckSettings.HasWorkToDo() %s", issue)
				}
			}
		})
	}
}

func TestIssueTypeAsString(t *testing.T) {
	tests := map[string]struct {
		i    cmd.CheckIssueType
		want string
	}{
		"unspecified": {i: cmd.CheckUnspecifiedIssue, want: "unspecified issue 0"},
		"empty":       {i: cmd.CheckEmptyIssue, want: "empty"},
		"files":       {i: cmd.CheckFilesIssue, want: "files"},
		"numbering":   {i: cmd.CheckNumberingIssue, want: "numbering"},
		"metadata":    {i: cmd.CheckConflictIssue, want: "metadata conflict"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cmd.IssueTypeAsString(tt.i); got != tt.want {
				t.Errorf("IssueTypeAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckedIssues_AddIssue(t *testing.T) {
	type args struct {
		source cmd.CheckIssueType
		issue  string
	}
	tests := map[string]struct {
		cI   cmd.CheckedIssues
		args args
	}{
		"add empty issue": {
			cI: cmd.NewCheckedIssues(),
			args: args{
				source: cmd.CheckEmptyIssue,
				issue:  "no albums"},
		},
		"add files issue": {
			cI: cmd.NewCheckedIssues(),
			args: args{
				source: cmd.CheckFilesIssue,
				issue:  "genre mismatch"},
		},
		"add numbering issue": {
			cI: cmd.NewCheckedIssues(),
			args: args{
				source: cmd.CheckNumberingIssue,
				issue:  "missing track 3"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.cI.HasIssues() {
				t.Errorf("CheckedIssues.AddIssue() has issues from the start")
			}
			tt.cI.AddIssue(tt.args.source, tt.args.issue)
			if !tt.cI.HasIssues() {
				t.Errorf("CheckedIssues.AddIssue() did not add an issue")
			}
		})
	}
}

func TestNewCheckedTrack(t *testing.T) {
	tests := map[string]struct {
		track          *files.Track
		wantValidValue bool
	}{
		"nil":  {track: nil, wantValidValue: false},
		"real": {track: sampleTrack, wantValidValue: true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := cmd.NewCheckedTrack(tt.track)
			if tt.wantValidValue {
				if got == nil {
					t.Errorf("NewCheckedTrack() = %v, want non-nil", got)
				} else {
					if got.HasIssues() {
						t.Errorf("NewCheckedTrack() has issues")
					}
					if got.Track() != tt.track {
						t.Errorf("NewCheckedTrack() has the wrong track")
					}
					got.AddIssue(cmd.CheckFilesIssue, "no metadata")
					if !got.HasIssues() {
						t.Errorf("NewCheckedTrack() does not reflect added issue")
					}
				}
			} else {
				if got != nil {
					t.Errorf("NewCheckedTrack() = %v, want nil", got)
				}
			}
		})
	}
}

func TestNewCheckedAlbum(t *testing.T) {
	var testAlbum *files.Album
	if albums := generateAlbums(1, 5); len(albums) > 0 {
		testAlbum = albums[0]
	}
	tests := map[string]struct {
		album          *files.Album
		wantValidAlbum bool
	}{
		"nil":  {album: nil, wantValidAlbum: false},
		"real": {album: testAlbum, wantValidAlbum: true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := cmd.NewCheckedAlbum(tt.album)
			if tt.wantValidAlbum {
				if got == nil {
					t.Errorf("NewCheckedAlbum() = %v, want non-nil", got)
				} else {
					if got.HasIssues() {
						t.Errorf("NewCheckedAlbum() created with issues")
					}
					if got.Album() != tt.album {
						t.Errorf("NewCheckedAlbum() created with wrong album: got %v, want %v", got.Album(), tt.album)
					}
					if len(got.Tracks()) != len(tt.album.Tracks()) {
						t.Errorf("NewCheckedAlbum() created with %d tracks, want %d", len(got.Tracks()), len(tt.album.Tracks()))
					}
					got.AddIssue(cmd.CheckNumberingIssue, "missing track 1")
					if !got.HasIssues() {
						t.Errorf("NewCheckedAlbum() cannot add issue")
					} else {
						got.CheckedIssues = cmd.NewCheckedIssues()
						if got.HasIssues() {
							t.Errorf("NewCheckedAlbum() has issues with clean map")
						}
						for _, track := range got.Tracks() {
							track.AddIssue(cmd.CheckFilesIssue, "missing metadata")
							break
						}
						if !got.HasIssues() {
							t.Errorf("NewCheckedAlbum() does not show issue assigned to track")
						}
					}
				}
			} else {
				if got != nil {
					t.Errorf("NewCheckedAlbum() = %v, want nil", got)
				}
			}
		})
	}
}

func TestNewCheckedArtist(t *testing.T) {
	tests := map[string]struct {
		artist          *files.Artist
		wantValidArtist bool
	}{
		"nil":  {artist: nil, wantValidArtist: false},
		"real": {artist: generateArtists(1, 4, 5)[0], wantValidArtist: true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := cmd.NewCheckedArtist(tt.artist)
			if tt.wantValidArtist {
				if got == nil {
					t.Errorf("NewCheckedArtist() = %v, want non-nil", got)
				} else {
					if got.HasIssues() {
						t.Errorf("NewCheckedArtist() created with issues")
					}
					if got.Artist() != tt.artist {
						t.Errorf("NewCheckedArtist() created with wrong artist: got %v, want %v", got.Artist(), tt.artist)
					}
					if len(got.Albums()) != len(tt.artist.Albums()) {
						t.Errorf("NewCheckedArtist() created with %d albums, want %d", len(got.Albums()), len(tt.artist.Albums()))
					}
					got.AddIssue(cmd.CheckEmptyIssue, "no albums!")
					if !got.HasIssues() {
						t.Errorf("NewCheckedArtist()) cannot add issue")
					} else {
						got.CheckedIssues = cmd.NewCheckedIssues()
						if got.HasIssues() {
							t.Errorf("NewCheckedArtist() has issues with clean map")
						}
						for _, track := range got.Albums() {
							track.AddIssue(cmd.CheckNumberingIssue, "missing track 909")
							break
						}
						if !got.HasIssues() {
							t.Errorf("NewCheckedArtist() does not show issue assigned to track")
						}
					}
				}
			} else {
				if got != nil {
					t.Errorf("NewCheckedArtist() = %v, want nil", got)
				}
			}
		})
	}
}

func TestCheckedIssues_OutputIssues(t *testing.T) {
	tests := map[string]struct {
		tab    int
		issues map[cmd.CheckIssueType][]string
		output.WantedRecording
	}{
		"no issues": {tab: 0, issues: nil, WantedRecording: output.WantedRecording{}},
		"lots of issues, untabbed": {
			tab: 0,
			issues: map[cmd.CheckIssueType][]string{
				cmd.CheckEmptyIssue:     {"no albums", "no tracks"},
				cmd.CheckFilesIssue:     {"track 1 no data", "track 0 no data"},
				cmd.CheckNumberingIssue: {"missing track 4", "missing track 1"},
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"* [empty] no albums\n" +
					"* [empty] no tracks\n" +
					"* [files] track 0 no data\n" +
					"* [files] track 1 no data\n" +
					"* [numbering] missing track 1\n" +
					"* [numbering] missing track 4\n",
			},
		},
		"lots of issues, indented": {
			tab: 2,
			issues: map[cmd.CheckIssueType][]string{
				cmd.CheckEmptyIssue:     {"no albums", "no tracks"},
				cmd.CheckFilesIssue:     {"track 1 no data", "track 0 no data"},
				cmd.CheckNumberingIssue: {"missing track 4", "missing track 1"},
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  * [empty] no albums\n" +
					"  * [empty] no tracks\n" +
					"  * [files] track 0 no data\n" +
					"  * [files] track 1 no data\n" +
					"  * [numbering] missing track 1\n" +
					"  * [numbering] missing track 4\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cI := cmd.NewCheckedIssues()
			o := output.NewRecorder()
			for k, v := range tt.issues {
				for _, s := range v {
					cI.AddIssue(k, s)
				}
			}
			cI.OutputIssues(o, tt.tab)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("CheckedIssue.OutputIssues() %s", issue)
				}
			}
		})
	}
}

func TestCheckedTrack_OutputIssues(t *testing.T) {
	tests := map[string]struct {
		cT     *cmd.CheckedTrack
		issues map[cmd.CheckIssueType][]string
		output.WantedRecording
	}{
		"no issues": {
			cT:              cmd.NewCheckedTrack(sampleTrack),
			issues:          nil,
			WantedRecording: output.WantedRecording{},
		},
		"some issues": {
			cT: cmd.NewCheckedTrack(sampleTrack),
			issues: map[cmd.CheckIssueType][]string{
				cmd.CheckFilesIssue: {"missing ID3V1 metadata", "missing ID3V2 metadata"},
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"    Track \"track 10\"\n" +
					"    * [files] missing ID3V1 metadata\n" +
					"    * [files] missing ID3V2 metadata\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			for k, v := range tt.issues {
				for _, s := range v {
					tt.cT.AddIssue(k, s)
				}
			}
			o := output.NewRecorder()
			tt.cT.OutputIssues(o)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("CheckedTrack.OutputIssues() %s", issue)
				}
			}
		})
	}
}

func TestCheckedAlbum_OutputIssues(t *testing.T) {
	var album1 *files.Album
	if albums := generateAlbums(1, 1); len(albums) > 0 {
		album1 = albums[0]
	}
	albumWithIssues := cmd.NewCheckedAlbum(album1)
	if albumWithIssues != nil {
		albumWithIssues.AddIssue(cmd.CheckNumberingIssue, "missing track 2")
	}
	var album2 *files.Album
	if albums := generateAlbums(1, 4); len(albums) > 0 {
		album2 = albums[0]
	}
	albumWithTrackIssues := cmd.NewCheckedAlbum(album2)
	if albumWithTrackIssues != nil {
		albumWithTrackIssues.Tracks()[3].AddIssue(cmd.CheckFilesIssue, "no metadata detected")
	}
	var nilAlbum *files.Album
	if albums := generateAlbums(1, 2); len(albums) > 0 {
		nilAlbum = albums[0]
	}
	tests := map[string]struct {
		cAl *cmd.CheckedAlbum
		output.WantedRecording
	}{
		"nil": {cAl: cmd.NewCheckedAlbum(nilAlbum)},
		"album with issues itself": {
			cAl: albumWithIssues,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  Album \"my album 00\"\n" +
					"  * [numbering] missing track 2\n",
			},
		},
		"album with track issues": {
			cAl: albumWithTrackIssues,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  Album \"my album 00\"\n" +
					"    Track \"my track 004\"\n" +
					"    * [files] no metadata detected\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.cAl.OutputIssues(o)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("CheckedAlbum.OutputIssues() %s", issue)
				}
			}
		})
	}
}

func TestCheckedArtist_OutputIssues(t *testing.T) {
	// artist without issues
	var artist1 *files.Artist
	if artists := generateArtists(1, 1, 1); len(artists) > 0 {
		artist1 = artists[0]
	}
	cAr000 := cmd.NewCheckedArtist(artist1)
	// artist with artist issues
	var artist2 *files.Artist
	if artists := generateArtists(1, 0, 0); len(artists) > 0 {
		artist2 = artists[0]
	}
	cAr001 := cmd.NewCheckedArtist(artist2)
	if cAr001 != nil {
		cAr001.AddIssue(cmd.CheckEmptyIssue, "no albums")
	}
	// artist with artist and album issues
	var artist3 *files.Artist
	if artists := generateArtists(1, 1, 0); len(artists) > 0 {
		artist3 = artists[0]
	}
	cAr011 := cmd.NewCheckedArtist(artist3)
	if cAr011 != nil {
		cAr011.AddIssue(cmd.CheckEmptyIssue, "expected no albums")
		cAr011.Albums()[0].AddIssue(cmd.CheckEmptyIssue, "no tracks")
	}
	// artist with artist, album, and track issues
	var artist4 *files.Artist
	if artists := generateArtists(1, 1, 1); len(artists) > 0 {
		artist4 = artists[0]
	}
	cAr111 := cmd.NewCheckedArtist(artist4)
	if cAr111 != nil {
		cAr111.AddIssue(cmd.CheckEmptyIssue, "expected no albums")
		cAr111.Albums()[0].AddIssue(cmd.CheckEmptyIssue, "expected no tracks")
		cAr111.Albums()[0].Tracks()[0].AddIssue(cmd.CheckFilesIssue, "no metadata")
	}
	// artist with artist and track issues
	var artist5 *files.Artist
	if artists := generateArtists(1, 1, 1); len(artists) > 0 {
		artist5 = artists[0]
	}
	cAr101 := cmd.NewCheckedArtist(artist5)
	if cAr101 != nil {
		cAr101.AddIssue(cmd.CheckEmptyIssue, "expected no albums")
		cAr101.Albums()[0].Tracks()[0].AddIssue(cmd.CheckFilesIssue, "no metadata")
	}
	// artist with album issues
	var artist6 *files.Artist
	if artists := generateArtists(1, 1, 1); len(artists) > 0 {
		artist6 = artists[0]
	}
	cAr010 := cmd.NewCheckedArtist(artist6)
	if cAr010 != nil {
		cAr010.Albums()[0].AddIssue(cmd.CheckEmptyIssue, "expected no tracks")
	}
	// artist with album and track issues
	var artist7 *files.Artist
	if artists := generateArtists(1, 1, 1); len(artists) > 0 {
		artist7 = artists[0]
	}
	cAr110 := cmd.NewCheckedArtist(artist7)
	if cAr110 != nil {
		cAr110.Albums()[0].AddIssue(cmd.CheckEmptyIssue, "expected no tracks")
		cAr110.Albums()[0].Tracks()[0].AddIssue(cmd.CheckFilesIssue, "no metadata")
	}
	// artist with track issues
	var artist8 *files.Artist
	if artists := generateArtists(1, 1, 1); len(artists) > 0 {
		artist8 = artists[0]
	}
	cAr100 := cmd.NewCheckedArtist(artist8)
	if cAr100 != nil {
		cAr100.Albums()[0].Tracks()[0].AddIssue(cmd.CheckFilesIssue, "no metadata")
	}
	tests := map[string]struct {
		cAr *cmd.CheckedArtist
		output.WantedRecording
	}{
		"nothing": {cAr: cAr000},
		"bad artist": {
			cAr: cAr001,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist \"my artist 0\"\n" +
					"* [empty] no albums\n",
			},
		},
		"bad artist, bad album": {
			cAr: cAr011,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist \"my artist 0\"\n" +
					"* [empty] expected no albums\n" +
					"  Album \"my album 00\"\n" +
					"  * [empty] no tracks\n",
			},
		},
		"bad artist, bad album, bad track": {
			cAr: cAr111,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist \"my artist 0\"\n" +
					"* [empty] expected no albums\n" +
					"  Album \"my album 00\"\n" +
					"  * [empty] expected no tracks\n" +
					"    Track \"my track 001\"\n" +
					"    * [files] no metadata\n",
			},
		},
		"bad artist, bad track": {
			cAr: cAr101,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist \"my artist 0\"\n" +
					"* [empty] expected no albums\n" +
					"  Album \"my album 00\"\n" +
					"    Track \"my track 001\"\n" +
					"    * [files] no metadata\n",
			},
		},
		"bad album": {
			cAr: cAr010,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist \"my artist 0\"\n" +
					"  Album \"my album 00\"\n" +
					"  * [empty] expected no tracks\n",
			},
		},
		"bad album, bad track": {
			cAr: cAr110,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist \"my artist 0\"\n" +
					"  Album \"my album 00\"\n" +
					"  * [empty] expected no tracks\n" +
					"    Track \"my track 001\"\n" +
					"    * [files] no metadata\n",
			},
		},
		"bad track": {
			cAr: cAr100,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist \"my artist 0\"\n" +
					"  Album \"my album 00\"\n" +
					"    Track \"my track 001\"\n" +
					"    * [files] no metadata\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.cAr.OutputIssues(o)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("CheckedArtist.OutputIssues() %s", issue)
				}
			}
		})
	}
}

func TestPrepareCheckedArtists(t *testing.T) {
	tests := map[string]struct {
		artists    []*files.Artist
		want       int
		wantAlbums int
		wantTracks int
	}{
		"empty": {},
		"plenty": {
			artists:    generateArtists(15, 16, 17),
			want:       15,
			wantAlbums: 15 * 16,
			wantTracks: 15 * 16 * 17,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cmd.PrepareCheckedArtists(tt.artists); len(got) != tt.want {
				t.Errorf("PrepareCheckedArtists() = %d, want %v", len(got), tt.want)
			} else {
				albums := 0
				tracks := 0
				collectedTracks := []*files.Track{}
				for _, artist := range got {
					albums += len(artist.Albums())
					for _, album := range artist.Albums() {
						tracks += len(album.Tracks())
						for _, cT := range album.Tracks() {
							collectedTracks = append(collectedTracks, cT.Track())
						}
					}
				}
				if albums != tt.wantAlbums {
					t.Errorf("PrepareCheckedArtists() = %d albums, want %v", albums, tt.wantAlbums)
				}
				if tracks != tt.wantTracks {
					t.Errorf("PrepareCheckedArtists() = %d tracks, want %v", tracks, tt.wantTracks)
				}
				for _, track := range collectedTracks {
					found := false
					for _, cAr := range got {
						if cAr.Lookup(track) != nil {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("PrepareCheckedArtists() cannot find track %q on %q by %q", track.FileName(), track.AlbumName(), track.RecordingArtist())
					}
				}
				copiedTracks := []*files.Track{}
				for _, artist := range tt.artists {
					copiedAr := artist.Copy()
					for _, album := range artist.Albums() {
						copiedAl := album.Copy(copiedAr, false)
						copiedAr.AddAlbum(copiedAl)
						for _, track := range album.Tracks() {
							copiedTr := track.Copy(copiedAl)
							copiedAl.AddTrack(copiedTr)
							copiedTracks = append(copiedTracks, copiedTr)
						}
					}
				}
				for _, track := range copiedTracks {
					found := false
					for _, cAr := range got {
						if cAr.Lookup(track) != nil {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("PrepareCheckedArtists() cannot find copied track %q on %q by %q", track.FileName(), track.AlbumName(), track.RecordingArtist())
					}
				}
			}
		})
	}
}

func TestCheckSettings_PerformEmptyAnalysis(t *testing.T) {
	tests := map[string]struct {
		cs             *cmd.CheckSettings
		checkedArtists []*cmd.CheckedArtist
		want           bool
	}{
		"do nothing": {cs: &cmd.CheckSettings{Empty: false}},
		"empty slice": {
			cs:             &cmd.CheckSettings{Empty: true},
			checkedArtists: nil,
		},
		"full slice, no issues": {
			cs:             &cmd.CheckSettings{Empty: true},
			checkedArtists: cmd.PrepareCheckedArtists(generateArtists(5, 6, 7)),
		},
		"empty artists": {
			cs:             &cmd.CheckSettings{Empty: true},
			checkedArtists: cmd.PrepareCheckedArtists(generateArtists(1, 0, 10)),
			want:           true,
		},
		"empty albums": {
			cs:             &cmd.CheckSettings{Empty: true},
			checkedArtists: cmd.PrepareCheckedArtists(generateArtists(4, 6, 0)),
			want:           true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.cs.PerformEmptyAnalysis(tt.checkedArtists); got != tt.want {
				t.Errorf("CheckSettings.PerformEmptyAnalysis() = %v, want %v", got, tt.want)
			}
			verifiedFound := false
			for _, artist := range tt.checkedArtists {
				if artist.HasIssues() {
					verifiedFound = true
				}
			}
			if verifiedFound != tt.want {
				t.Errorf("CheckSettings.PerformEmptyAnalysis() verified = %v, want %v", verifiedFound, tt.want)
			}
		})
	}
}

func TestGenerateMissingNumbers(t *testing.T) {
	type args struct {
		low  int
		high int
	}
	tests := map[string]struct {
		args
		want string
	}{
		"equal":   {args: args{low: 2, high: 2}, want: "2"},
		"inequal": {args: args{low: 2, high: 3}, want: "2-3"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cmd.GenerateMissingNumbers(tt.args.low, tt.args.high); got != tt.want {
				t.Errorf("GenerateMissingNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateNumberingIssues(t *testing.T) {
	type args struct {
		m        map[int][]string
		maxTrack int
	}
	tests := map[string]struct {
		args
		want []string
	}{
		"empty": {
			args: args{m: nil, maxTrack: 0},
			want: []string{},
		},
		"clean": {
			args: args{
				m: map[int][]string{
					1: {"track 1"},
					2: {"track 2"},
					3: {"track 3"},
					4: {"track 4"},
					5: {"track 5"},
				},
				maxTrack: 5,
			},
			want: []string{},
		},
		"problematic": {
			args: args{
				m: map[int][]string{
					3:  {"track 3"},
					5:  {"track 4", "track 5", "some other track"},
					8:  {"track 8"},
					9:  {},
					10: {"track 10"},
					19: {"track 19"},
				},
				maxTrack: 20,
			},
			want: []string{
				"multiple tracks identified as track 5: \"some other track\", \"track 4\" and \"track 5\"",
				"missing tracks identified: 1-2, 4, 6-7, 9, 11-18, 20",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cmd.GenerateNumberingIssues(tt.args.m, tt.args.maxTrack); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateNumberingIssues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckSettings_PerformNumberingAnalysis(t *testing.T) {
	defectiveArtists := []*files.Artist{}
	for r := 0; r < 4; r++ {
		artistName := fmt.Sprintf("my artist %d", r)
		artist := files.NewArtist(artistName, filepath.Join("Music", artistName))
		for k := 0; k < 5; k++ {
			albumName := fmt.Sprintf("my album %d%d", r, k)
			album := files.NewAlbum(albumName, artist, filepath.Join("Music", "my artist", albumName))
			for j := 1; j <= 6; j += 2 {
				trackName := fmt.Sprintf("my track %d%d%d", r, k, j)
				track := files.NewTrack(album, fmt.Sprintf("%d %s.mp3", j, trackName), trackName, j)
				album.AddTrack(track)
			}
			artist.AddAlbum(album)
		}
		defectiveArtists = append(defectiveArtists, artist)
	}

	tests := map[string]struct {
		cs             *cmd.CheckSettings
		checkedArtists []*cmd.CheckedArtist
		want           bool
	}{
		"no analysis": {
			cs:             &cmd.CheckSettings{Numbering: false},
			checkedArtists: cmd.PrepareCheckedArtists(generateArtists(5, 6, 7)),
			want:           false,
		},
		"ok analysis": {
			cs:             &cmd.CheckSettings{Numbering: true},
			checkedArtists: cmd.PrepareCheckedArtists(generateArtists(5, 6, 7)),
			want:           false,
		},
		"missing numbers found": {
			cs:             &cmd.CheckSettings{Numbering: true},
			checkedArtists: cmd.PrepareCheckedArtists(defectiveArtists),
			want:           true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.cs.PerformNumberingAnalysis(tt.checkedArtists); got != tt.want {
				t.Errorf("CheckSettings.PerformNumberingAnalysis() = %v, want %v", got, tt.want)
			}
			verifiedFound := false
			for _, artist := range tt.checkedArtists {
				if artist.HasIssues() {
					verifiedFound = true
				}
			}
			if verifiedFound != tt.want {
				t.Errorf("CheckSettings.PerformNumberingAnalysis() verified = %v, want %v", verifiedFound, tt.want)
			}
		})
	}
}

func TestRecordFileIssues(t *testing.T) {
	originalArtists := generateArtists(5, 6, 7)
	tracks := []*files.Track{}
	for _, artist := range originalArtists {
		copiedArtist := artist.Copy()
		for _, album := range artist.Albums() {
			copiedAlbum := album.Copy(copiedArtist, true)
			copiedArtist.AddAlbum(copiedAlbum)
			tracks = append(tracks, copiedAlbum.Tracks()...)
		}
	}
	type args struct {
		checkedArtists []*cmd.CheckedArtist
		track          *files.Track
		issues         []string
	}
	tests := map[string]struct {
		args
		wantFoundIssues bool
	}{
		"no issues": {
			args:            args{checkedArtists: nil, track: nil, issues: nil},
			wantFoundIssues: false,
		},
		"issues": {
			args: args{
				checkedArtists: cmd.PrepareCheckedArtists(originalArtists),
				track:          tracks[len(tracks)-1],
				issues:         []string{"mismatched artist", "mismatched album"},
			},
			wantFoundIssues: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if gotFoundIssues := cmd.RecordFileIssues(tt.args.checkedArtists, tt.args.track, tt.args.issues); gotFoundIssues != tt.wantFoundIssues {
				t.Errorf("RecordFileIssues() = %v, want %v", gotFoundIssues, tt.wantFoundIssues)
			}
			if tt.wantFoundIssues {
				hasIssues := false
				for _, cAr := range tt.args.checkedArtists {
					if cAr.HasIssues() {
						hasIssues = true
					}
				}
				if !hasIssues {
					t.Errorf("RecordFileIssues() true, but no issues actually recorded")
				}
			}
		})
	}
}

func TestCheckSettings_PerformFileAnalysis(t *testing.T) {
	originalReadMetadata := cmd.ReadMetadata
	defer func() {
		cmd.ReadMetadata = originalReadMetadata
	}()
	cmd.ReadMetadata = func(_ output.Bus, _ []*files.Artist) {}
	type args struct {
		checkedArtists []*cmd.CheckedArtist
		ss             *cmd.SearchSettings
	}
	tests := map[string]struct {
		cs *cmd.CheckSettings
		args
		want bool
		output.WantedRecording
	}{
		"not permitted to do anything": {
			cs:              &cmd.CheckSettings{Files: false},
			args:            args{},
			want:            false,
			WantedRecording: output.WantedRecording{},
		},
		"allowed, but nothing to check": {
			cs: &cmd.CheckSettings{Files: true},
			args: args{
				checkedArtists: []*cmd.CheckedArtist{},
				ss:             &cmd.SearchSettings{},
			},
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No music files remain after filtering.\n" +
					"Why?\n" +
					"After applying --artistFilter=<nil>, --albumFilter=<nil>, and --trackFilter=<nil>, no files remained.\n" +
					"What to do:\n" +
					"Use less restrictive filter settings.\n",
				Log: "level='error' --albumFilter='<nil>' --artistFilter='<nil>' --trackFilter='<nil>' msg='no files remain after filtering'\n",
			},
		},
		"work to do": {
			cs: &cmd.CheckSettings{Files: true},
			args: args{
				checkedArtists: cmd.PrepareCheckedArtists(generateArtists(4, 5, 6)),
				ss: &cmd.SearchSettings{
					ArtistFilter: regexp.MustCompile(".*"),
					AlbumFilter:  regexp.MustCompile(".*"),
					TrackFilter:  regexp.MustCompile(".*"),
				},
			},
			want:            true,
			WantedRecording: output.WantedRecording{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			if got := tt.cs.PerformFileAnalysis(o, tt.args.checkedArtists, tt.args.ss); got != tt.want {
				t.Errorf("CheckSettings.PerformFileAnalysis() = %v, want %v", got, tt.want)
			}
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("CheckSettings.PerformFileAnalysis() %s", issue)
				}
			}
		})
	}
}

func TestCheckSettings_MaybeReportCleanResults(t *testing.T) {
	type args struct {
		emptyIssues     bool
		numberingIssues bool
		fileIssues      bool
	}
	tests := map[string]struct {
		cs *cmd.CheckSettings
		args
		output.WantedRecording
	}{
		"no issues found because nothing was checked": {
			cs:              &cmd.CheckSettings{},
			args:            args{},
			WantedRecording: output.WantedRecording{},
		},
		"all issues found, everything was checked": {
			cs:              &cmd.CheckSettings{Empty: true, Numbering: true, Files: true},
			args:            args{emptyIssues: true, numberingIssues: true, fileIssues: true},
			WantedRecording: output.WantedRecording{},
		},
		"no issues found, everything was checked": {
			cs:   &cmd.CheckSettings{Empty: true, Numbering: true, Files: true},
			args: args{emptyIssues: false, numberingIssues: false, fileIssues: false},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Empty Folder Analysis: no empty folders found.\n" +
					"Numbering Analysis: no missing or duplicate tracks found.\n" +
					"File Analysis: no inconsistencies found.\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.cs.MaybeReportCleanResults(o, tt.args.emptyIssues, tt.args.numberingIssues, tt.args.fileIssues)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("CheckSettings.MaybeReportCleanResults() %s", issue)
				}
			}
		})
	}
}

func TestCheckSettings_PerformChecks(t *testing.T) {
	originalReadMetadata := cmd.ReadMetadata
	defer func() {
		cmd.ReadMetadata = originalReadMetadata
	}()
	cmd.ReadMetadata = func(_ output.Bus, _ []*files.Artist) {}
	type args struct {
		artists       []*files.Artist
		artistsLoaded bool
		ss            *cmd.SearchSettings
	}
	tests := map[string]struct {
		cs *cmd.CheckSettings
		args
		output.WantedRecording
	}{
		"no artists loaded": {
			cs:              nil,
			args:            args{artists: generateArtists(1, 1, 1), artistsLoaded: false, ss: nil},
			WantedRecording: output.WantedRecording{},
		},
		"no artists": {
			cs:              nil,
			args:            args{artists: nil, artistsLoaded: true, ss: nil},
			WantedRecording: output.WantedRecording{},
		},
		"artists to check, check everything": {
			cs: &cmd.CheckSettings{Empty: true, Numbering: true, Files: true},
			args: args{
				artists:       generateArtists(1, 2, 3),
				artistsLoaded: true,
				ss: &cmd.SearchSettings{
					ArtistFilter: regexp.MustCompile(".*"),
					AlbumFilter:  regexp.MustCompile(".*"),
					TrackFilter:  regexp.MustCompile(".*"),
				},
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist \"my artist 0\"\n" +
					"  Album \"my album 00\"\n" +
					"    Track \"my track 001\"\n" +
					"    * [files] differences cannot be determined: metadata has not been read\n" +
					"    Track \"my track 002\"\n" +
					"    * [files] differences cannot be determined: metadata has not been read\n" +
					"    Track \"my track 003\"\n" +
					"    * [files] differences cannot be determined: metadata has not been read\n" +
					"  Album \"my album 01\"\n" +
					"    Track \"my track 011\"\n" +
					"    * [files] differences cannot be determined: metadata has not been read\n" +
					"    Track \"my track 012\"\n" +
					"    * [files] differences cannot be determined: metadata has not been read\n" +
					"    Track \"my track 013\"\n" +
					"    * [files] differences cannot be determined: metadata has not been read\n" +
					"Empty Folder Analysis: no empty folders found.\n" +
					"Numbering Analysis: no missing or duplicate tracks found.\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.cs.PerformChecks(o, tt.args.artists, tt.args.artistsLoaded, tt.args.ss)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("CheckSettings.PerformChecks() %s", issue)
				}
			}
		})
	}
}

func TestCheckSettings_MaybeDoWork(t *testing.T) {
	tests := map[string]struct {
		cs *cmd.CheckSettings
		ss *cmd.SearchSettings
		output.WantedRecording
	}{
		"nothing to do": {
			cs: &cmd.CheckSettings{},
			ss: nil,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"The flags --empty, --files, and --numbering are all configured false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
			},
		},
		"try a little work": {
			cs: &cmd.CheckSettings{Empty: true},
			ss: &cmd.SearchSettings{
				TopDirectory:   filepath.Join(".", "no dir"),
				FileExtensions: []string{".mp3"},
				AlbumFilter:    regexp.MustCompile(".*"),
				ArtistFilter:   regexp.MustCompile(".*"),
				TrackFilter:    regexp.MustCompile(".*"),
			},
			WantedRecording: output.WantedRecording{
				Error: "" +
					"The directory \"no dir\" cannot be read: open no dir: The system cannot find the file specified.\n" +
					"No music files could be found using the specified parameters.\n" +
					"Why?\n" +
					"There were no directories found in \"no dir\" (the --topDir value).\n" +
					"What to do:\n" +
					"Set --topDir to the path of a directory that contains artist directories.\n",
				Log: "" +
					"level='error' directory='no dir' error='open no dir: The system cannot find the file specified.' msg='cannot read directory'\n" +
					"level='error' --topDir='no dir' msg='cannot find any artist directories'\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.cs.MaybeDoWork(o, tt.ss)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("CheckSettings.MaybeDoWork() %s", issue)
				}
			}
		})
	}
}

func TestCheckRun(t *testing.T) {
	cmd.InitGlobals()
	originalBus := cmd.Bus
	originalSearchFlags := cmd.SearchFlags
	defer func() {
		cmd.Bus = originalBus
		cmd.SearchFlags = originalSearchFlags
	}()
	cmd.SearchFlags = safeSearchFlags
	checkFlags := cmd.SectionFlags{
		SectionName: cmd.CheckCommand,
		Flags: map[string]*cmd.FlagDetails{
			cmd.CheckEmpty: {
				AbbreviatedName: cmd.CheckEmptyAbbr,
				Usage:           "report empty album and artist directories",
				ExpectedType:    cmd.BoolType,
				DefaultValue:    false,
			},
			cmd.CheckFiles: {
				AbbreviatedName: cmd.CheckFilesAbbr,
				Usage:           "report metadata/file inconsistencies",
				ExpectedType:    cmd.BoolType,
				DefaultValue:    false,
			},
			cmd.CheckNumbering: {
				AbbreviatedName: cmd.CheckNumberingAbbr,
				Usage:           "report missing track numbers and duplicated track numbering",
				ExpectedType:    cmd.BoolType,
				DefaultValue:    false,
			},
		},
	}
	command := &cobra.Command{}
	cmd.AddFlags(output.NewNilBus(), cmd_toolkit.EmptyConfiguration(), command.Flags(), checkFlags, true)
	type args struct {
		cmd *cobra.Command
		in1 []string
	}
	tests := map[string]struct {
		args
		output.WantedRecording
	}{
		"default case": {
			args: args{cmd: command},
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No checks will be executed.\n" +
					"Why?\n" +
					"The flags --empty, --files, and --numbering are all configured false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command line.\n",
				Log: "" +
					"level='info'" +
					" --albumFilter='.*'" +
					" --artistFilter='.*'" +
					" --empty='false'" +
					" --files='false'" +
					" --numbering='false'" +
					" --topDir='.'" +
					" --trackFilter='.*'" +
					" command='check'" +
					" empty-user-set='false'" +
					" files-user-set='false'" +
					" numbering-user-set='false'" +
					" msg='executing command'\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			cmd.Bus = o // cook getBus()
			cmd.CheckRun(tt.args.cmd, tt.args.in1)
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("ListRun() %s", issue)
				}
			}
		})
	}
}

func cloneCommand(original *cobra.Command) *cobra.Command {
	clone := &cobra.Command{
		Use:                   original.Use,
		DisableFlagsInUseLine: original.DisableFlagsInUseLine,
		Short:                 original.Short,
		Long:                  original.Long,
		Example:               original.Example,
		Run:                   original.Run,
	}
	return clone
}

func TestCheckHelp(t *testing.T) {
	originalSearchFlags := cmd.SearchFlags
	defer func() {
		cmd.SearchFlags = originalSearchFlags
	}()
	cmd.SearchFlags = safeSearchFlags
	commandUnderTest := cloneCommand(cmd.CheckCmd)
	cmd.AddFlags(output.NewNilBus(), cmd_toolkit.EmptyConfiguration(), commandUnderTest.Flags(), cmd.CheckFlags, true)
	tests := map[string]struct {
		output.WantedRecording
	}{
		"good": {
			WantedRecording: output.WantedRecording{
				Console: "" +
					"\"check\" runs checks on mp3 files and their containing directories and reports any problems detected\n" +
					"\n" +
					"Usage:\n" +
					"  check [--empty] [--files] [--numbering] [--albumFilter regex] [--artistFilter regex] [--trackFilter regex] [--topDir dir] [--extensions extensions]\n" +
					"\n" +
					"Examples:\n" +
					"check --empty\n" +
					"  reports empty artist and album directories\n" +
					"check --files\n" +
					"  reads each mp3 file's metadata and reports any inconsistencies found\n" +
					"check --numbering\n" +
					"  reports errors in the track numbers of mp3 files\n" +
					"\n" +
					"Flags:\n" +
					"      --albumFilter string    regular expression specifying which albums to select (default \".*\")\n" +
					"      --artistFilter string   regular expression specifying which artists to select (default \".*\")\n" +
					"  -e, --empty                 report empty album and artist directories (default false)\n" +
					"      --extensions string     comma-delimited list of file extensions used by mp3 files (default \".mp3\")\n" +
					"  -f, --files                 report metadata/file inconsistencies (default false)\n" +
					"  -n, --numbering             report missing track numbers and duplicated track numbering (default false)\n" +
					"      --topDir string         top directory specifying where to find mp3 files (default \".\")\n" +
					"      --trackFilter string    regular expression specifying which tracks to select (default \".*\")\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			command := commandUnderTest
			enableCommandRecording(o, command)
			command.Help()
			if issues, ok := o.Verify(tt.WantedRecording); !ok {
				for _, issue := range issues {
					t.Errorf("check Help() %s", issue)
				}
			}
		})
	}
}