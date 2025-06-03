# gotd-session-storage

## Overview
`gotd-session` is a lightweight Go library designed to facilitate session management for the gotd Telegram client. It offers flexible session storage solutions, including SQLite and in-memory options, and supports importing and converting session data from popular Telegram clients like Telethon and Pyrogram.

## Features
- Session storage implementations for gotd client.
- Support for SQLite with configurable journaling mode.
- In-memory session storage for testing or ephemeral usage.
- Import and conversion utilities for Telethon and Pyrogram sessions.
- Simple, concurrent-safe design.

## Installation

```bash
go get github.com/pageton/gotd-session-storage
```

## Usage

```go
import (
    "context"
    "log"
    "github.com/pageton/gotd-session-storage/storage"
    "github.com/pageton/gotd-session-storage/session"
)

func main() {
    ctx := context.Background()

    sessionStorage, err := storage.NewSQLiteSessionStorage("session.db")
    if err != nil {
        log.Fatal(err)
    }

    err = session.ImportTelethonSession(ctx, "<TELETHON_SESSION_STRING>", sessionStorage)
    if err != nil {
        log.Fatal(err)
    }

    // Use sessionStorage with gotd client...
}
```

## Supported Session Formats

- Telethon string sessions.
- Pyrogram string sessions.

## Contributing
Contributions and suggestions are welcome! Please open an issue or submit a pull request.
