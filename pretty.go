package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

type RowData struct {
	header []string
	rows   [][]string
}

func (r RowData) Pretty() {
	if len(r.header) == 0 {
		return
	}
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	for _, c := range r.header {
		fmt.Fprintf(writer, "%s\t", c)
	}

	fmt.Fprintln(writer)

	for _, r := range r.rows {
		for _, c := range r {
			fmt.Fprintf(writer, "%s\t", c)
		}
		fmt.Fprintln(writer)
	}
	writer.Flush()
}

func queryResponse(r *bufio.Reader) {
	rd := &RowData{make([]string, 0), make([][]string, 0)}
	msg := ""
	tag, _ := r.ReadByte()
	rcvQueryResponseStream(tag, msg, r, rd)
}

func rcvQueryResponseStream(tag byte, msg string, r *bufio.Reader, rd *RowData) {
	for tag != 'Z' {
		switch tag {
		case 'T':
			rowDescription(msg, r, rd)
		case 'D':
			rowData(r, rd)
		case 'C':
			commandDescription(r)
		case 'E':
			errorDescription(r)
		default:
			readMsg(r)
		}
		tag, _ = r.ReadByte()
	}
	readMsg(r)
	rd.Pretty()
}

func rowDescription(msg string, r *bufio.Reader, rd *RowData) {
	messageLength := readMsgLen(4, r)
	fieldsNum := make([]byte, 2)
	io.ReadFull(r, fieldsNum)
	fnum := binary.BigEndian.Uint16(fieldsNum)

	{
		columnData := make([]byte, messageLength-2)
		io.ReadFull(r, columnData)

		k := 0
		for i := 0; i < int(fnum); i++ {
			fieldName := ""
			letter := columnData[k]
			for letter >= 65 && letter <= 122 {
				fieldName += string(letter)
				k++
				letter = columnData[k]
			}
			k += 19
			rd.header = append(rd.header, fieldName)
		}

	}
}

func rowData(r *bufio.Reader, rd *RowData) {
	readMsgLen(4, r)

	fields := make([]byte, 2)
	io.ReadFull(r, fields)
	fnum := binary.BigEndian.Uint16(fields)

	columns := make([]string, 0)
	for i := 0; i < int(fnum); i++ {
		msg := ""
		ftype := make([]byte, 4)
		io.ReadFull(r, ftype)
		ftypelen := binary.BigEndian.Uint32(ftype)
		if ftypelen == 4294967295 {
			columns = append(columns, "NULL")
			continue
		}
		for j := 0; j < int(ftypelen); j++ {
			k, _ := r.ReadByte()
			msg += string(k)
		}
		columns = append(columns, msg)

	}
	rd.rows = append(rd.rows, columns)
}

func commandDescription(r *bufio.Reader) {

	messageLength := readMsgLen(4, r)
	commandBuffer := make([]byte, messageLength)
	io.ReadFull(r, commandBuffer)

	command := string(commandBuffer)
	if strings.Contains(command, "SELECT") {
		command = strings.ReplaceAll(command, "SELECT", "Rows returned: ")
	}

	fmt.Println("\n", command)
	fmt.Println()
}

func errorDescription(r *bufio.Reader) {
	messageLength := readMsgLen(4, r)
	r.ReadByte()

	msg := ""
	for i := 0; i < messageLength-1; i++ {
		r, _ := r.ReadByte()
		msg += string(r)
	}

	fmt.Println("Error: ", msg)
}
