package files

import (
	"fmt"
	"mp3/internal"
	"path/filepath"
	"sort"
)

// NOTE: the functions in this file are strictly for testing purposes. Do not
// call them from production code.

// CreateAllArtistsForTesting creates a well-defined slice of artists, with
// albums and tracks.
func CreateAllArtistsForTesting(topDir string, addExtras bool) []*Artist {
	var artists []*Artist
	for k := 0; k < 10; k++ {
		artistName := internal.CreateArtistNameForTesting(k)
		artistDir := filepath.Join(topDir, artistName)
		artist := NewArtist(artistName, artistDir)
		for n := 0; n < 10; n++ {
			albumName := internal.CreateAlbumNameForTesting(n)
			albumDir := filepath.Join(artistDir, albumName)
			album := NewAlbum(albumName, artist, albumDir)
			for p := 0; p < 10; p++ {
				trackName := internal.CreateTrackNameForTesting(p)
				name, trackNo, _ := parseTrackName(nil, trackName, album, defaultFileExtension)
				album.AddTrack(NewTrack(album, trackName, name, trackNo))
			}
			artist.AddAlbum(album)
		}
		if addExtras {
			albumName := internal.CreateAlbumNameForTesting(999)
			album := NewAlbum(albumName, artist, artist.subDirectory(albumName))
			artist.AddAlbum(album)
		}
		artists = append(artists, artist)
	}
	if addExtras {
		artistName := internal.CreateArtistNameForTesting(999)
		artist := NewArtist(artistName, filepath.Join(topDir, artistName))
		artists = append(artists, artist)
	}
	return artists
}

var (
	nameToID3V2Tag = map[string]string{
		"artist": "TPE1",
		"album":  "TALB",
		"title":  "TIT2",
		"genre":  "TCON",
		"year":   "TYER",
		"track":  "TRCK",
	}
	recognizedTagNames = []string{"artist", "album", "title", "genre", "year", "track"}
)

// CreateID3V2TaggedDataForTesting creates ID3V2-tagged content. This code is
// based on reading https://id3.org/id3v2.3.0 and on looking at a hex dump of a
// real mp3 file.
func CreateID3V2TaggedDataForTesting(audio []byte, frames map[string]string) []byte {
	content := make([]byte, 0)
	// block off tag header
	content = append(content, []byte("ID3")...)
	content = append(content, []byte{3, 0, 0, 0, 0, 0, 0}...)
	// add some text frames; order is fixed for testing
	var keys []string
	for key := range frames {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		content = append(content, makeTextFrame(key, frames[key])...)
	}
	contentLength := len(content) - 10
	factor := 128 * 128 * 128
	for k := 0; k < 4; k++ {
		content[6+k] = byte(contentLength / factor)
		contentLength %= factor
		factor /= 128
	}
	// add payload
	content = append(content, audio...)
	return content
}

func makeTextFrame(id, content string) []byte {
	frame := make([]byte, 0)
	frame = append(frame, []byte(id)...)
	contentSize := 1 + len(content)
	factor := 256 * 256 * 256
	for k := 0; k < 4; k++ {
		frame = append(frame, byte(contentSize/factor))
		contentSize %= factor
		factor /= 256
	}
	frame = append(frame, []byte{0, 0, 0}...)
	frame = append(frame, []byte(content)...)
	return frame
}

// CreateConsistentlyTaggedDataForTesting creates a file with a consistent set
// of ID3V2 and ID3V1 tags
func CreateConsistentlyTaggedDataForTesting(audio []byte, m map[string]any) []byte {
	var frames = map[string]string{}
	for _, tagName := range recognizedTagNames {
		if value, ok := m[tagName]; ok {
			switch tagName {
			case "track":
				frames[nameToID3V2Tag[tagName]] = fmt.Sprintf("%d", value.(int))
			default:
				frames[nameToID3V2Tag[tagName]] = value.(string)
			}
		}
	}
	data := CreateID3V2TaggedDataForTesting(audio, frames)
	data = append(data, createID3V1TaggedDataForTesting(m)...)
	return data
}

func createID3V1TaggedDataForTesting(m map[string]any) []byte {
	v1 := newID3v1Metadata()
	v1.writeString("TAG", id3v1Tag)
	for _, tagName := range recognizedTagNames {
		if value, ok := m[tagName]; ok {
			switch tagName {
			case "artist":
				v1.setArtist(value.(string))
			case "album":
				v1.setAlbum(value.(string))
			case "title":
				v1.setTitle(value.(string))
			case "genre":
				v1.setGenre(value.(string))
			case "year":
				v1.setYear(value.(string))
			case "track":
				v1.setTrack(value.(int))
			}
		}
	}
	return v1.data
}
