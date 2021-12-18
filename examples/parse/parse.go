////////////////////////////////////////////////////////////
// Extracts the cover from a .mkv file, and displays it
////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rescene/rescene"
	"golang.org/x/text/message"
)

func run() error {
	filename := os.Args[1]
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	switch filepath.Ext(filename) {
	case ".srs":
		fmt.Printf("Parse : %s (SRS)\n", filename)
		s := &rescene.SrsFile{}
		err := s.Unmarshal(b)
		if err != nil {
			return err
		}
	case ".srr":
		fmt.Printf("Parse : %s (SRR)\n", filename)
		s := &rescene.SrrFile{}
		err := s.Unmarshal(b)
		if err != nil {
			return err
		}

		fmt.Printf("Creating Application:\n\t%s\n\n", s.ApplicationName)

		if s.RarCompressed {
			fmt.Printf("SRR for compressed RARs.\n\n")
		}

		if len(s.StoredFiles) > 0 {
			p := message.NewPrinter(message.MatchLanguage("en"))
			p.Printf("Stored files:\n")
			for _, v := range s.StoredFiles {
				p.Printf("\t%9d  %s\n", len(v.Data), v.Path)
			}
			fmt.Printf("\n")
		}

		if len(s.RarFiles) > 0 {
			fmt.Printf("RAR files:\n")
			for _, v := range s.RarFiles {
				if v.CRC != 0 {
					fmt.Printf("\t%s %08X %d\n", v.Path, v.CRC, v.Size)
				} else {
					fmt.Printf("\t%s %d\n", v.Path, v.Size)
				}
			}
			fmt.Printf("\n")
		}

		if len(s.PackedFiles) > 0 {
			fmt.Printf("Archived files:\n")
			for _, v := range s.PackedFiles {
				fmt.Printf("\t%s %08X %d\n", v.Path, v.CRC, v.Size)
			}
			fmt.Printf("\n")
		}

		if len(s.OSOHashes) > 0 {
			fmt.Printf("ISDb hashes:\n")
			for _, v := range s.OSOHashes {
				fmt.Printf("\t%s %016x %d\n", v.Path, v.Hash, v.Size)
			}
			fmt.Printf("\n")
		}

		if len(s.SFVComments) > 0 {
			fmt.Printf("SFV comments:\n")
			for _, v := range s.SFVComments {
				fmt.Printf("\t%s\n", v)
			}
			fmt.Printf("\n")
		}
	default:
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
