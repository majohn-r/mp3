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
	"sort"
	"testing"

	cmd_toolkit "github.com/majohn-r/cmd-toolkit"
	"github.com/majohn-r/output"
	"github.com/spf13/cobra"
)

func TestProcessListFlags(t *testing.T) {
	tests := map[string]struct {
		values map[string]*cmd.FlagValue
		want   *cmd.ListSettings
		want1  bool
		output.WantedRecording
	}{
		"no data": {
			values: map[string]*cmd.FlagValue{},
			want:   cmd.NewListSettings(),
			WantedRecording: output.WantedRecording{
				Error: "An internal error occurred: flag \"albums\" is not found.\n" +
					"An internal error occurred: flag \"annotate\" is not found.\n" +
					"An internal error occurred: flag \"artists\" is not found.\n" +
					"An internal error occurred: flag \"details\" is not found.\n" +
					"An internal error occurred: flag \"diagnostic\" is not found.\n" +
					"An internal error occurred: flag \"byNumber\" is not found.\n" +
					"An internal error occurred: flag \"byTitle\" is not found.\n" +
					"An internal error occurred: flag \"tracks\" is not found.\n",
				Log: "level='error'" +
					" error='flag not found'" +
					" flag='albums'" +
					" msg='internal error'\n" +
					"level='error'" +
					" error='flag not found'" +
					" flag='annotate'" +
					" msg='internal error'\n" +
					"level='error'" +
					" error='flag not found'" +
					" flag='artists'" +
					" msg='internal error'\n" +
					"level='error'" +
					" error='flag not found'" +
					" flag='details'" +
					" msg='internal error'\n" +
					"level='error'" +
					" error='flag not found'" +
					" flag='diagnostic'" +
					" msg='internal error'\n" +
					"level='error'" +
					" error='flag not found'" +
					" flag='byNumber'" +
					" msg='internal error'\n" +
					"level='error'" +
					" error='flag not found'" +
					" flag='byTitle'" +
					" msg='internal error'\n" +
					"level='error'" +
					" error='flag not found'" +
					" flag='tracks'" +
					" msg='internal error'\n",
			},
		},
		"configured": {
			values: map[string]*cmd.FlagValue{
				"albums":     cmd.NewFlagValue().WithValue(true),
				"annotate":   cmd.NewFlagValue().WithValue(true),
				"artists":    cmd.NewFlagValue().WithValue(true),
				"details":    cmd.NewFlagValue().WithValue(true),
				"diagnostic": cmd.NewFlagValue().WithValue(true),
				"byNumber":   cmd.NewFlagValue().WithValue(true),
				"byTitle":    cmd.NewFlagValue().WithValue(true),
				"tracks":     cmd.NewFlagValue().WithValue(true),
			},
			want: cmd.NewListSettings().WithAlbums(true).WithAlbumsUserSet(
				false).WithAnnotate(true).WithArtists(true).WithArtistsUserSet(
				false).WithDetails(true).WithDiagnostic(true).WithSortByNumber(
				true).WithSortByNumberUserSet(false).WithSortByTitle(
				true).WithSortByTitleUserSet(false).WithTracks(
				true).WithTracksUserSet(false),
			want1: true,
		},
		"user set": {
			values: map[string]*cmd.FlagValue{
				"albums":     cmd.NewFlagValue().WithValue(false).WithExplicitlySet(true),
				"annotate":   cmd.NewFlagValue().WithValue(false).WithExplicitlySet(true),
				"artists":    cmd.NewFlagValue().WithValue(false).WithExplicitlySet(true),
				"details":    cmd.NewFlagValue().WithValue(false).WithExplicitlySet(true),
				"diagnostic": cmd.NewFlagValue().WithValue(false).WithExplicitlySet(true),
				"byNumber":   cmd.NewFlagValue().WithValue(false).WithExplicitlySet(true),
				"byTitle":    cmd.NewFlagValue().WithValue(false).WithExplicitlySet(true),
				"tracks":     cmd.NewFlagValue().WithValue(false).WithExplicitlySet(true),
			},
			want: cmd.NewListSettings().WithAlbums(false).WithAlbumsUserSet(
				true).WithAnnotate(false).WithArtists(false).WithArtistsUserSet(
				true).WithDetails(false).WithDiagnostic(false).WithSortByNumber(
				false).WithSortByNumberUserSet(true).WithSortByTitle(
				false).WithSortByTitleUserSet(true).WithTracks(false).WithTracksUserSet(
				true),
			want1: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			got, got1 := cmd.ProcessListFlags(o, tt.values)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessListFlags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ProcessListFlags() got1 = %v, want %v", got1, tt.want1)
			}
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ProcessListFlags() %s", difference)
				}
			}
		})
	}
}

func TestListSettingsHasWorkToDo(t *testing.T) {
	tests := map[string]struct {
		ls   *cmd.ListSettings
		want bool
		output.WantedRecording
	}{
		"none true, none explicitly set": {
			ls: cmd.NewListSettings(),
			WantedRecording: output.WantedRecording{
				Error: "No listing will be output.\n" +
					"Why?\n" +
					"The flags --albums, --artists, and --tracks are all configured" +
					" false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags" +
					" is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command" +
					" line.\n",
			},
		},
		"none true, tracks explicitly set": {
			ls: cmd.NewListSettings().WithTracksUserSet(true),
			WantedRecording: output.WantedRecording{
				Error: "No listing will be output.\n" +
					"Why?\n" +
					"In addition to --albums and --artists configured false, you" +
					" explicitly set --tracks false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags" +
					" is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command" +
					" line.\n",
			},
		},
		"none true, artists explicitly set": {
			ls: cmd.NewListSettings().WithArtistsUserSet(true),
			WantedRecording: output.WantedRecording{
				Error: "No listing will be output.\n" +
					"Why?\n" +
					"In addition to --albums and --tracks configured false, you" +
					" explicitly set --artists false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags" +
					" is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command" +
					" line.\n",
			},
		},
		"none true, artists and tracks explicitly set": {
			ls: cmd.NewListSettings().WithArtistsUserSet(true).WithTracksUserSet(true),
			WantedRecording: output.WantedRecording{
				Error: "No listing will be output.\n" +
					"Why?\n" +
					"In addition to --albums configured false, you explicitly set" +
					" --artists and --tracks false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags" +
					" is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command" +
					" line.\n",
			},
		},
		"none true, albums explicitly set": {
			ls: cmd.NewListSettings().WithAlbumsUserSet(true),
			WantedRecording: output.WantedRecording{
				Error: "No listing will be output.\n" +
					"Why?\n" +
					"In addition to --artists and --tracks configured false, you" +
					" explicitly set --albums false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags" +
					" is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command" +
					" line.\n",
			},
		},
		"none true, albums and tracks explicitly set": {
			ls: cmd.NewListSettings().WithAlbumsUserSet(true).WithTracksUserSet(true),
			WantedRecording: output.WantedRecording{
				Error: "No listing will be output.\n" +
					"Why?\n" +
					"In addition to --artists configured false, you explicitly set" +
					" --albums and --tracks false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags" +
					" is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command" +
					" line.\n",
			},
		},
		"none true, albums and artists explicitly set": {
			ls: cmd.NewListSettings().WithAlbumsUserSet(true).WithArtistsUserSet(true),
			WantedRecording: output.WantedRecording{
				Error: "No listing will be output.\n" +
					"Why?\n" +
					"In addition to --tracks configured false, you explicitly set" +
					" --albums and --artists false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags" +
					" is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command" +
					" line.\n",
			},
		},
		"none true, albums and artists and tracks explicitly set": {
			ls: cmd.NewListSettings().WithAlbumsUserSet(true).WithArtistsUserSet(
				true).WithTracksUserSet(true),
			WantedRecording: output.WantedRecording{
				Error: "No listing will be output.\n" +
					"Why?\n" +
					"You explicitly set --albums, --artists, and --tracks false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags" +
					" is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command" +
					" line.\n",
			},
		},
		"tracks true": {
			ls:   cmd.NewListSettings().WithTracks(true),
			want: true,
		},
		"artists true": {
			ls:   cmd.NewListSettings().WithArtists(true),
			want: true,
		},
		"artists and tracks true": {
			ls:   cmd.NewListSettings().WithArtists(true).WithTracks(true),
			want: true,
		},
		"albums true": {
			ls:   cmd.NewListSettings().WithAlbums(true),
			want: true,
		},
		"albums and tracks true": {
			ls:   cmd.NewListSettings().WithAlbums(true).WithTracks(true),
			want: true,
		},
		"albums and artists true": {
			ls:   cmd.NewListSettings().WithAlbums(true).WithArtists(true),
			want: true,
		},
		"albums and artists and tracks true": {
			ls:   cmd.NewListSettings().WithAlbums(true).WithArtists(true).WithTracks(true),
			want: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			if got := tt.ls.HasWorkToDo(o); got != tt.want {
				t.Errorf("ListSettings.HasWorkToDo() = %v, want %v", got, tt.want)
			}
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.HasWorkToDo() %s", difference)
				}
			}
		})
	}
}

