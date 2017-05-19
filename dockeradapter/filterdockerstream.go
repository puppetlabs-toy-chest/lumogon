package dockeradapter

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"

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

// frameHeader holds the steam type and payload size for a Docker stream frame
type frameHeader struct {
	streamType  int
	payloadSize int
}

// FilterDockerStream reads the requested stream from *bufio.Reader, strips the
// frameHeader added by Dockers stdWriter.Write(p []byte) and returns a slice
// of strings for each line received (splitting on '/n').
//
// If the stream processes without error but no lines are found then it will
// return an empty slice of strings.
//
// The following streamTypes are supported (from the Docker stdcopy package):
//  - Stdin    - 0
//  - Stderr   - 1
//  - Stderr   - 2
//  - Systemerr - 3
// The function works by walking the supplied stream alternating between reading
// the prefix header bytes followed by the number of bytes specified in the prefix
// size bytes.
func FilterDockerStream(reader io.Reader, streamType int) ([]string, error) {
	logging.Stderr("[FilterDockerstream] filtering reader for steam type: %d", streamType)
	defer logging.Stderr("[FilterDockerstream] leaving")
	result := []string{}
	h, err := readFrameHeader(reader)
	if err != nil && err != io.EOF {
		logging.Stderr("[FilterDockerstream] stream contains no initial header: %s", err)
		return nil, err
	}
	for err == nil {
		payload, err := readFramePayload(reader, h)
		if err != nil {
			logging.Stderr("[FilterDockerstream] error reading payload: %s", err)
			return result, err
		}
		// Discard payload if streamType doesn't match requested
		if h.streamType == streamType {
			lines := bytes.Split(payload, []byte("\n"))
			for _, line := range lines {
				if len(line) != 0 {
					logging.Stderr("[FilterDockerstream] extracted line: %s", string(line))
					result = append(result, string(line))
				}
			}
		}
		h, err = readFrameHeader(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// readFrameHeader returns a Docker stream frameHeader from the supplied Reader
// See - https://docs.docker.com/engine/api/v1.29/#operation/ContainerAttach
// It returns io.EOF if the buffer is empty, any other errors returned should be
// handled by the caller
func readFrameHeader(reader io.Reader) (*frameHeader, error) {
	logging.Stderr("[readFrameHeader] reading header")
	prefix := make([]byte, stdWriterPrefixLen)

	_, err := reader.Read(prefix)
	if err != nil {
		logging.Stderr("[readFrameHeader] error thrown reading frameHeader: %s", err)
		return nil, err
	}

	header := frameHeader{
		streamType:  int(prefix[0]),
		payloadSize: int(binary.BigEndian.Uint32(prefix[stdWriterSizeIndex:])),
	}
	logging.Stderr("[readFrameHeader] extracted streamType: %d, payloadSize: %d", header.streamType, header.payloadSize)
	return &header, nil
}

// readFrameHeader returns a payload byte array from the Reader whose length is
// specified in the frameHeader
// It does not return io.EOF at the end of the buffer as it uses ioutil.ReadAll,
// any other errors returned should be handled by the caller
func readFramePayload(r io.Reader, h *frameHeader) ([]byte, error) {
	logging.Stderr("[readFramePayload] reading payloadSize: %d", h.payloadSize)
	lr := io.LimitReader(r, int64(h.payloadSize))
	payload, err := ioutil.ReadAll(lr)
	if err != nil {
		logging.Stderr("[readFramePayload] Error reading payload: %s", err)
		return nil, err
	}
	logging.Stderr("[readFramePayload] payload size: %d", len(payload))
	logging.Stderr("[readFramePayload] extracted payload: %s", payload)
	return payload, nil
}
