package challenge_response

import (
	"alexander.rassanov/pow-tcp-server/pkg/pow"
	"alexander.rassanov/pow-tcp-server/pkg/protocol"
	"io"
)

// GetServiceByStream requests a quote of wisdom by a protocol using stream
func GetServiceByStream(stream io.ReadWriter) (interface{}, error) {
	err := protocol.SendPackage(protocol.NewMessage(RequestChallenge, ""), stream)
	if err != nil {
		return "", err
	}
	m, err := protocol.ReadPackage(stream)
	if err != nil {
		return "", err
	}
	hd, ok := m.Payload.(pow.HashCashData)
	if !ok {
		return "", protocol.ErrBadPayload
	}
	resolvedHd, err := hd.Resolve()
	if err != nil {
		return "", err
	}
	err = protocol.SendPackage(protocol.NewMessage(RequestService, resolvedHd), stream)
	if err != nil {
		return "", err
	}
	m, err = protocol.ReadPackage(stream)
	if err != nil {
		return "", err
	}
	quote, ok := m.Payload.(string)
	if !ok {
		return "", protocol.ErrBadPayload
	}
	return quote, nil
}