func TestListSettingsTracksSortable(t *testing.T) {
	tests := map[string]struct {
		ls      *cmd.ListSettings
		want    bool
		lsFinal *cmd.ListSettings
		output.WantedRecording
	}{
		"tracks listed, both options set, neither explicitly": {
			ls: cmd.NewListSettings().WithTracks(true).WithSortByNumber(
				true).WithSortByTitle(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Track sorting cannot be done.\n" +
					"Why?\n" +
					"The --byNumber and --byTitle flags are both configured true.\n" +
					"What to do:\n" +
					"Either edit the configuration file and use those default values, or" +
					" use appropriate command line values.\n",
			},
		},
		"tracks listed, both options set, by number explicitly": {
			ls: cmd.NewListSettings().WithTracks(true).WithSortByNumber(
				true).WithSortByNumberUserSet(true).WithSortByTitle(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Track sorting cannot be done.\n" +
					"Why?\n" +
					"The --byTitle flag is configured true and you explicitly set" +
					" --byNumber true.\n" +
					"What to do:\n" +
					"Either edit the configuration file and use those default values, or" +
					" use appropriate command line values.\n",
			},
		},
		"tracks listed, both options set, by title explicitly": {
			ls: cmd.NewListSettings().WithTracks(true).WithSortByNumber(
				true).WithSortByTitle(true).WithSortByTitleUserSet(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Track sorting cannot be done.\n" +
					"Why?\n" +
					"The --byNumber flag is configured true and you explicitly set" +
					" --byTitle true.\n" +
					"What to do:\n" +
					"Either edit the configuration file and use those default values, or" +
					" use appropriate command line values.\n",
			},
		},
		"tracks listed, both options set, both explicitly": {
			ls: cmd.NewListSettings().WithTracks(true).WithSortByNumber(
				true).WithSortByNumberUserSet(true).WithSortByTitle(
				true).WithSortByTitleUserSet(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Track sorting cannot be done.\n" +
					"Why?\n" +
					"You explicitly set --byNumber and --byTitle true.\n" +
					"What to do:\n" +
					"Either edit the configuration file and use those default values, or" +
					" use appropriate command line values.\n",
			},
		},
		"tracks listed, no albums, sort by number, neither explicitly": {
			ls:   cmd.NewListSettings().WithTracks(true).WithSortByNumber(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Sorting tracks by number not possible.\n" +
					"Why?\n" +
					"Track numbers are only relevant if albums are also output.\n" +
					"--albums is configured as false, and --byNumber is configured as" +
					" true.\n" +
					"What to do:\n" +
					"Either edit the configuration file or change which flags you set on" +
					" the command line.\n",
			},
		},
		"tracks listed, no albums, sort by number, albums explicitly": {
			ls: cmd.NewListSettings().WithAlbumsUserSet(true).WithTracks(
				true).WithSortByNumber(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Sorting tracks by number not possible.\n" +
					"Why?\n" +
					"Track numbers are only relevant if albums are also output.\n" +
					"You set --albums false and --byNumber is configured as true.\n" +
					"What to do:\n" +
					"Either edit the configuration file or change which flags you set on" +
					" the command line.\n",
			},
		},
		"tracks listed, no albums, sort by number, sort explicitly": {
			ls: cmd.NewListSettings().WithTracks(true).WithSortByNumber(
				true).WithSortByNumberUserSet(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Sorting tracks by number not possible.\n" +
					"Why?\n" +
					"Track numbers are only relevant if albums are also output.\n" +
					"You set --byNumber true and --albums is configured as false.\n" +
					"What to do:\n" +
					"Either edit the configuration file or change which flags you set on" +
					" the command line.\n",
			},
		},
		"tracks listed, no albums, sort by number, both explicitly": {
			ls: cmd.NewListSettings().WithAlbumsUserSet(true).WithTracks(
				true).WithSortByNumber(true).WithSortByNumberUserSet(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Sorting tracks by number not possible.\n" +
					"Why?\n" +
					"Track numbers are only relevant if albums are also output.\n" +
					"You set --byNumber true and --albums false.\n" +
					"What to do:\n" +
					"Either edit the configuration file or change which flags you set on" +
					" the command line.\n",
			},
		},
		"tracks listed, both sorting options explicitly false": {
			ls: cmd.NewListSettings().WithTracks(true).WithSortByNumberUserSet(
				true).WithSortByTitleUserSet(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "A listing of tracks is not possible.\n" +
					"Why?\n" +
					"Tracks are enabled, but you set both --byNumber and --byTitle false.\n" +
					"What to do:\n" +
					"Enable one of the sorting flags.\n",
			},
		},
		"tracks listed, no sorting, user said no to number": {
			ls:   cmd.NewListSettings().WithTracks(true).WithSortByNumberUserSet(true),
			want: true,
			lsFinal: cmd.NewListSettings().WithTracks(true).WithSortByNumberUserSet(
				true).WithSortByTitle(true),
			WantedRecording: output.WantedRecording{
				Log: "level='info'" +
					" --albums='false'" +
					" --byTitle='true'" +
					" byNumber='false'" +
					" msg='no track sorting set, providing a sensible value'\n",
			},
		},
		"tracks listed, no sorting, user said no to title": {
			ls:   cmd.NewListSettings().WithTracks(true).WithSortByTitleUserSet(true),
			want: true,
			lsFinal: cmd.NewListSettings().WithTracks(true).WithSortByNumber(
				true).WithSortByTitleUserSet(true),
			WantedRecording: output.WantedRecording{
				Log: "level='info'" +
					" --albums='false'" +
					" --byTitle='false'" +
					" byNumber='true'" +
					" msg='no track sorting set, providing a sensible value'\n",
			},
		},
		"tracks listed, no sorting, albums included": {
			ls:   cmd.NewListSettings().WithAlbums(true).WithTracks(true),
			want: true,
			lsFinal: cmd.NewListSettings().WithAlbums(true).WithTracks(
				true).WithSortByNumber(true),
			WantedRecording: output.WantedRecording{
				Log: "level='info'" +
					" --albums='true'" +
					" --byTitle='false'" +
					" byNumber='true'" +
					" msg='no track sorting set, providing a sensible value'\n",
			},
		},
		"tracks listed, no sorting, no albums": {
			ls:      cmd.NewListSettings().WithTracks(true),
			want:    true,
			lsFinal: cmd.NewListSettings().WithTracks(true).WithSortByTitle(true),
			WantedRecording: output.WantedRecording{
				Log: "level='info'" +
					" --albums='false'" +
					" --byTitle='true'" +
					" byNumber='false'" +
					" msg='no track sorting set, providing a sensible value'\n",
			},
		},
		"tracks not listed, no sorting explicitly called for": {
			ls: cmd.NewListSettings().WithSortByNumberUserSet(
				true).WithSortByTitleUserSet(true),
			want: true,
			lsFinal: cmd.NewListSettings().WithSortByNumberUserSet(
				true).WithSortByTitleUserSet(true),
		},
		"tracks not listed, sort by number explicitly called for": {
			ls:   cmd.NewListSettings().WithSortByNumber(true).WithSortByNumberUserSet(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Your sorting preferences are not relevant.\n" +
					"Why?\n" +
					"Tracks are not included in the output, but you explicitly set" +
					" --byNumber or --byTitle true.\n" +
					"What to do:\n" +
					"Either set --tracks true or remove the sorting flags from the" +
					" command line.\n",
			},
		},
		"tracks not listed, sort by title explicitly called for": {
			ls:   cmd.NewListSettings().WithSortByTitle(true).WithSortByTitleUserSet(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Your sorting preferences are not relevant.\n" +
					"Why?\n" +
					"Tracks are not included in the output, but you explicitly set" +
					" --byNumber or --byTitle true.\nWhat to do:\n" +
					"Either set --tracks true or remove the sorting flags from the" +
					" command line.\n",
			},
		},
		"tracks not listed, sort by number and title explicitly called for": {
			ls: cmd.NewListSettings().WithSortByNumber(true).WithSortByNumberUserSet(
				true).WithSortByTitle(true).WithSortByTitleUserSet(true),
			want: false,
			WantedRecording: output.WantedRecording{
				Error: "Your sorting preferences are not relevant.\n" +
					"Why?\n" +
					"Tracks are not included in the output, but you explicitly set" +
					" --byNumber or --byTitle true.\n" +
					"What to do:\n" +
					"Either set --tracks true or remove the sorting flags from the" +
					" command line.\n",
			},
		},
		"tracks listed, albums too, just sort by number": {
			ls: cmd.NewListSettings().WithAlbums(true).WithTracks(
				true).WithSortByNumber(true),
			want: true,
			lsFinal: cmd.NewListSettings().WithAlbums(true).WithTracks(
				true).WithSortByNumber(true),
		},
		"tracks listed, just sort by title": {
			ls:      cmd.NewListSettings().WithTracks(true).WithSortByTitle(true),
			want:    true,
			lsFinal: cmd.NewListSettings().WithTracks(true).WithSortByTitle(true),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			if got := tt.ls.TracksSortable(o); got != tt.want {
				t.Errorf("ListSettings.TracksSortable() = %v, want %v", got, tt.want)
			}
			if tt.want {
				if *tt.ls != *tt.lsFinal {
					t.Errorf("ListSettings.TracksSortable() ls = %v, want %v", tt.ls,
						tt.lsFinal)
				}
			}
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.TracksSortable() %s", difference)
				}
			}
		})
	}
}

