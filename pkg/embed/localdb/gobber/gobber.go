package gobber

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
)

type GobReader struct {
	reader io.Reader
}

type GobWriter struct {
	writer io.Writer
}

func NewGobReaderIO(r io.Reader) *GobReader {
	return &GobReader{reader: r}
}

func (gr *GobReader) Read(p interface{}) error {
	var length int64
	err := binary.Read(gr.reader, binary.LittleEndian, &length)
	if err != nil {
		return err
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(gr.reader, buf)
	if err != nil {
		return err
	}
	bReader := bytes.NewReader(buf)
	decoder := gob.NewDecoder(bReader)
	return decoder.Decode(p)
}

func NewGobWriterIO(w io.Writer) *GobWriter {
	return &GobWriter{writer: w}
}

func (gw *GobWriter) Write(p interface{}) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(p)
	if err != nil {
		return err
	}

	// Write the length of the buffer
	err = binary.Write(gw.writer, binary.LittleEndian, int64(buf.Len()))
	if err != nil {
		return err
	}

	// Write the buffer itself
	_, err = gw.writer.Write(buf.Bytes())
	return err
}
