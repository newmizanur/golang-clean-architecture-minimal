## Setup

Install Taskfile CLI:
```
  go install github.com/go-task/task/v3/cmd/task@latest
```

Start PostgreSQL (Docker):
```
  task postgres:docker
```

Install Goose tool:
```
  task tools:install
```

Run migrations:
```
  task goose:up
```

Run the API:
```
  task run
```

Build:
```
  task build
```

Build for Amazon Linux (static):
```
  task build:amazon-linux
```

## Config

Update `config.json` for database and JWT settings.