var (
	sampleTrack = files.NewTrack(
		files.NewAlbum(
			"my album",
			files.NewArtist("my artist", "music/my artist"),
			"music/my artist/my album"),
		"10 track 10.mp3", "track 10", 10)
	safeSearchFlags = cmd.NewSectionFlags().WithSectionName("search").WithFlags(
		map[string]*cmd.FlagDetails{
			cmd.SearchAlbumFilter: cmd.NewFlagDetails().WithUsage(
				"regular expression specifying which albums to select").WithExpectedType(
				cmd.StringType).WithDefaultValue(".*"),
			cmd.SearchArtistFilter: cmd.NewFlagDetails().WithUsage(
				"regular expression specifying which artists to select").WithExpectedType(
				cmd.StringType).WithDefaultValue(".*"),
			cmd.SearchTrackFilter: cmd.NewFlagDetails().WithUsage(
				"regular expression specifying which tracks to select").WithExpectedType(
				cmd.StringType).WithDefaultValue(".*"),
			cmd.SearchTopDir: cmd.NewFlagDetails().WithUsage(
				"top directory specifying where to find mp3 files").WithExpectedType(
				cmd.StringType).WithDefaultValue("."),
			cmd.SearchFileExtensions: cmd.NewFlagDetails().WithUsage(
				"comma-delimited list of file extensions used by mp3" +
					" files").WithExpectedType(cmd.StringType).WithDefaultValue(".mp3"),
		},
	)
)

func TestShowID3V1Diagnostics(t *testing.T) {
	type args struct {
		track *files.Track
		tags  []string
		err   error
		tab   int
	}
	tests := map[string]struct {
		args
		output.WantedRecording
	}{
		"with error": {
			args: args{
				track: sampleTrack,
				err:   fmt.Errorf("could not read track"),
				tab:   2,
			},
			WantedRecording: output.WantedRecording{
				Log: "level='error'" +
					" error='could not read track'" +
					" metadata='ID3V1'" +
					" track='music\\my artist\\my album\\10 track 10.mp3'" +
					" msg='metadata read error'\n",
			},
		},
		"without error": {
			args: args{
				track: sampleTrack,
				tags: []string{
					"artist=my artist",
					"album=my album",
					"track=track 10",
					"number=10",
				},
				tab: 2,
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  ID3V1 artist=my artist\n" +
					"  ID3V1 album=my album\n" +
					"  ID3V1 track=track 10\n" +
					"  ID3V1 number=10\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			cmd.ShowID3V1Diagnostics(o, tt.args.track, tt.args.tags, tt.args.err,
				tt.args.tab)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ShowID3V1Diagnostics() %s", difference)
				}
			}
		})
	}
}

func TestShowID3V2Diagnostics(t *testing.T) {
	type args struct {
		track    *files.Track
		version  byte
		encoding string
		frames   []string
		err      error
		tab      int
	}
	tests := map[string]struct {
		args
		output.WantedRecording
	}{
		"error": {
			args: args{
				track: sampleTrack,
				err:   fmt.Errorf("no ID3V2 data found"),
			},
			WantedRecording: output.WantedRecording{
				Log: "level='error'" +
					" error='no ID3V2 data found'" +
					" metadata='ID3V2'" +
					" track='music\\my artist\\my album\\10 track 10.mp3'" +
					" msg='metadata read error'\n",
			},
		},
		"empty frames": {
			args: args{
				track:    sampleTrack,
				version:  1,
				encoding: "UTF-8",
				tab:      2,
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  ID3V2 Version: 1\n" +
					"  ID3V2 Encoding: \"UTF-8\"\n",
			},
		},
		"with frames": {
			args: args{
				track:    sampleTrack,
				version:  1,
				encoding: "UTF-8",
				frames:   []string{"FRAME1", "FRAME2"},
				tab:      2,
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  ID3V2 Version: 1\n" +
					"  ID3V2 Encoding: \"UTF-8\"\n" +
					"  ID3V2 FRAME1\n" +
					"  ID3V2 FRAME2\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			cmd.ShowID3V2Diagnostics(o, tt.args.track, tt.args.version, tt.args.encoding,
				tt.args.frames, tt.args.err, tt.args.tab)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ShowID3V2Diagnostics() %s", difference)
				}
			}
		})
	}
}

func TestListSettingsListTrackDiagnostics(t *testing.T) {
	type args struct {
		track *files.Track
		tab   int
	}
	tests := map[string]struct {
		ls *cmd.ListSettings
		args
		output.WantedRecording
	}{
		"permitted": {
			ls:   cmd.NewListSettings().WithDiagnostic(true),
			args: args{track: sampleTrack, tab: 2},
			WantedRecording: output.WantedRecording{
				Log: "level='error'" +
					" error='open music\\my artist\\my album\\10 track 10.mp3: The system" +
					" cannot find the path specified.'" +
					" metadata='ID3V2'" +
					" track='music\\my artist\\my album\\10 track 10.mp3'" +
					" msg='metadata read error'\n" +
					"level='error'" +
					" error='open music\\my artist\\my album\\10 track 10.mp3: The system" +
					" cannot find the path specified.'" +
					" metadata='ID3V1'" +
					" track='music\\my artist\\my album\\10 track 10.mp3'" +
					" msg='metadata read error'\n",
			},
		},
		"not permitted": {
			ls: cmd.NewListSettings().WithDiagnostic(false),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.ls.ListTrackDiagnostics(o, tt.args.track, tt.args.tab)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.ListTrackDiagnostics() %s", difference)
				}
			}
		})
	}
}

func TestShowDetails(t *testing.T) {
	type args struct {
		track        *files.Track
		details      map[string]string
		detailsError error
		tab          int
	}
	tests := map[string]struct {
		args
		output.WantedRecording
	}{
		"error": {
			args: args{
				track:        sampleTrack,
				detailsError: fmt.Errorf("details service offline"),
			},
			WantedRecording: output.WantedRecording{
				Error: "The details are not available for track \"track 10\" on album" +
					" \"my album\" by artist \"my artist\": \"details service offline\".\n",
				Log: "level='error'" +
					" error='details service offline'" +
					" track='music\\my artist\\my album\\10 track 10.mp3'" +
					" msg='cannot get details'\n",
			},
		},
		"no error, and no details": {args: args{track: sampleTrack}},
		"no error, with details": {
			args: args{
				track: sampleTrack,
				details: map[string]string{
					"composer": "some German",
					"producer": "A True Genius",
				},
				tab: 2,
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  Details:\n" +
					"    composer = \"some German\"\n" +
					"    producer = \"A True Genius\"\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			cmd.ShowDetails(o, tt.args.track, tt.args.details, tt.args.detailsError,
				tt.args.tab)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ShowDetails() %s", difference)
				}
			}
		})
	}
}

func TestListSettingsListTrackDetails(t *testing.T) {
	type args struct {
		track *files.Track
		tab   int
	}
	tests := map[string]struct {
		ls *cmd.ListSettings
		args
		output.WantedRecording
	}{
		"not wanted": {ls: cmd.NewListSettings().WithDetails(false)},
		"wanted": {
			ls:   cmd.NewListSettings().WithDetails(true),
			args: args{track: sampleTrack, tab: 2},
			WantedRecording: output.WantedRecording{
				Error: "The details are not available for track \"track 10\" on album" +
					" \"my album\" by artist \"my artist\":" +
					" \"open music\\\\my artist\\\\my album\\\\10 track 10.mp3: The" +
					" system cannot find the path specified.\".\n",
				Log: "level='error'" +
					" error='open music\\my artist\\my album\\10 track 10.mp3: The system" +
					" cannot find the path specified.'" +
					" track='music\\my artist\\my album\\10 track 10.mp3'" +
					" msg='cannot get details'\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.ls.ListTrackDetails(o, tt.args.track, tt.args.tab)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.ListTrackDetails() %s", difference)
				}
			}
		})
	}
}

