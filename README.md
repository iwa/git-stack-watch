# git-stack-watch

Periodically push changes of my docker compose stacks.

## Usage

Options:
```
  --repo /path/to/repo
        Path to the git repository to watch (required)
  --push
        Push changes after committing
```

### Docker Compose

```yaml
services:
  git-stack-watch:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: git-stack-watch
    restart: unless-stopped
    volumes:
      - /path/to/repo:/repo:rw
      - /path/to/.gitconfig:/root/.gitconfig:ro
      - /path/to/.ssh:/root/.ssh:ro
    environment:
      - TZ=Europe/Paris
    command: ["--repo", "/repo"]
```

> Note that due to Docker I/O latency on binded volumes, creation of commits can take a bit of time.
> This behavior is expected and certainly might not be fixed.

### Binary

```
go run main.go --repo /path/to/repo
```
