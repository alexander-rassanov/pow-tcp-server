package pow

import (
	"alexander.rassanov/pow-tcp-server/pkg/protocol"
	"io"
)

// GetServiceByStream requests a quote of wisdom by a protocol using stream.
func GetServiceByStream(stream io.ReadWriter) (interface{}, error) {
	err := protocol.SendPackage(protocol.NewMessage(RequestChallenge, nil), stream)
	if err != nil {
		return nil, err
	}
	m, err := protocol.ReadPackage(stream)
	if err != nil {
		return nil, err
	}
	if m.Header != ResponseChallenge {
		return nil, ErrBadResponse
	}
	hd, ok := m.Payload.(HashCashData)
	if !ok {
		return nil, protocol.ErrBadPayload
	}
	resolvedHd, err := hd.Resolve()
	if err != nil {
		return nil, err
	}
	err = protocol.SendPackage(protocol.NewMessage(RequestService, resolvedHd), stream)
	if err != nil {
		return nil, err
	}
	m, err = protocol.ReadPackage(stream)
	if err != nil {
		return m, err
	}
	if m.Header != ResponseService {
		return m, ErrBadResponse
	}
	return m.Payload, nil
}
