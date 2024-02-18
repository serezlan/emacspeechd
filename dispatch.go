package main

import (
	"fmt"
	"log"

	"github.com/ilyapashuk/go-speechd"
)

type Synthesizer struct {
	Session          *speechd.SpeechdSession
	audio            *AudioManager
	text             []string
	codeText         []string
	isPunctuationAll bool
	isProcessing     bool
	isSplitCaps      bool
}

func OpenSpeechSession() (*Synthesizer, error) {
	Session, err := speechd.Open()
		

	if err != nil {
		log.Println(err)
		return nil, err
	}

	audioManager := &AudioManager{isOpen: false}
	result := &Synthesizer{Session: Session, audio: audioManager}
	return result, nil
}

type AudioManager struct {
	isOpen bool
}

type commandHandler func(*Synthesizer, []string)

var commandTable = map[string]commandHandler{
	"l":                    SpeakLetter,
	"d":                    Dispatch,
	"q":                    QueueText,
	"tts_set_speech_rate":  SetRate,
	"tts_set_punctuations": SetPunctuation,
	"tts_set_voice":        SetVoice,
	"tts_set_volume":       SetVoiceVolume,
	"tts_say":              Say,
	"s":                    StopSpeaking,
	// "c":                    QueueText,
	"tts_sync_state": SetTtsSyncState,
	"p":              PlayAudio,
}

func (s *Synthesizer) ProcessCommand(cmd Command) {
	cmdDispatch, ok := commandTable[cmd.Name]
	if !ok {
		log.Printf("Undefined command '%v'", cmd)
		return
	}

	if isUnimplementedFunction(cmdDispatch) {
		log.Printf("Unimplemented command %v\n", cmd)
		return
	}
	log.Printf("command '%s' %v", cmd.Name, cmd.Args)
	cmdDispatch(s, cmd.Args)
}

func (s *Synthesizer) Close() {
	s.Session.Close()
}

// isUnimplementedFunction (fn dispat) ...
func isUnimplementedFunction(fn commandHandler) bool {
	p1 := fmt.Sprintf("%v", NotImplemented)
	p2 := fmt.Sprintf("%v", fn)

	if p1 == p2 {
		return true
	}
	return false
}
