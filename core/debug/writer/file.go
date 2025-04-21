package writer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileWriter struct {
	file         *os.File
	FilenameFunc func(*FileWriter) string
	num          uint
	opts         Options
	input        chan []byte
}

const timeFormat = "2025-04-21"

func NewFileWriter(opts ...Option) (*FileWriter, error) {
	p := &FileWriter{
		opts: setup(),
	}

	for _, o := range opts {
		o(&p.opts)
	}

	var filename string

	for {
		filename = p.getFilename()
		_, err := os.Stat(filename)

		if err != nil {
			if os.IsNotExist(err) {
				if p.num > 0 {
					p.num--
					filename = p.getFilename()
				}
				break
			}
			return nil, err
		}

		p.num++

		if p.opts.cap == 0 {
			break
		}
	}

	var err error
	p.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		return nil, err
	}

	p.input = make(chan []byte, 100)
	go p.write()

	return p, nil
}

func (p *FileWriter) getFilename() string {
	if p.FilenameFunc != nil {
		return p.FilenameFunc(p)
	}

	if p.opts.cap == 0 {
		return filepath.Join(
			p.opts.path,
			fmt.Sprintf("%s.%s", time.Now().Format(timeFormat), p.opts.suffix),
		)
	}

	return filepath.Join(
		p.opts.path,
		fmt.Sprintf("%s-[%d].%s", time.Now().Format(timeFormat), p.num, p.opts.suffix),
	)
}

func (p *FileWriter) write() {
	for {
		select {
		case d := <-p.input:
			_, err := p.file.Write(d)
			if err != nil {
				log.Printf("write file error: %s\n", err.Error())
			}
			p.checkFile()
		}
	}
}

func (p *FileWriter) checkFile() {
	info, _ := p.file.Stat()
	if strings.Index(p.file.Name(), time.Now().Format(timeFormat)) < 0 ||
		(p.opts.cap > 0 && uint(info.Size()) > p.opts.cap) {
		if uint(info.Size()) > p.opts.cap {
			p.num++
		} else {
			p.num = 0
		}

		filename := p.getFilename()
		_ = p.file.Close()
		p.file, _ = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	}
}

func (p *FileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}

	if p.file == nil {
		return 0, errors.New("file not opened")
	}

	n = len(data)
	go func() {
		p.input <- data
	}()

	return n, nil
}