func TestListSettingsAnnotateTrackName(t *testing.T) {
	tests := map[string]struct {
		ls    *cmd.ListSettings
		track *files.Track
		want  string
	}{
		"no annotations": {
			ls:    cmd.NewListSettings().WithAnnotate(false),
			track: sampleTrack,
			want:  "track 10",
		},
		"annotations, albums printed": {
			ls:    cmd.NewListSettings().WithAnnotate(true).WithAlbums(true),
			track: sampleTrack,
			want:  "track 10",
		},
		"annotations, no albums, artists included": {
			ls: cmd.NewListSettings().WithAnnotate(true).WithAlbums(
				false).WithArtists(true),
			track: sampleTrack,
			want:  "\"track 10\" on \"my album\"",
		},
		"annotations, no albums, no artists": {
			ls: cmd.NewListSettings().WithAnnotate(true).WithAlbums(
				false).WithArtists(false),
			track: sampleTrack,
			want:  "\"track 10\" on \"my album\" by \"my artist\"",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.ls.AnnotateTrackName(tt.track); got != tt.want {
				t.Errorf("ListSettings.AnnotateTrackName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateTracks(count int) []*files.Track {
	albums := generateAlbums(1, count)
	for _, album := range albums {
		return album.Tracks()
	}
	return nil
}

func TestListSettingsListTracksByName(t *testing.T) {
	type args struct {
		tracks []*files.Track
		tab    int
	}
	tests := map[string]struct {
		ls *cmd.ListSettings
		args
		output.WantedRecording
	}{
		"no tracks": {
			ls:   cmd.NewListSettings(),
			args: args{tracks: nil, tab: 2},
		},
		"multiple tracks": {
			ls:   cmd.NewListSettings().WithAnnotate(true),
			args: args{tracks: generateTracks(25), tab: 0},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"\"my track 001\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0010\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0011\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0012\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0013\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0014\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0015\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0016\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0017\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0018\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0019\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 002\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0020\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0021\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0022\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0023\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0024\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 0025\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 003\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 004\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 005\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 006\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 007\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 008\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 009\" on \"my album 00\" by \"my artist 0\"\n",
			},
		},
		"https://github.com/majohn-r/mp3/issues/147": {
			ls: cmd.NewListSettings().WithAnnotate(true),
			args: args{
				tracks: []*files.Track{
					files.NewEmptyTrack().WithName("Old Brown Shoe").WithAlbum(
						files.NewEmptyAlbum().WithTitle("Anthology 3 [Disc 2]").WithArtist(
							files.NewEmptyArtist().WithFileName("The Beatles"))),
					files.NewEmptyTrack().WithName("Old Brown Shoe").WithAlbum(
						files.NewEmptyAlbum().WithTitle("Live In Japan [Disc 1]").WithArtist(
							files.NewEmptyArtist().WithFileName("George Harrison & Eric Clapton"))),
					files.NewEmptyTrack().WithName("Old Brown Shoe").WithAlbum(
						files.NewEmptyAlbum().WithTitle("Past Masters, Vol. 2").WithArtist(
							files.NewEmptyArtist().WithFileName("The Beatles"))),
					files.NewEmptyTrack().WithName("Old Brown Shoe").WithAlbum(
						files.NewEmptyAlbum().WithTitle("Songs From The Material World - A Tribute To George Harrison").WithArtist(
							files.NewEmptyArtist().WithFileName("Various Artists"))),
					files.NewEmptyTrack().WithName("Old Brown Shoe (Take 2)").WithAlbum(
						files.NewEmptyAlbum().WithTitle("Abbey Road- Sessions [Disc 2]").WithArtist(
							files.NewEmptyArtist().WithFileName("The Beatles"))),
				},
				tab: 0,
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					`"Old Brown Shoe" on "Anthology 3 [Disc 2]" by "The Beatles"
"Old Brown Shoe" on "Live In Japan [Disc 1]" by "George Harrison & Eric Clapton"
"Old Brown Shoe" on "Past Masters, Vol. 2" by "The Beatles"
"Old Brown Shoe" on "Songs From The Material World - A Tribute To George Harrison" by "Various Artists"
"Old Brown Shoe (Take 2)" on "Abbey Road- Sessions [Disc 2]" by "The Beatles"
`,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.ls.ListTracksByName(o, tt.args.tracks, tt.args.tab)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.ListTracksByName() %s", difference)
				}
			}
		})
	}
}

func TestListSettingsListTracksByNumber(t *testing.T) {
	type args struct {
		tracks []*files.Track
		tab    int
	}
	tests := map[string]struct {
		ls *cmd.ListSettings
		args
		output.WantedRecording
	}{
		"no tracks": {
			ls:   cmd.NewListSettings(),
			args: args{},
		},
		"lots of tracks": {
			ls:   cmd.NewListSettings(),
			args: args{tracks: generateTracks(17), tab: 2},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"   1. my track 001\n" +
					"   2. my track 002\n" +
					"   3. my track 003\n" +
					"   4. my track 004\n" +
					"   5. my track 005\n" +
					"   6. my track 006\n" +
					"   7. my track 007\n" +
					"   8. my track 008\n" +
					"   9. my track 009\n" +
					"  10. my track 0010\n" +
					"  11. my track 0011\n" +
					"  12. my track 0012\n" +
					"  13. my track 0013\n" +
					"  14. my track 0014\n" +
					"  15. my track 0015\n" +
					"  16. my track 0016\n" +
					"  17. my track 0017\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.ls.ListTracksByNumber(o, tt.args.tracks, tt.args.tab)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.ListTracksByNumber() %s", difference)
				}
			}
		})
	}
}

func TestListSettingsListTracks(t *testing.T) {
	type args struct {
		tracks []*files.Track
		tab    int
	}
	tests := map[string]struct {
		ls *cmd.ListSettings
		args
		output.WantedRecording
	}{
		"no tracks": {
			ls:   cmd.NewListSettings().WithTracks(true),
			args: args{},
		},
		"do not list tracks": {
			ls:   cmd.NewListSettings().WithTracks(false).WithSortByNumber(true),
			args: args{tracks: generateTracks(99)},
		},
		"list tracks by number": {
			ls:   cmd.NewListSettings().WithTracks(true).WithSortByNumber(true),
			args: args{tracks: generateTracks(25), tab: 2},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"   1. my track 001\n" +
					"   2. my track 002\n" +
					"   3. my track 003\n" +
					"   4. my track 004\n" +
					"   5. my track 005\n" +
					"   6. my track 006\n" +
					"   7. my track 007\n" +
					"   8. my track 008\n" +
					"   9. my track 009\n" +
					"  10. my track 0010\n" +
					"  11. my track 0011\n" +
					"  12. my track 0012\n" +
					"  13. my track 0013\n" +
					"  14. my track 0014\n" +
					"  15. my track 0015\n" +
					"  16. my track 0016\n" +
					"  17. my track 0017\n" +
					"  18. my track 0018\n" +
					"  19. my track 0019\n" +
					"  20. my track 0020\n" +
					"  21. my track 0021\n" +
					"  22. my track 0022\n" +
					"  23. my track 0023\n" +
					"  24. my track 0024\n" +
					"  25. my track 0025\n",
			},
		},
		"list tracks by name": {
			ls:   cmd.NewListSettings().WithTracks(true).WithSortByTitle(true),
			args: args{tracks: generateTracks(25), tab: 2},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  my track 001\n" +
					"  my track 0010\n" +
					"  my track 0011\n" +
					"  my track 0012\n" +
					"  my track 0013\n" +
					"  my track 0014\n" +
					"  my track 0015\n" +
					"  my track 0016\n" +
					"  my track 0017\n" +
					"  my track 0018\n" +
					"  my track 0019\n" +
					"  my track 002\n" +
					"  my track 0020\n" +
					"  my track 0021\n" +
					"  my track 0022\n" +
					"  my track 0023\n" +
					"  my track 0024\n" +
					"  my track 0025\n" +
					"  my track 003\n" +
					"  my track 004\n" +
					"  my track 005\n" +
					"  my track 006\n" +
					"  my track 007\n" +
					"  my track 008\n" +
					"  my track 009\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.ls.ListTracks(o, tt.args.tracks, tt.args.tab)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.ListTracks() %s", difference)
				}
			}
		})
	}
}

