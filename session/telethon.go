package session

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/session"
)

func ImportTelethonSession(ctx context.Context, base64Session string, storage session.Storage) error {
	data, err := session.TelethonSession(base64Session)
	if err != nil {
		return errors.Wrap(err, "parse telethon session")
	}

	loader := session.Loader{Storage: storage}
	if err := loader.Save(ctx, data); err != nil {
		return errors.Wrap(err, "store session")
	}

	return nil
}
