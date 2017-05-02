package helper

import "encoding/binary"

// AddDockerStreamHeader prefixes the supplied byte slice with a docker stream header as seen
// in Dockers stdWriter.Write(p []byte).
//
// This takes the following form:
//
//  header := [8]byte{STREAM_TYPE, 0, 0, 0, SIZE1, SIZE2, SIZE3, SIZE4}
//
// `STREAM_TYPE` can be:
//   - 0: stdin
//   - 1: stdout
//   - 2: stderr
//   - 3: Systemer
//
// `SIZE1, SIZE2, SIZE3, SIZE4` are the 4 bytes of the uint32 size encoded as big endian.
func AddDockerStreamHeader(payload []byte, streamType int) []byte {
	const (
		prefixLen = 8
		sizeIndex = 4
		fdIndex   = 0
	)
	header := [prefixLen]byte{fdIndex: byte(streamType)}
	binary.BigEndian.PutUint32(header[sizeIndex:], uint32(len(payload)))
	return append(header[:], payload...)
}

// AddCustomDockerStreamHeader behaves indentically to addHeader but allows you to alter
// the calculated payload size in the header by supplying a sizeOffset to
// increase/decrease the size stored in the header.
func AddCustomDockerStreamHeader(payload []byte, streamType int, sizeOffset int) []byte {
	const (
		prefixLen = 8
		sizeIndex = 4
		fdIndex   = 0
	)
	header := [prefixLen]byte{fdIndex: byte(streamType)}
	binary.BigEndian.PutUint32(header[sizeIndex:], uint32(len(payload)+sizeOffset))
	return append(header[:], payload...)
}
