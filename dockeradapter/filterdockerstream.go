package dockeradapter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/puppetlabs/lumogon/logging"
)

const (
	// Stdin represents standard input stream type.
	Stdin = iota
	// Stdout represents standard output stream type.
	Stdout
	// Stderr represents standard error steam type.
	Stderr
	// Systemerr represents errors originating from the system that make it
	// into the the multiplexed stream.
	Systemerr

	stdWriterPrefixLen = 8
	stdWriterSizeIndex = 4
)

// FilterDockerStream reads the requested stream from *bufio.Reader, strips the
// prefix headers added by Dockers stdWriter.Write(p []byte) and returns a slice
// of strings for each line received (splitting on '/n').
//
// If the stream processes without error but no lines are found then it will
// return an empty slice of strings.
//
// The following streamTypes are supported (from the Docker stdcopy package):
//  - Stdin    - 0
//  - Stdout   - 1
//  - Stderr   - 2
//  - Systemerr - 3
// The function works by walking the supplied stream alternating between reading
// the prefix header bytes followed by the number of bytes specified in the prefix
// size bytes.
func FilterDockerStream(reader io.Reader, streamType int) ([]string, error) {
	logging.Stderr("[filterdockerstream] filtering reader for steam type: ", streamType)
	var (
		result        = []string{}
		prefix        = make([]byte, stdWriterPrefixLen)
		payloadLength int
		// TODO how come I have to declare this rather than use the assignment
		// operator in the loop below? Doing so causes the function to return
		// and empty slice?
		bytesReturned int
	)
	_, err := reader.Read(prefix)
	for err == nil {
		payloadLength = int(binary.BigEndian.Uint32(prefix[stdWriterSizeIndex:]))
		payload := make([]byte, payloadLength)
		bytesReturned, err = reader.Read(payload)
		if err != nil {
			return nil, fmt.Errorf("Error reading payload")
		}
		if bytesReturned != payloadLength {
			return nil,
				fmt.Errorf("Bytes returned from stream [%d] does not match the length specified in the header [%d]",
					bytesReturned,
					payloadLength,
				)
		}
		lines := bytes.Split(payload, []byte("\n"))
		// Only extract the requested stream type
		if int(prefix[0]) == streamType {
			for _, line := range lines {
				if len(line) != 0 {
					logging.Stderr("[filterdockerstream] read line from stream: ", string(line))
					result = append(result, string(line))
				}
			}
		}
		_, err = reader.Read(prefix)
		if err != nil {
			continue
		}
	}

	return result, nil
}
