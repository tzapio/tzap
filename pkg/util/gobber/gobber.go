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

	if _, err = io.ReadFull(gr.reader, buf); err != nil {
		println("Warning - Binary file reader crashed. Run `tzap reset` or submit a bug report if error persists.")
		return err
	}
	bReader := bytes.NewReader(buf)
	decoder := gob.NewDecoder(bReader)

	if err := decoder.Decode(p); err != nil {
		println("Warning - Binary file reader crashed. Run `tzap reset` or submit a bug report if error persists.")
		return err
	}
	return nil
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
