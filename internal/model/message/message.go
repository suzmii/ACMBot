package message

import (
	"github.com/suzmii/ACMBot/internal/model/race"
	zMsg "github.com/wdvxdr1123/ZeroBot/message"
)

type Message interface {
	ToZeroMessage() zMsg.Message
}

type Text string

func (t Text) ToZeroMessage() zMsg.Message {
	return zMsg.Message{zMsg.Text(t)}
}

type Image []byte

func (i Image) ToZeroMessage() zMsg.Message {
	return zMsg.Message{zMsg.ImageBytes(i)}
}

type Races []race.Race

func (r Races) ToZeroMessage() zMsg.Message {
	var result string

	for _, race_ := range r[:min(8, len(r))] {
		result += "\n" + race_.NoUrlString()
	}
	return zMsg.Message{zMsg.Text(result)}
}
