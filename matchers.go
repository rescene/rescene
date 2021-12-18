package rescene

import (
	"github.com/h2non/filetype"
)

// TypeSrs for SRS files
var TypeSrs = filetype.NewType("srs", "application/resample")

// TypeID3v1 for ID3v1 tags
var TypeID3v1 = filetype.NewType("id3v1", "audio/id3v1")

// TypeLyrics200 for Lyricsv2 tags
var TypeLyrics200 = filetype.NewType("lyrics200", "audio/lyrics200")

func SrsMatcher(buf []byte) bool {
	return len(buf) > 4 && buf[0] == 'S' && buf[1] == 'R' && buf[2] == 'S' && (buf[3] == 'F' || buf[3] == 'T' || buf[3] == 'P')
}

func ID3v1Matcher(buf []byte) bool {
	return len(buf) >= 128 && buf[0] == 'T' && buf[1] == 'A' && buf[2] == 'G'
}

func Lyrics200Matcher(buf []byte) bool {
	return len(buf) >= 11 && buf[0] == 'L' && buf[1] == 'Y' && buf[2] == 'R' && buf[3] == 'I' && buf[4] == 'C' && buf[5] == 'S' && buf[6] == 'B' && buf[7] == 'E' && buf[8] == 'G' && buf[9] == 'I' && buf[10] == 'N'
}

func init() {
	filetype.AddMatcher(TypeSrs, SrsMatcher)
	filetype.AddMatcher(TypeID3v1, ID3v1Matcher)
	filetype.AddMatcher(TypeLyrics200, Lyrics200Matcher)
}