func TestListSettingsAnnotateAlbumName(t *testing.T) {
	tests := map[string]struct {
		ls   *cmd.ListSettings
		want string
	}{
		"no annotation, no artist": {
			ls:   cmd.NewListSettings().WithAnnotate(false).WithArtists(false),
			want: "my album",
		},
		"no annotation, with artist": {
			ls:   cmd.NewListSettings().WithAnnotate(false).WithArtists(true),
			want: "my album",
		},
		"with annotation, no artist": {
			ls:   cmd.NewListSettings().WithAnnotate(true).WithArtists(false),
			want: "\"my album\" by \"my artist\"",
		},
		"with annotation, with artist": {
			ls:   cmd.NewListSettings().WithAnnotate(true).WithArtists(true),
			want: "my album",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			album := files.NewAlbum("my album",
				files.NewArtist("my artist", filepath.Join("Music", "my artist")),
				filepath.Join("Music", "my artist", "my album"))
			if got := tt.ls.AnnotateAlbumName(album); got != tt.want {
				t.Errorf("ListSettings.AnnotateAlbumName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateArtists(artistCount, albumCount, trackCount int) []*files.Artist {
	artists := []*files.Artist{}
	for r := 0; r < artistCount; r++ {
		artistName := fmt.Sprintf("my artist %d", r)
		artist := files.NewArtist(artistName, filepath.Join("Music", artistName))
		for k := 0; k < albumCount; k++ {
			albumName := fmt.Sprintf("my album %d%d", r, k)
			album := files.NewAlbum(albumName, artist,
				filepath.Join("Music", "my artist", albumName))
			for j := 1; j <= trackCount; j++ {
				trackName := fmt.Sprintf("my track %d%d%d", r, k, j)
				track := files.NewTrack(album, fmt.Sprintf("%d %s.mp3", j, trackName),
					trackName, j)
				album.AddTrack(track)
			}
			artist.AddAlbum(album)
		}
		artists = append(artists, artist)
	}
	return artists
}

func generateAlbums(albumCount, trackCount int) []*files.Album {
	artists := generateArtists(1, albumCount, trackCount)
	for _, artist := range artists {
		return artist.Albums()
	}
	return nil
}

func TestListSettingsListAlbums(t *testing.T) {
	type args struct {
		albums []*files.Album
		tab    int
	}
	tests := map[string]struct {
		ls *cmd.ListSettings
		args
		output.WantedRecording
	}{
		"no albums": {
			ls:   cmd.NewListSettings(),
			args: args{albums: nil, tab: 0},
		},
		"list albums without tracks": {
			ls: cmd.NewListSettings().WithAlbums(true),
			args: args{
				albums: generateAlbums(3, 3),
				tab:    2,
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  Album: my album 00\n" +
					"  Album: my album 01\n" +
					"  Album: my album 02\n",
			},
		},
		"list tracks only": {
			ls: cmd.NewListSettings().WithArtists(true).WithTracks(true).WithAnnotate(
				true).WithSortByTitle(true),
			args: args{
				albums: generateAlbums(2, 2),
				tab:    2,
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"  \"my track 001\" on \"my album 00\"\n" +
					"  \"my track 002\" on \"my album 00\"\n" +
					"  \"my track 011\" on \"my album 01\"\n" +
					"  \"my track 012\" on \"my album 01\"\n",
			},
		},
		"list albums and tracks": {
			ls: cmd.NewListSettings().WithAlbums(true).WithTracks(true).WithAnnotate(
				true).WithSortByNumber(true),
			args: args{
				albums: generateAlbums(3, 3),
				tab:    0,
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Album: \"my album 00\" by \"my artist 0\"\n" +
					"   1. my track 001\n" +
					"   2. my track 002\n" +
					"   3. my track 003\n" +
					"Album: \"my album 01\" by \"my artist 0\"\n" +
					"   1. my track 011\n" +
					"   2. my track 012\n" +
					"   3. my track 013\n" +
					"Album: \"my album 02\" by \"my artist 0\"\n" +
					"   1. my track 021\n" +
					"   2. my track 022\n" +
					"   3. my track 023\n",
			},
		},
		"https://github.com/majohn-r/mp3/issues/147": {
			ls: cmd.NewListSettings().WithAlbums(true).WithAnnotate(true),
			args: args{
				albums: []*files.Album{
					files.NewEmptyAlbum().WithTitle("Live Rhymin' [Bonus Tracks]").WithArtist(
						files.NewEmptyArtist().WithFileName("Paul Simon")),
					files.NewEmptyAlbum().WithTitle("Live In Paris & Toronto [Disc 2]").WithArtist(
						files.NewEmptyArtist().WithFileName("Loreena McKennitt")),
					files.NewEmptyAlbum().WithTitle("Live In Paris & Toronto [Disc 1]").WithArtist(
						files.NewEmptyArtist().WithFileName("Loreena McKennitt")),
					files.NewEmptyAlbum().WithTitle("Live In Japan [Disc 2]").WithArtist(
						files.NewEmptyArtist().WithFileName("George Harrison & Eric Clapton")),
					files.NewEmptyAlbum().WithTitle("Live In Japan [Disc 1]").WithArtist(
						files.NewEmptyArtist().WithFileName("George Harrison & Eric Clapton")),
					files.NewEmptyAlbum().WithTitle("Live From New York City, 1967").WithArtist(
						files.NewEmptyArtist().WithFileName("Simon & Garfunkel")),
					files.NewEmptyAlbum().WithTitle("Live At The Circle Room").WithArtist(
						files.NewEmptyArtist().WithFileName("Nat King Cole")),
					files.NewEmptyAlbum().WithTitle("Live At The BBC [Disc 2]").WithArtist(
						files.NewEmptyArtist().WithFileName("The Beatles")),
					files.NewEmptyAlbum().WithTitle("Live At The BBC [Disc 1]").WithArtist(
						files.NewEmptyArtist().WithFileName("The Beatles")),
					files.NewEmptyAlbum().WithTitle("Live 1975-85 [Disc 3]").WithArtist(
						files.NewEmptyArtist().WithFileName("Bruce Springsteen & The E Street Band")),
					files.NewEmptyAlbum().WithTitle("Live 1975-85 [Disc 2]").WithArtist(
						files.NewEmptyArtist().WithFileName("Bruce Springsteen & The E Street Band")),
					files.NewEmptyAlbum().WithTitle("Live 1975-85 [Disc 1]").WithArtist(
						files.NewEmptyArtist().WithFileName("Bruce Springsteen & The E Street Band")),
					files.NewEmptyAlbum().WithTitle("Live").WithArtist(
						files.NewEmptyArtist().WithFileName("Roger Whittaker")),
					files.NewEmptyAlbum().WithTitle("Live").WithArtist(
						files.NewEmptyArtist().WithFileName("Blondie")),
					files.NewEmptyAlbum().WithTitle("Live").WithArtist(
						files.NewEmptyArtist().WithFileName("Big Bad Voodoo Daddy")),
				},
				tab: 0,
			},
			WantedRecording: output.WantedRecording{
				Console: "" +
					`Album: "Live" by "Big Bad Voodoo Daddy"
Album: "Live" by "Blondie"
Album: "Live" by "Roger Whittaker"
Album: "Live 1975-85 [Disc 1]" by "Bruce Springsteen & The E Street Band"
Album: "Live 1975-85 [Disc 2]" by "Bruce Springsteen & The E Street Band"
Album: "Live 1975-85 [Disc 3]" by "Bruce Springsteen & The E Street Band"
Album: "Live At The BBC [Disc 1]" by "The Beatles"
Album: "Live At The BBC [Disc 2]" by "The Beatles"
Album: "Live At The Circle Room" by "Nat King Cole"
Album: "Live From New York City, 1967" by "Simon & Garfunkel"
Album: "Live In Japan [Disc 1]" by "George Harrison & Eric Clapton"
Album: "Live In Japan [Disc 2]" by "George Harrison & Eric Clapton"
Album: "Live In Paris & Toronto [Disc 1]" by "Loreena McKennitt"
Album: "Live In Paris & Toronto [Disc 2]" by "Loreena McKennitt"
Album: "Live Rhymin' [Bonus Tracks]" by "Paul Simon"
`,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.ls.ListAlbums(o, tt.args.albums, tt.args.tab)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.ListAlbums() %s", difference)
				}
			}
		})
	}
}

func TestListSettingsListArtists(t *testing.T) {
	tests := map[string]struct {
		ls      *cmd.ListSettings
		artists []*files.Artist
		output.WantedRecording
	}{
		"no artists": {ls: cmd.NewListSettings()},
		"tracks": {
			ls: cmd.NewListSettings().WithTracks(true).WithAnnotate(
				true).WithSortByTitle(true),
			artists: generateArtists(3, 3, 3),
			WantedRecording: output.WantedRecording{
				Console: "" +
					"\"my track 001\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 002\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 003\" on \"my album 00\" by \"my artist 0\"\n" +
					"\"my track 011\" on \"my album 01\" by \"my artist 0\"\n" +
					"\"my track 012\" on \"my album 01\" by \"my artist 0\"\n" +
					"\"my track 013\" on \"my album 01\" by \"my artist 0\"\n" +
					"\"my track 021\" on \"my album 02\" by \"my artist 0\"\n" +
					"\"my track 022\" on \"my album 02\" by \"my artist 0\"\n" +
					"\"my track 023\" on \"my album 02\" by \"my artist 0\"\n" +
					"\"my track 101\" on \"my album 10\" by \"my artist 1\"\n" +
					"\"my track 102\" on \"my album 10\" by \"my artist 1\"\n" +
					"\"my track 103\" on \"my album 10\" by \"my artist 1\"\n" +
					"\"my track 111\" on \"my album 11\" by \"my artist 1\"\n" +
					"\"my track 112\" on \"my album 11\" by \"my artist 1\"\n" +
					"\"my track 113\" on \"my album 11\" by \"my artist 1\"\n" +
					"\"my track 121\" on \"my album 12\" by \"my artist 1\"\n" +
					"\"my track 122\" on \"my album 12\" by \"my artist 1\"\n" +
					"\"my track 123\" on \"my album 12\" by \"my artist 1\"\n" +
					"\"my track 201\" on \"my album 20\" by \"my artist 2\"\n" +
					"\"my track 202\" on \"my album 20\" by \"my artist 2\"\n" +
					"\"my track 203\" on \"my album 20\" by \"my artist 2\"\n" +
					"\"my track 211\" on \"my album 21\" by \"my artist 2\"\n" +
					"\"my track 212\" on \"my album 21\" by \"my artist 2\"\n" +
					"\"my track 213\" on \"my album 21\" by \"my artist 2\"\n" +
					"\"my track 221\" on \"my album 22\" by \"my artist 2\"\n" +
					"\"my track 222\" on \"my album 22\" by \"my artist 2\"\n" +
					"\"my track 223\" on \"my album 22\" by \"my artist 2\"\n",
			},
		},
		"albums": {
			ls:      cmd.NewListSettings().WithAlbums(true).WithAnnotate(true),
			artists: generateArtists(3, 3, 3),
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Album: \"my album 00\" by \"my artist 0\"\n" +
					"Album: \"my album 01\" by \"my artist 0\"\n" +
					"Album: \"my album 02\" by \"my artist 0\"\n" +
					"Album: \"my album 10\" by \"my artist 1\"\n" +
					"Album: \"my album 11\" by \"my artist 1\"\n" +
					"Album: \"my album 12\" by \"my artist 1\"\n" +
					"Album: \"my album 20\" by \"my artist 2\"\n" +
					"Album: \"my album 21\" by \"my artist 2\"\n" +
					"Album: \"my album 22\" by \"my artist 2\"\n",
			},
		},
		"albums and tracks": {
			ls: cmd.NewListSettings().WithAlbums(true).WithTracks(
				true).WithAnnotate(true).WithSortByNumber(true),
			artists: generateArtists(3, 3, 3),
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Album: \"my album 00\" by \"my artist 0\"\n" +
					"   1. my track 001\n" +
					"   2. my track 002\n" +
					"   3. my track 003\n" +
					"Album: \"my album 01\" by \"my artist 0\"\n" +
					"   1. my track 011\n" +
					"   2. my track 012\n" +
					"   3. my track 013\n" +
					"Album: \"my album 02\" by \"my artist 0\"\n" +
					"   1. my track 021\n" +
					"   2. my track 022\n" +
					"   3. my track 023\n" +
					"Album: \"my album 10\" by \"my artist 1\"\n" +
					"   1. my track 101\n" +
					"   2. my track 102\n" +
					"   3. my track 103\n" +
					"Album: \"my album 11\" by \"my artist 1\"\n" +
					"   1. my track 111\n" +
					"   2. my track 112\n" +
					"   3. my track 113\n" +
					"Album: \"my album 12\" by \"my artist 1\"\n" +
					"   1. my track 121\n" +
					"   2. my track 122\n" +
					"   3. my track 123\n" +
					"Album: \"my album 20\" by \"my artist 2\"\n" +
					"   1. my track 201\n" +
					"   2. my track 202\n" +
					"   3. my track 203\n" +
					"Album: \"my album 21\" by \"my artist 2\"\n" +
					"   1. my track 211\n" +
					"   2. my track 212\n" +
					"   3. my track 213\n" +
					"Album: \"my album 22\" by \"my artist 2\"\n" +
					"   1. my track 221\n" +
					"   2. my track 222\n" +
					"   3. my track 223\n",
			},
		},
		"artists": {
			ls:      cmd.NewListSettings().WithArtists(true),
			artists: generateArtists(3, 3, 3),
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist: my artist 0\n" +
					"Artist: my artist 1\n" +
					"Artist: my artist 2\n",
			},
		},
		"artists and tracks": {
			ls: cmd.NewListSettings().WithArtists(true).WithTracks(
				true).WithAnnotate(true).WithSortByTitle(true),
			artists: generateArtists(3, 3, 3),
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist: my artist 0\n" +
					"  \"my track 001\" on \"my album 00\"\n" +
					"  \"my track 002\" on \"my album 00\"\n" +
					"  \"my track 003\" on \"my album 00\"\n" +
					"  \"my track 011\" on \"my album 01\"\n" +
					"  \"my track 012\" on \"my album 01\"\n" +
					"  \"my track 013\" on \"my album 01\"\n" +
					"  \"my track 021\" on \"my album 02\"\n" +
					"  \"my track 022\" on \"my album 02\"\n" +
					"  \"my track 023\" on \"my album 02\"\n" +
					"Artist: my artist 1\n" +
					"  \"my track 101\" on \"my album 10\"\n" +
					"  \"my track 102\" on \"my album 10\"\n" +
					"  \"my track 103\" on \"my album 10\"\n" +
					"  \"my track 111\" on \"my album 11\"\n" +
					"  \"my track 112\" on \"my album 11\"\n" +
					"  \"my track 113\" on \"my album 11\"\n" +
					"  \"my track 121\" on \"my album 12\"\n" +
					"  \"my track 122\" on \"my album 12\"\n" +
					"  \"my track 123\" on \"my album 12\"\n" +
					"Artist: my artist 2\n" +
					"  \"my track 201\" on \"my album 20\"\n" +
					"  \"my track 202\" on \"my album 20\"\n" +
					"  \"my track 203\" on \"my album 20\"\n" +
					"  \"my track 211\" on \"my album 21\"\n" +
					"  \"my track 212\" on \"my album 21\"\n" +
					"  \"my track 213\" on \"my album 21\"\n" +
					"  \"my track 221\" on \"my album 22\"\n" +
					"  \"my track 222\" on \"my album 22\"\n" +
					"  \"my track 223\" on \"my album 22\"\n",
			},
		},
		"artists and albums": {
			ls:      cmd.NewListSettings().WithArtists(true).WithAlbums(true),
			artists: generateArtists(3, 3, 3),
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist: my artist 0\n" +
					"  Album: my album 00\n" +
					"  Album: my album 01\n" +
					"  Album: my album 02\n" +
					"Artist: my artist 1\n" +
					"  Album: my album 10\n" +
					"  Album: my album 11\n" +
					"  Album: my album 12\n" +
					"Artist: my artist 2\n" +
					"  Album: my album 20\n" +
					"  Album: my album 21\n" +
					"  Album: my album 22\n",
			},
		},
		"everything": {
			ls: cmd.NewListSettings().WithArtists(true).WithAlbums(
				true).WithTracks(true).WithSortByNumber(true),
			artists: generateArtists(3, 3, 3),
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist: my artist 0\n" +
					"  Album: my album 00\n" +
					"     1. my track 001\n" +
					"     2. my track 002\n" +
					"     3. my track 003\n" +
					"  Album: my album 01\n" +
					"     1. my track 011\n" +
					"     2. my track 012\n" +
					"     3. my track 013\n" +
					"  Album: my album 02\n" +
					"     1. my track 021\n" +
					"     2. my track 022\n" +
					"     3. my track 023\n" +
					"Artist: my artist 1\n" +
					"  Album: my album 10\n" +
					"     1. my track 101\n" +
					"     2. my track 102\n" +
					"     3. my track 103\n" +
					"  Album: my album 11\n" +
					"     1. my track 111\n" +
					"     2. my track 112\n" +
					"     3. my track 113\n" +
					"  Album: my album 12\n" +
					"     1. my track 121\n" +
					"     2. my track 122\n" +
					"     3. my track 123\n" +
					"Artist: my artist 2\n" +
					"  Album: my album 20\n" +
					"     1. my track 201\n" +
					"     2. my track 202\n" +
					"     3. my track 203\n" +
					"  Album: my album 21\n" +
					"     1. my track 211\n" +
					"     2. my track 212\n" +
					"     3. my track 213\n" +
					"  Album: my album 22\n" +
					"     1. my track 221\n" +
					"     2. my track 222\n" +
					"     3. my track 223\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			tt.ls.ListArtists(o, tt.artists)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.ListArtists() %s", difference)
				}
			}
		})
	}
}

