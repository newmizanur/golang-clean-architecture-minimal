## Setup

Install Taskfile CLI:
```
  go install github.com/go-task/task/v3/cmd/task@latest
```

Start MySQL (Docker):
```
  task mysql:docker
```

Install Goose + Jet tools:
```
  task tools:install
```

Run migrations:
```
  task goose:up
```

Generate Jet models:
```
  task jet:generate
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
