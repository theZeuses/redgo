package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file   *os.File
	parser *Reader
	mutex  sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file:   file,
		parser: NewReader(file),
	}

	go func() {
		for {
			aof.mutex.Lock()

			aof.file.Sync()

			aof.mutex.Unlock()

			// Sleep for a while to avoid busy waiting
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (a *Aof) Write(data Value) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	_, err := a.file.Write(data.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (a *Aof) Read() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for {
		value, err := a.parser.ParseFromRespString()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		ProcessCommand(value.(ArrayValue).Val[0].(BulkStringValue).Val, value.(ArrayValue).Val[1:], nil)
	}

	return nil
}

func (a *Aof) Close() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.file.Close()
}

func InitAof() (*Aof, error) {
	fmt.Println("Loading AOF file....")
	aof, err := NewAof("database.aof")
	if err != nil {
		return nil, err
	}

	err = aof.Read()
	if err != nil {
		return nil, err
	}
	fmt.Println("AOF file loaded successfully....")

	return aof, nil
}