func Test_ListRun(t *testing.T) {
	cmd.InitGlobals()
	originalBus := cmd.Bus
	originalSearchFlags := cmd.SearchFlags
	defer func() {
		cmd.Bus = originalBus
		cmd.SearchFlags = originalSearchFlags
	}()
	cmd.SearchFlags = safeSearchFlags

	testListFlags := cmd.NewSectionFlags().WithSectionName(cmd.ListCommand).WithFlags(
		map[string]*cmd.FlagDetails{
			cmd.ListAlbums: cmd.NewFlagDetails().WithAbbreviatedName(
				"l").WithUsage("include album names in listing").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListArtists: cmd.NewFlagDetails().WithAbbreviatedName(
				"r").WithUsage("include artist names in listing").WithExpectedType(
				cmd.BoolType).WithDefaultValue(true),
			cmd.ListTracks: cmd.NewFlagDetails().WithAbbreviatedName(
				"t").WithUsage("include track names in listing").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListSortByNumber: cmd.NewFlagDetails().WithUsage(
				"sort tracks by track number").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListSortByTitle: cmd.NewFlagDetails().WithUsage(
				"sort tracks by track title").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListAnnotate: cmd.NewFlagDetails().WithUsage(
				"annotate listings with album and artist names").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListDetails: cmd.NewFlagDetails().WithUsage(
				"include details with tracks").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListDiagnostic: cmd.NewFlagDetails().WithUsage(
				"include diagnostic information with tracks").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
		},
	)
	testCmd := &cobra.Command{}
	cmd.AddFlags(output.NewNilBus(), cmd_toolkit.EmptyConfiguration(), testCmd.Flags(),
		testListFlags, cmd.SearchFlags)

	testListFlags2 := cmd.NewSectionFlags().WithSectionName(cmd.ListCommand).WithFlags(
		map[string]*cmd.FlagDetails{
			cmd.ListAlbums: cmd.NewFlagDetails().WithAbbreviatedName(
				"l").WithUsage("include album names in listing").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListArtists: cmd.NewFlagDetails().WithAbbreviatedName(
				"r").WithUsage("include artist names in listing").WithExpectedType(
				cmd.BoolType).WithDefaultValue(true),
			cmd.ListTracks: cmd.NewFlagDetails().WithAbbreviatedName(
				"t").WithUsage("include track names in listing").WithExpectedType(
				cmd.BoolType).WithDefaultValue(true),
			cmd.ListSortByNumber: cmd.NewFlagDetails().WithUsage(
				"sort tracks by track number").WithExpectedType(
				cmd.BoolType).WithDefaultValue(true),
			cmd.ListSortByTitle: cmd.NewFlagDetails().WithUsage(
				"sort tracks by track title").WithExpectedType(
				cmd.BoolType).WithDefaultValue(true),
			cmd.ListAnnotate: cmd.NewFlagDetails().WithUsage(
				"annotate listings with album and artist names").WithExpectedType(cmd.BoolType).WithDefaultValue(false),
			cmd.ListDetails: cmd.NewFlagDetails().WithUsage(
				"include details with tracks").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListDiagnostic: cmd.NewFlagDetails().WithUsage(
				"include diagnostic information with tracks").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
		},
	)
	testCmd2 := &cobra.Command{}
	cmd.AddFlags(output.NewNilBus(), cmd_toolkit.EmptyConfiguration(), testCmd2.Flags(),
		testListFlags2, cmd.SearchFlags)

	testListFlags3 := cmd.NewSectionFlags().WithSectionName(cmd.ListCommand).WithFlags(
		map[string]*cmd.FlagDetails{
			cmd.ListAlbums: cmd.NewFlagDetails().WithAbbreviatedName(
				"l").WithUsage("include album names in listing").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListArtists: cmd.NewFlagDetails().WithAbbreviatedName(
				"r").WithUsage("include artist names in listing").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListTracks: cmd.NewFlagDetails().WithAbbreviatedName(
				"t").WithUsage("include track names in listing").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListSortByNumber: cmd.NewFlagDetails().WithUsage(
				"sort tracks by track number").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListSortByTitle: cmd.NewFlagDetails().WithUsage(
				"sort tracks by track title").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListAnnotate: cmd.NewFlagDetails().WithUsage(
				"annotate listings with album and artist names").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListDetails: cmd.NewFlagDetails().WithUsage(
				"include details with tracks").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
			cmd.ListDiagnostic: cmd.NewFlagDetails().WithUsage(
				"include diagnostic information with tracks").WithExpectedType(
				cmd.BoolType).WithDefaultValue(false),
		},
	)
	testCmd3 := &cobra.Command{}
	cmd.AddFlags(output.NewNilBus(), cmd_toolkit.EmptyConfiguration(), testCmd3.Flags(),
		testListFlags3, cmd.SearchFlags)

	tests := map[string]struct {
		cmd *cobra.Command
		in1 []string
		output.WantedRecording
	}{
		"typical": {
			cmd: testCmd,
			in1: nil,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No music files could be found using the specified parameters.\n" +
					"Why?\n" +
					"There were no directories found in \".\" (the --topDir value).\n" +
					"What to do:\n" +
					"Set --topDir to the path of a directory that contains artist" +
					" directories.\n",
				Log: "" +
					"level='info'" +
					" --albumFilter='.*'" +
					" --albums='false'" +
					" --annotate='false'" +
					" --artistFilter='.*'" +
					" --artists='true'" +
					" --byNumber='false'" +
					" --byTitle='false'" +
					" --details='false'" +
					" --diagnostic='false'" +
					" --extensions='[.mp3]'" +
					" --topDir='.'" +
					" --trackFilter='.*'" +
					" --tracks='false'" +
					" albums-user-set='false'" +
					" artists-user-set='false'" +
					" byNumber-user-set='false'" +
					" byTitle-user-set='false'" +
					" command='list'" +
					" tracks-user-set='false'" +
					" msg='executing command'\n" +
					"level='error'" +
					" --topDir='.'" +
					" msg='cannot find any artist directories'\n",
			},
		},
		"typical but sorting is screwy": {
			cmd: testCmd2,
			in1: nil,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"Track sorting cannot be done.\n" +
					"Why?\n" +
					"The --byNumber and --byTitle flags are both configured true.\n" +
					"What to do:\n" +
					"Either edit the configuration file and use those default values, or" +
					" use appropriate command line values.\n",
				Log: "" +
					"level='info'" +
					" --albumFilter='.*'" +
					" --albums='false'" +
					" --annotate='false'" +
					" --artistFilter='.*'" +
					" --artists='true'" +
					" --byNumber='true'" +
					" --byTitle='true'" +
					" --details='false'" +
					" --diagnostic='false'" +
					" --extensions='[.mp3]'" +
					" --topDir='.'" +
					" --trackFilter='.*'" +
					" --tracks='true'" +
					" albums-user-set='false'" +
					" artists-user-set='false'" +
					" byNumber-user-set='false'" +
					" byTitle-user-set='false'" +
					" command='list'" +
					" tracks-user-set='false'" +
					" msg='executing command'\n",
			},
		},
		"no work to do": {
			cmd: testCmd3,
			in1: nil,
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No listing will be output.\n" +
					"Why?\n" +
					"The flags --albums, --artists, and --tracks are all configured" +
					" false.\n" +
					"What to do:\n" +
					"Either:\n" +
					"[1] Edit the configuration file so that at least one of these flags" +
					" is true, or\n" +
					"[2] explicitly set at least one of these flags true on the command" +
					" line.\n",
				Log: "" +
					"level='info'" +
					" --albumFilter='.*'" +
					" --albums='false'" +
					" --annotate='false'" +
					" --artistFilter='.*'" +
					" --artists='false'" +
					" --byNumber='false'" +
					" --byTitle='false'" +
					" --details='false'" +
					" --diagnostic='false'" +
					" --extensions='[.mp3]'" +
					" --topDir='.'" +
					" --trackFilter='.*'" +
					" --tracks='false'" +
					" albums-user-set='false'" +
					" artists-user-set='false'" +
					" byNumber-user-set='false'" +
					" byTitle-user-set='false'" +
					" command='list'" +
					" tracks-user-set='false'" +
					" msg='executing command'\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			cmd.Bus = o // cook getBus()
			cmd.ListRun(tt.cmd, tt.in1)
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListRun() %s", difference)
				}
			}
		})
	}
}

