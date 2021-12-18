package rescene

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"strconv"

	"github.com/rescene/mkvparse"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	id3v1 "github.com/mikkyang/id3-go/v1"
	id3v2 "github.com/mikkyang/id3-go/v2"
	"golang.org/x/image/riff"
)

type SrsFile struct {
	Blocks []interface{}
}

type ID3v1Block struct {
	Size int
	Data []byte
}

type ID3v2Block struct {
	Size int
	Data []byte
}

type Lyrics200Block struct {
	Size int
	Data []byte
}

type Lyrics200SubBlock struct {
	Lyrics200SubHeader
	Data []byte
}

type Lyrics200SubHeader struct {
	Head [3]byte
	Len  [5]byte
}

type SrsBlock struct {
	SrsHeader
	Size int
}

type MkvBlock struct {
	Size int
	Data []byte
}

type AviBlock struct {
	Size int
	Data []byte
}

type SrsHeader struct {
	Head   [4]byte
	Length uint32
}

func (sb *Lyrics200SubBlock) Size() (int, error) {
	i, err := strconv.Atoi(string(sb.Lyrics200SubHeader.Len[:]))
	if err != nil {
		return 0, err
	} else {
		return 8 + i, nil
	}
}

func (f *SrsFile) Unmarshal(b []byte) (err error) {
	f.Blocks = make([]interface{}, 0)
	offset := 0
	for offset < len(b) {
		t, err := filetype.Get(b[offset:])
		if err != nil {
			return err
		}
		log.Printf("Offset %.5x (%.5d): %x (%s): %v (len : %d)\n", offset, offset, b[offset:offset+4], string(b[offset:offset+4]), t, len(b))

		switch t {
		case matchers.TypeMp3:
			block := &ID3v2Block{}
			err = block.Unmarshal(b[offset:])
			if err != nil {
				return err
			}
			f.Blocks = append(f.Blocks, block)
			offset += block.Size
		case TypeID3v1:
			block := &ID3v1Block{}
			err = block.Unmarshal(b[offset:])
			if err != nil {
				return err
			}
			f.Blocks = append(f.Blocks, block)
			offset += block.Size
		case TypeSrs:
			block := &SrsBlock{}
			err = block.Unmarshal(b[offset:])
			if err != nil {
				return err
			}
			f.Blocks = append(f.Blocks, block)
			offset += block.Size
			log.Printf("Block %s : Len %d\n", string(block.SrsHeader.Head[:]), block.Size)
		case TypeLyrics200:
			block := Lyrics200Block{}
			err = block.Unmarshal(b[offset:])
			if err != nil {
				return err
			}
			f.Blocks = append(f.Blocks, block)
			offset += block.Size
		case matchers.TypeMkv:
			block := MkvBlock{}
			err = block.Unmarshal(b[offset:])
			if err != nil {
				return err
			}
			f.Blocks = append(f.Blocks, block)
			offset += block.Size
		case matchers.TypeFlac:
			return nil
		case matchers.TypeAvi:
			block := AviBlock{}
			err = block.Unmarshal(b[offset:])
			if err != nil {
				return err
			}
			f.Blocks = append(f.Blocks, block)
			offset += block.Size
		default:
			return nil
		}

	}
	return nil
}

func (block *ID3v2Block) Unmarshal(b []byte) (err error) {
	buf := bytes.NewBuffer(b)
	readSeeker := bytes.NewReader(buf.Bytes())

	if v2Tag := id3v2.ParseTag(readSeeker); v2Tag != nil {
		block.Size = v2Tag.Size() + 10
		for _, f := range v2Tag.AllFrames() {
			log.Printf("id3v2 - Tag : %s : %s\n", f.Id(), f.String())
		}
	}
	return nil
}

func (block *ID3v1Block) Unmarshal(b []byte) (err error) {
	buf := bytes.NewBuffer(b)
	readSeeker := bytes.NewReader(buf.Bytes())

	if v1Tag := id3v1.ParseTag(readSeeker); v1Tag != nil {
		block.Size = v1Tag.Size()
		log.Printf("id3v1 - Tag : Title : %s\n", v1Tag.Title())
		log.Printf("id3v1 - Tag : Album : %s\n", v1Tag.Album())
		log.Printf("id3v1 - Tag : Artist : %s\n", v1Tag.Artist())
		log.Printf("id3v1 - Tag : Year : %s\n", v1Tag.Year())
		log.Printf("id3v1 - Tag : Genre : %s\n", v1Tag.Genre())
		log.Printf("id3v1 - Tag : Comments : %s\n", v1Tag.Comments())
	}
	return nil
}

func (block *SrsBlock) Unmarshal(b []byte) (err error) {
	buffer := bytes.NewBuffer(b[0:8])
	header := &SrsHeader{}
	err = binary.Read(buffer, binary.LittleEndian, header)
	if err != nil {
		return err
	}
	block.SrsHeader = *header
	block.Size = int(header.Length)
	return nil
}

func (block *Lyrics200Block) Unmarshal(b []byte) (err error) {
	offset := 11
	for offset < len(b) {
		if len(b[offset:]) < 15 {
			return io.ErrUnexpectedEOF
		}
		if b[offset+6] == 'L' && b[offset+7] == 'Y' && b[offset+8] == 'R' && b[offset+9] == 'I' && b[offset+10] == 'C' && b[offset+11] == 'S' && b[offset+12] == '2' && b[offset+13] == '0' && b[offset+14] == '0' {
			block.Data = b[0 : offset+15]
			block.Size = len(block.Data)
			return nil
		}
		buffer := bytes.NewBuffer(b[offset : offset+8])
		header := &Lyrics200SubHeader{}
		err = binary.Read(buffer, binary.LittleEndian, header)
		subblock := &Lyrics200SubBlock{
			Lyrics200SubHeader: *header,
		}
		s, err := subblock.Size()
		if err != nil {
			return err
		}
		subblock.Data = b[offset+8 : offset+s]
		offset += s
		log.Printf("Lyrics3v2 - Tag : %s : %s\n", string(subblock.Head[:]), string(subblock.Data[:]))
	}
	return nil
}

func (block *MkvBlock) Unmarshal(b []byte) (err error) {
	handler := MyParser{}
	buf := bytes.NewBuffer(b)
	err = mkvparse.Parse(buf, &handler)
	block.Size = handler.Size
	return nil
}

func (block *AviBlock) Unmarshal(b []byte) (err error) {
	buf := bytes.NewBuffer(b)
	formType, r, err := riff.NewReader(buf)
	if err != nil {
		return err
	}
	log.Printf("RIFF(%s)\n", formType)
	if err = dumpRIFF(r, ".\t"); err != nil {
		return err
	}
	return nil
}
