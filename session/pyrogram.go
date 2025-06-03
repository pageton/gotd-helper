package session

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/gotd/td/crypto"

	"github.com/gotd/td/session"
)

func getDCAddress(dc int) (string, error) {
	lookup := map[int]string{
		1: "149.154.175.50",
		2: "149.154.167.51",
		3: "149.154.175.100",
		4: "149.154.167.91",
		5: "91.108.56.130",
	}

	addr, exists := lookup[dc]
	if !exists {
		return "", fmt.Errorf("unknown dc id: %d", dc)
	}
	return addr, nil
}

func PyrogramSession(sessionStr string) (*session.Data, error) {
	const (
		dcLen     = 1
		apiIDLen  = 4
		flagLen   = 1
		keyLen    = 256
		userIDLen = 8
		isBotLen  = 1
	)

	for len(sessionStr)%4 != 0 {
		sessionStr += "="
	}

	rawData, err := base64.URLEncoding.DecodeString(sessionStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode base64")
	}

	expectedLen := dcLen + apiIDLen + flagLen + keyLen + userIDLen + isBotLen
	if len(rawData) != expectedLen {
		return nil, errors.Errorf("length mismatch: got %d expected %d", len(rawData), expectedLen)
	}

	dc := int(rawData[0])
	ip, err := getDCAddress(dc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dc ip")
	}

	const defaultPort = 443
	keyStart := dcLen + apiIDLen + flagLen
	var authKey crypto.Key
	copy(authKey[:], rawData[keyStart:keyStart+keyLen])
	keyID := authKey.WithID().ID

	return &session.Data{
		DC:        dc,
		Addr:      net.JoinHostPort(ip, strconv.Itoa(defaultPort)),
		AuthKey:   authKey[:],
		AuthKeyID: keyID[:],
	}, nil
}

func ImportPyrogramSession(ctx context.Context, base64Session string, storage session.Storage) error {
	data, err := PyrogramSession(base64Session)
	if err != nil {
		return errors.Wrap(err, "parse pyrogram session")
	}

	loader := session.Loader{Storage: storage}
	if err := loader.Save(ctx, data); err != nil {
		return errors.Wrap(err, "store session")
	}

	return nil
}
