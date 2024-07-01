package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var SymbolReplacement = map[string]string{
	".":  "dot",
	"/":  "slash",
	";":  "semicolon",
	"-":  "dash",
	"?":  "question mark",
	":":  "colon",
	"_":  "underscore",
	"\"": "quote",
	"(":  "left parent",
	")":  "right parent",
	",":  "comma",
	"'":  "apostopry",
	"=":  "equals",
	"*":  "star",
	"!":  "exclaim",
	"&":  "and",
}

// SpeakLetterSpeaks a letter.
func SpeakLetter(synt *Synthesizer, args []string) {
	runeArray := []rune(args[0])

	if unicode.IsUpper(runeArray[0]) {
		synt.Session.Speak(fmt.Sprintf("Cap %s", args[0]))
		return
	}
	synt.Session.Speak(args[0])
}

// Dispatch Sends queue to speech dispatcher.
func Dispatch(s *Synthesizer, args []string) {
	if len(s.text) == 0 || s.isProcessing {
		return
	}

	fullText := strings.Join(s.text, "")
	if s.isPunctuationAll {
		for k, v := range SymbolReplacement {
			fullText = strings.Replace(fullText, k, fmt.Sprintf(" %s ", v), -1)
		}
	} else {
		fullText = strings.Replace(fullText, "-", " dash ", -1)
	}

	s.Session.Speak(fullText)
	s.isProcessing = true
}

// QueueText Queue text into buffer.
func QueueText(s *Synthesizer, args []string) {
	s.text = append(s.text, args[0])
}

// StopSpeaking Stops all speech.
func StopSpeaking(s *Synthesizer, args []string) {
	s.Session.Cancel()
	s.isProcessing = false
	s.text = []string{}
}

// SetRate Sets speech rate.
func SetRate(s *Synthesizer, args []string) {
	rate, err := strconv.Atoi(args[0])
	if err != nil {
		log.Printf("Error converting value %s : %s\n", args[0], err)
		return
	}

	err = s.Session.SetRate(clipValue(rate))
	if err != nil {
		log.Printf("Error while set speech rate: %s\n", err)
	}
}

// SetPitch Set speech pitch.
func SetPitch(s *Synthesizer, args []string) {
	value, err := strconv.Atoi(args[0])
	if err != nil {
		log.Printf("Error converting value %s : %s\n", args[0], err)
		return
	}

	err = s.Session.SetPitch(value)
	if err != nil {
		log.Printf("Error while set pitch value: %s\n", err)
	}
}

// Say Says something.
func Say(s *Synthesizer, args []string) {
	s.Session.Cancel()
	go s.Session.Speak(args[0])
}

// clipValue will cap any value outside range +- 100
func clipValue(val int) int {
	if val > 100 {
		return 100
	}
	if val < -100 {
		return -100
	}
	return val
}

// NotImplemented Dummy placeholder
func NotImplemented(s *Synthesizer, args []string) {

}

// SetPunctuation Set speech punctuation mode.
func SetPunctuation(synt *Synthesizer, args []string) {
	switch args[0] {
	case "some":
		synt.isPunctuationAll = false
		log.Println("Set punctuation to some")
	case "all":
		synt.isPunctuationAll = true
		log.Println("Set punctuation to all")
	default:
		log.Printf("ERROR: Unknown punctuation mode '%s'\n", args[0])
	}
}

// SetTtsSyncState Process tts_sink_state command.
func SetTtsSyncState(s *Synthesizer, args []string) {
	SetPunctuation(s, args)
	SetRate(s, []string{args[2]})
	// TODO: works on split caps
}

// QueueCode Not yet implemented.
func QueueCode(s *Synthesizer, args []string) {
	s.codeText = append(s.codeText, args...)
}

// PlayAudio Command handler to play audio file.
func PlayAudio(s *Synthesizer, args []string) {
	filename := args[0]

	if !s.audio.isOpen {
		go playOnSpeaker(filename, true)
		s.audio.isOpen = true
		return
	}
	go playOnSpeaker(filename, false)
}

// playOnSpeaker Plays audio file.
func playOnSpeaker(filename string, initSpeaker bool) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error while opening file %s with error %s\n", filename, err)
		return
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Printf("Error while opening stream for file %s with error %s\n", filename, err)
		return
	}

	if initSpeaker {
		log.Println("Initializing speaker")
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	}

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		log.Printf("Closing stream of %s\n", filename)
		streamer.Close()
	})))
}

// SetVoice Set speech voice.
func SetVoice(s *Synthesizer, args []string) {
	log.Printf("Set output module to %s and voice to %s\n", args[0], args[1])
	s.Session.SetOutputModule(args[0])
	s.Session.SetSynthVoice(args[1])
}

// SetVoiceVolume Sets speech volume.
func SetVoiceVolume(s *Synthesizer, args []string) {
	vol, err := strconv.Atoi(args[0])
	if err != nil {
		log.Printf("Value is not integer: %s\n", args[0])
		return
	}
	s.Session.SetVolume(clipValue(vol))
}