func compareExitErrors(e1, e2 *cmd.ExitError) bool {
	if e1 == nil {
		return e2 == nil
	}
	if e2 == nil {
		return false
	}
	return e1.Error() == e2.Error()
}

func TestListSettingsProcessArtists(t *testing.T) {
	type args struct {
		allArtists     []*files.Artist
		loaded         bool
		searchSettings *cmd.SearchSettings
	}
	tests := map[string]struct {
		ls *cmd.ListSettings
		args
		wantStatus *cmd.ExitError
		output.WantedRecording
	}{
		"no data": {
			ls: cmd.NewListSettings().WithArtists(true),
			args: args{
				allArtists: nil,
				loaded:     true,
				searchSettings: cmd.NewSearchSettings().WithArtistFilter(
					regexp.MustCompile(".*")).WithAlbumFilter(
					regexp.MustCompile(".*")).WithTrackFilter(regexp.MustCompile(".*")),
			},
			wantStatus: cmd.NewExitUserError(cmd.ListCommand),
			WantedRecording: output.WantedRecording{
				Error: "" +
					"No music files remain after filtering.\n" +
					"Why?\n" +
					"After applying --artistFilter=\".*\", --albumFilter=\".*\", and" +
					" --trackFilter=\".*\", no files remained.\n" +
					"What to do:\n" +
					"Use less restrictive filter settings.\n",
				Log: "level='error'" +
					" --albumFilter='.*'" +
					" --artistFilter='.*'" +
					" --trackFilter='.*'" +
					" msg='no files remain after filtering'\n",
			},
		},
		"with data": {
			ls: cmd.NewListSettings().WithArtists(true),
			args: args{
				allArtists: generateArtists(3, 4, 5),
				loaded:     true,
				searchSettings: cmd.NewSearchSettings().WithArtistFilter(
					regexp.MustCompile(".*")).WithAlbumFilter(
					regexp.MustCompile(".*")).WithTrackFilter(regexp.MustCompile(".*")),
			},
			wantStatus: nil,
			WantedRecording: output.WantedRecording{
				Console: "" +
					"Artist: my artist 0\n" +
					"Artist: my artist 1\n" +
					"Artist: my artist 2\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			if got := tt.ls.ProcessArtists(o, tt.args.allArtists, tt.args.loaded,
				tt.args.searchSettings); !compareExitErrors(got, tt.wantStatus) {
				t.Errorf("ListSettings.ProcessArtists() got %s want %s", got,
					tt.wantStatus)
			}
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("ListSettings.ProcessArtists() %s", difference)
				}
			}
		})
	}
}

