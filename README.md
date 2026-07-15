# Preview CLI

Preview CLI is a modular command-line tool designed to fetch, list, and play video content from various sources. It uses a plugin-like architecture, allowing for easy extension to support new content providers.

## Features

- **Modular Design**: Easily add new content providers by implementing core interfaces.
- **Flexible Playback**: Supports piping content directly to your preferred media player.
- **Configurable**: Customize proxy settings and player arguments via `~/.preview.toml`.

## Installation

1. Ensure you have Go installed.
2. Clone the repository.
3. Build the project:
   ```bash
   go build -o preview cmd/main.go
   ```

## Usage

The CLI provides two main commands: `show` and `play`.

### Show Content
Lists available content from a provider:
