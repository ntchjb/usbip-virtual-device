package protocol

import "io"

type Serializer interface {
	Decode(reader io.Reader) error
	Encode(writer io.Writer) error
}
