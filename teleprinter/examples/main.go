package main

import (
	"os"
	"time"

	. "github.com/SimonRichardson/colossus/teleprinter"
	"github.com/SimonRichardson/colossus/teleprinter/logs/plaintext"
)

func main() {
	DefaultLog()

	L.Error().Println("Error: Rubber Ducky")
	L.Error().Printf("Error: Hello %s", "world!\n")
	L.Error().HR()

	L.Info().Println("Info: Rubber Ducky")
	L.Info().Printf("Info: Hello %s", "world!\n")
	L.Info().HR()

	seg := L.Error().Segment()

	seg.Println("Error Segment: Rubber Ducky")
	seg.Printf("Error Segment: Hello %s", "world!\n")
	seg.HR()

	seg.Flush()

	seg = L.Info().Segment()

	seg.Println("Info Segment: Rubber Ducky")
	seg.Printf("Info Segment: Hello %s", "world!\n")
	seg.HR()

	seg.Flush()

	// EMOJI

	emoji := plaintext.NewEmojiSync(os.Stdout)
	emoji.Error().Println("Hello :+1: World! :-1: : Should print :arrow_right: :+1:")

	timer := time.NewTimer(time.Second * 2).C
	for {
		select {
		case <-timer:
			return
		}
	}
}
