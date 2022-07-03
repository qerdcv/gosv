package gosv

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrStructExpected = errors.New("struct expected")
)

type Writer struct {
	source           io.Writer
	fieldNames       []string
	headerFieldDelim rune
	delim            rune
	wrtHeading       bool
	isHeadingWrited  bool
}

func NewWriter(source io.Writer) *Writer {
	return &Writer{
		source:           source,
		headerFieldDelim: '.',
		delim:            ',',
		wrtHeading:       false,
		isHeadingWrited:  false,
	}
}

// SetDelimiter set delimiter of the csv file
func (w *Writer) SetDelimiter(delim rune) *Writer {
	w.delim = delim
	return w
}

// SetHeadingFieldDelim set delimiter for header field names
// Currently nested field namesa are not implemented
func (w *Writer) SetHeadingFieldDelim(delim rune) *Writer {
	w.headerFieldDelim = delim
	return w
}

// SetWriteHeading set defines write or not to write heading line
func (w *Writer) SetWriteHeading(writeHeading bool) *Writer {
	w.wrtHeading = writeHeading
	return w
}

func (w *Writer) setFieldNames(t reflect.Type) {
	numField := t.NumField()
	fieldNames := make([]string, 0, numField)

	for i := 0; i < numField; i++ {
		fieldNames = append(fieldNames, t.Field(i).Tag.Get("csv"))
	}

	w.fieldNames = fieldNames
}

func (w *Writer) writeln(line string) (n int, err error) {
	return w.source.Write([]byte(line + "\n"))
}

func (w Writer) writeDelSlice(fields []string) (n int, err error) {
	return w.writeln(strings.Join(fields, string(w.delim)))
}

func (w Writer) writeHeading() (n int, err error) {
	return w.writeDelSlice(w.fieldNames)
}

func (w *Writer) writeStruct(t reflect.Value) (n int, err error) {
	values := make([]string, 0, len(w.fieldNames))

	for fIdx := range w.fieldNames {
		val := t.Field(fIdx).Interface()

		switch val := val.(type) {
		case string:
			values = append(values, val)
		case int:
			values = append(values, strconv.FormatInt(int64(val), 10))
		case int32:
			values = append(values, strconv.FormatInt(int64(val), 10))
		case int64:
			values = append(values, strconv.FormatInt(int64(val), 10))
		case float32:
			values = append(values, strconv.FormatFloat(float64(val), 'f', 2, 32))
		case float64:
			values = append(values, strconv.FormatFloat(float64(val), 'f', 2, 64))
		case bool:
			values = append(values, strconv.FormatBool(val))
		}
	}

	return w.writeDelSlice(values)
}

// Write writes passed structure to the writer
func (w *Writer) Write(v any) (n int, err error) {
	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Struct {
		return n, ErrStructExpected
	}

	w.setFieldNames(t)

	if w.wrtHeading && !w.isHeadingWrited {
		n, err = w.writeHeading()
		if err != nil {
			return n, fmt.Errorf("write heading")
		}

		w.isHeadingWrited = true
	}

	n, err = w.writeStruct(reflect.ValueOf(v))
	if err != nil {
		return n, fmt.Errorf("write struct: %w", err)
	}

	return n, nil
}
