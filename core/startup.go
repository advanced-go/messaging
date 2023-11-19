package core

import (
	"errors"
	"fmt"
)

var Directory = NewEntryDirectory()

// Register - function to register a startup uri
func Register(uri string, c chan Message) error {
	if uri == "" {
		return errors.New("invalid argument: uri is empty")
	}
	if c == nil {
		return errors.New(fmt.Sprintf("invalid argument: channel is nil for [%v]", uri))
	}
	RegisterUnchecked(uri, c)
	return nil
}

func RegisterUnchecked(uri string, c chan Message) error {
	Directory.Add(uri, c)
	return nil
}

// Shutdown - startup shutdown
func Shutdown() {
	Directory.shutdown()
}
