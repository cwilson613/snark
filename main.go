package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.design/x/clipboard"
)

type Release struct {
	Name      string
	ID        int64
	AssetsURL string
}

type Emote struct {
	Name string
	Text string
}

var EmoteMap map[string]Emote = makeEmoteMap()

var Version string = ""
var BuildTime string = ""

func main() {
	argNum := len(os.Args)
	if argNum > 1 {
		arg1 := os.Args[1]
		switch {
		//init
		case arg1 == "list":
			listAllEmotes()
		case arg1 == "print":
			if argNum == 2 {
				log.Fatal("No name provided to `print`")
			} else if argNum == 3 {
				emoteName := os.Args[2]
				printEmote(emoteName)
			}
		case arg1 == "version":
			fmt.Println("Version:	" + Version)
			fmt.Println("Date:		" + BuildTime)
			fmt.Println("OS:		" + runtime.GOOS)
			fmt.Println("Arch:		" + runtime.GOARCH)
		default:
			if argNum == 1 {
				printDefaultUsage()
			} else if argNum == 2 {
				snarkString := os.Args[1]
				makeSarcastic(snarkString)
			}
		}

	} else {
		printDefaultUsage()
	}
}

func makeSarcastic(text string) {
	resultString := ""
	for _, c := range text {
		randSeed := rand.NewSource(time.Now().UnixNano())
		randNum := rand.New(randSeed).Intn(100)
		if randNum < 50 {
			resultString += strings.ToLower(string(c))
		} else {
			resultString += strings.ToUpper(string(c))
		}
	}
	fmt.Println(resultString)
	copyStringToClipboard(resultString)
}

func listAllEmotes() {
	for _, emote := range EmoteMap {
		fmt.Println(emote.Name + "	=>	" + emote.Text)
	}
}

func printEmote(emoteName string) {

	emoteText := EmoteMap[emoteName].Text
	fmt.Println(emoteText)
	copyStringToClipboard(emoteText)
}

func copyStringToClipboard(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}

func printDefaultUsage() {
	fmt.Printf("Usage:\n" +
		" snark		 		prints usage\n" +
		" snark version			see build and version information\n" +
		// " snark help			provides help information\n" +
		" snark list			list available emote names and values\n" +
		" snark print <emote name>	print the value for arg emote name and copies it to the clipboard\n" +
		" snark \"<string>\"		prints the given string in a sarcastic font of shIfTinG caPiTaLIsatIoN and copies to the clipboard")
}

func makeEmoteMap() map[string]Emote {
	var emoteMap = map[string]Emote{
		"lenny": {
			"lenny",
			"( ͡° ͜ʖ ͡°)",
		},
		"shrug": {
			"shrug",
			"¯\\_(ツ)_/¯",
		},
		"cat": {
			"cat",
			"(•ㅅ•)",
		},
		"shock": {
			"shock",
			"ಠ_ಠ",
		},
	}
	return emoteMap

}
