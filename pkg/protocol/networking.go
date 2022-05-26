package protocol

import (
	"bufio"
	"io"
)

// delimiter is a symbol which indicates the end of a package.
const delimiter = '\n'

// SendPackage sends a proper generator network package.
// Network package is a message with appended delimiter.
func SendPackage(message Message, conn io.Writer) error {
	encodedMessage, err := message.Encode()
	if err != nil {
		return err
	}
	ret := append(make([]byte, 0), encodedMessage...)
	ret = append(ret, delimiter)
	_, err = conn.Write(ret)
	return err
}

// ReadPackage reads a network package.
func ReadPackage(reader io.Reader) (Message, error) {
	b, err := bufio.NewReader(reader).ReadBytes(delimiter)
	if err != nil {
		return Message{}, err
	}
	return ParseMessage(b)
}
