package writer

import "os"

type FileWriter struct {
	file         *os.File
	FilenameFunc func(*FileWriter) string
	num          uint
	opts         Options
	input        chan []byte
}
