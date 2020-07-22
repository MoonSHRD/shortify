# Shortify

## Requirements

- Go 1.14.4
- MongoDB

## Building

```bash
go build ./cmd/shortify
```

## Running

1. Firstly, you need to generate config:
    ```bash
    ./shortify -genconfig > backend.toml
    ```

2. Then just edit it, and start executable with specified config:
    ```bash
    ./shortify -config ./backend.toml
    ```