func TestListHelp(t *testing.T) {
	originalSearchFlags := cmd.SearchFlags
	defer func() {
		cmd.SearchFlags = originalSearchFlags
	}()
	cmd.SearchFlags = safeSearchFlags
	commandUnderTest := cloneCommand(cmd.ListCmd)
	cmd.AddFlags(output.NewNilBus(), cmd_toolkit.EmptyConfiguration(),
		commandUnderTest.Flags(), cmd.ListFlags, cmd.SearchFlags)
	tests := map[string]struct {
		output.WantedRecording
	}{
		"good": {
			WantedRecording: output.WantedRecording{
				Console: "" +
					"\"list\" lists mp3 files and containing album and artist" +
					" directories\n" +
					"\n" +
					"Usage:\n" +
					"  list [--albums] [--artists] [--tracks] [--annotate] [--details]" +
					" [--diagnostic] [--byNumber | --byTitle] [--albumFilter regex]" +
					" [--artistFilter regex] [--trackFilter regex] [--topDir dir]" +
					" [--extensions extensions]\n" +
					"\n" +
					"Examples:\n" +
					"list --annotate\n" +
					"  Annotate tracks with album and artist data and albums with artist" +
					" data\n" +
					"list --details\n" +
					"  Include detailed information, if available, for each track. This" +
					" includes composer,\n" +
					"  conductor, key, lyricist, orchestra/band, and subtitle\n" +
					"list --albums\n" +
					"  Include the album names in the output\n" +
					"list --artists\n" +
					"  Include the artist names in the output\n" +
					"list --tracks\n" +
					"  Include the track names in the output\n" +
					"list --byTitle\n" +
					"  Sort tracks by name, ignoring track numbers\n" +
					"list --byNumber\n" +
					"  Sort tracks by track number\n" +
					"\n" +
					"Flags:\n" +
					"      --albumFilter string    " +
					"regular expression specifying which albums to select (default \".*\")\n" +
					"  -l, --albums                " +
					"include album names in listing (default false)\n" +
					"      --annotate              " +
					"annotate listings with album and artist names (default false)\n" +
					"      --artistFilter string   " +
					"regular expression specifying which artists to select (default \".*\")\n" +
					"  -r, --artists               " +
					"include artist names in listing (default false)\n" +
					"      --byNumber              " +
					"sort tracks by track number (default false)\n" +
					"      --byTitle               " +
					"sort tracks by track title (default false)\n" +
					"      --details               " +
					"include details with tracks (default false)\n" +
					"      --diagnostic            " +
					"include diagnostic information with tracks (default false)\n" +
					"      --extensions string     " +
					"comma-delimited list of file extensions used by mp3 files (default \".mp3\")\n" +
					"      --topDir string         " +
					"top directory specifying where to find mp3 files (default \".\")\n" +
					"      --trackFilter string    " +
					"regular expression specifying which tracks to select (default \".*\")\n" +
					"  -t, --tracks                " +
					"include track names in listing (default false)\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := output.NewRecorder()
			command := commandUnderTest
			enableCommandRecording(o, command)
			command.Help()
			if differences, ok := o.Verify(tt.WantedRecording); !ok {
				for _, difference := range differences {
					t.Errorf("list Help() %s", difference)
				}
			}
		})
	}
}

func TestTrackSliceSort(t *testing.T) {
	tests := map[string]struct {
		ts   []*files.Track
		want []*files.Track
	}{
		"https://github.com/majohn-r/mp3/issues/147": {
			ts: []*files.Track{
				files.NewEmptyTrack().WithName("b").WithAlbum(
					files.NewEmptyAlbum().WithTitle("b").WithArtist(
						files.NewEmptyArtist().WithFileName("c"))),
				files.NewEmptyTrack().WithName("b").WithAlbum(
					files.NewEmptyAlbum().WithTitle("a").WithArtist(
						files.NewEmptyArtist().WithFileName("c"))),
				files.NewEmptyTrack().WithName("b").WithAlbum(
					files.NewEmptyAlbum().WithTitle("b").WithArtist(
						files.NewEmptyArtist().WithFileName("a"))),
				files.NewEmptyTrack().WithName("a").WithAlbum(
					files.NewEmptyAlbum().WithTitle("b").WithArtist(
						files.NewEmptyArtist().WithFileName("c"))),
				files.NewEmptyTrack().WithName("a").WithAlbum(
					files.NewEmptyAlbum().WithTitle("a").WithArtist(
						files.NewEmptyArtist().WithFileName("c"))),
				files.NewEmptyTrack().WithName("a").WithAlbum(
					files.NewEmptyAlbum().WithTitle("b").WithArtist(
						files.NewEmptyArtist().WithFileName("a"))),
			},
			want: []*files.Track{
				files.NewEmptyTrack().WithName("a").WithAlbum(
					files.NewEmptyAlbum().WithTitle("a").WithArtist(
						files.NewEmptyArtist().WithFileName("c"))),
				files.NewEmptyTrack().WithName("a").WithAlbum(
					files.NewEmptyAlbum().WithTitle("b").WithArtist(
						files.NewEmptyArtist().WithFileName("a"))),
				files.NewEmptyTrack().WithName("a").WithAlbum(
					files.NewEmptyAlbum().WithTitle("b").WithArtist(
						files.NewEmptyArtist().WithFileName("c"))),
				files.NewEmptyTrack().WithName("b").WithAlbum(
					files.NewEmptyAlbum().WithTitle("a").WithArtist(
						files.NewEmptyArtist().WithFileName("c"))),
				files.NewEmptyTrack().WithName("b").WithAlbum(
					files.NewEmptyAlbum().WithTitle("b").WithArtist(
						files.NewEmptyArtist().WithFileName("a"))),
				files.NewEmptyTrack().WithName("b").WithAlbum(
					files.NewEmptyAlbum().WithTitle("b").WithArtist(
						files.NewEmptyArtist().WithFileName("c"))),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sort.Sort(cmd.TrackSlice(tt.ts))
			if !reflect.DeepEqual(tt.ts, tt.want) {
				t.Errorf("TrackSlice.Sort = %v, want %v", tt.ts, tt.want)
			}
		})
	}
}

func TestAlbumSliceSort(t *testing.T) {
	tests := map[string]struct {
		ts   []*files.Album
		want []*files.Album
	}{
		"https://github.com/majohn-r/mp3/issues/147": {
			ts: []*files.Album{
				files.NewEmptyAlbum().WithTitle("b").WithArtist(
					files.NewEmptyArtist().WithFileName("c")),
				files.NewEmptyAlbum().WithTitle("a").WithArtist(
					files.NewEmptyArtist().WithFileName("c")),
				files.NewEmptyAlbum().WithTitle("b").WithArtist(
					files.NewEmptyArtist().WithFileName("a")),
			},
			want: []*files.Album{
				files.NewEmptyAlbum().WithTitle("a").WithArtist(
					files.NewEmptyArtist().WithFileName("c")),
				files.NewEmptyAlbum().WithTitle("b").WithArtist(
					files.NewEmptyArtist().WithFileName("a")),
				files.NewEmptyAlbum().WithTitle("b").WithArtist(
					files.NewEmptyArtist().WithFileName("c")),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sort.Sort(cmd.AlbumSlice(tt.ts))
			if !reflect.DeepEqual(tt.ts, tt.want) {
				t.Errorf("TrackSlice.Sort = %v, want %v", tt.ts, tt.want)
			}
		})
	}
}
