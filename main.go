package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var DEBUG = false

func main() {
	if DEBUG {
		f, err := ioutil.TempFile("/tmp", "emacspeechd")
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetFlags(0)
		wr := logWriter{
			f: f,
		}
		log.SetOutput(wr)
	}
	c, err := OpenSpeechSession()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = c.Session.SetClientName("emacspeak", "emacspeak", "emacs")
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.Session.Speak("starting emacspeechd")
	if err != nil {
		log.Println(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic %s\n", err)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		line, _ := reader.ReadString('\n')
		cmd := Parse(line)
		c.ProcessCommand(cmd)
	}
}

type logWriter struct {
	f *os.File
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return writer.f.WriteString(time.Now().Format("15:04:05.999 ") + string(bytes))
}
