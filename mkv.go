package rescene

import (
	"log"
	"strings"
	"time"

	"github.com/rescene/mkvparse"
)

type MyParser struct {
	Size int
}

func (p *MyParser) HandleMasterBegin(id mkvparse.ElementID, info mkvparse.ElementInfo) (bool, error) {
	p.Size = int(info.Offset + info.Size)
	log.Printf("%s- %s:\n", indent(info.Level), mkvparse.NameForElementID(id))
	return true, nil
}

func (p *MyParser) HandleMasterEnd(id mkvparse.ElementID, info mkvparse.ElementInfo) error {
	return nil
}

func (p *MyParser) HandleString(id mkvparse.ElementID, value string, info mkvparse.ElementInfo) error {
	log.Printf("%s- %v: %q\n", indent(info.Level), mkvparse.NameForElementID(id), value)
	return nil
}

func (p *MyParser) HandleInteger(id mkvparse.ElementID, value int64, info mkvparse.ElementInfo) error {
	log.Printf("%s- %v: %v\n", indent(info.Level), mkvparse.NameForElementID(id), value)
	return nil
}

func (p *MyParser) HandleFloat(id mkvparse.ElementID, value float64, info mkvparse.ElementInfo) error {
	log.Printf("%s- %v: %v\n", indent(info.Level), mkvparse.NameForElementID(id), value)
	return nil
}

func (p *MyParser) HandleDate(id mkvparse.ElementID, value time.Time, info mkvparse.ElementInfo) error {
	log.Printf("%s- %v: %v\n", indent(info.Level), mkvparse.NameForElementID(id), value)
	return nil
}

func (p *MyParser) HandleBinary(id mkvparse.ElementID, value []byte, info mkvparse.ElementInfo) error {
	if id == mkvparse.SeekIDElement {
		log.Printf("%s- %v: %x\n", indent(info.Level), mkvparse.NameForElementID(id), value)
	} else {
		log.Printf("%s- %v: <binary>\n", indent(info.Level), mkvparse.NameForElementID(id))
	}
	return nil
}

func indent(n int) string {
	return strings.Repeat("  ", n)
}
