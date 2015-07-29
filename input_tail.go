package main

import (
	"errors"

	"github.com/ActiveState/tail"
)

func TailValidateConfig(config *ConfStruct) error {
	if config.TailPath == "" {
		return errors.New("input mode 'tail' requires VEILLE_TAIL_PATH environment variable")
	}
	return nil
}

// InputTail pulls in a file line by line and then watches for new lines
// and pulls those in too.
//
// InputTail implements the InputMode interface.
type TailInput struct {
	// Path is the filesystem path to the file that we'll be tailing
	Path string

	tail     *tail.Tail
	lineChan chan string
}

// Start makes the TailInput start reading the file and writing any new
// lines to LineChan
func (inp *TailInput) Start() (err error) {
	inp.tail, err = tail.TailFile(inp.Path, tail.Config{
		ReOpen:    true,
		MustExist: false,
		Follow:    true,
		Logger:    tail.DiscardingLogger,
	})
	if err != nil {
		return
	}
	go inp.start()
	return
}
func (inp *TailInput) start() {
	var line *tail.Line
	for line = range inp.tail.Lines {
		inp.lineChan <- line.Text
	}
}

// LineChan returns the channel to which lines read from the file will
// be written.
func (inp *TailInput) LineChan() chan string {
	return inp.lineChan
}

// NewTailInput generates a new TailInput struct.
func NewTailInput(conf *ConfStruct) *TailInput {
	return &TailInput{
		Path:     conf.TailPath,
		lineChan: make(chan string),
	}
}
