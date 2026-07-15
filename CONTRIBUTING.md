# Contributing to Preview CLI

This guide is intended for developers and AI agents looking to extend the project by adding new content providers.

## Adding a New Provider

To add a new content source, follow these steps:

1. **Create a Package**: Create a new directory (e.g., `newsite/`) and implement the `Preview` interface defined in `preview.go`.
2. **Implement Logic**:
   - Ensure your struct has `SetSource`, `Name`, and `Get` methods.
   - If the site has multiple sub-domains, implement the `Domain` interface.
   - If the site supports search, implement the `Search` interface.
3. **Register**: Open `barrel/barrel.go` and add your new provider instance to the `Domain` slice.
4. **Test**: Run the CLI to ensure your provider appears in the list and functions correctly.

## Core Interfaces

Your provider must satisfy one of the following interfaces defined in `preview.go`:

- **`Preview`**: The base interface.
    - `SetSource(source string)`: Configures the specific endpoint or sub-domain.
    - `Name() string`: Returns the provider's name.
    - `Get(index int) ([]ContentResource, error)`: Fetches the content list.
- **`Domain`**: Extends `Preview` for providers with multiple sub-domains.
- **`Search`**: Extends `Preview` for providers that support search functionality.

## Coding Standards

- **HTTP Requests**: Use the `resty` library for network requests.
- **Registry**: All new providers must be instantiated and added to the `barrel.Domain` slice in `barrel/barrel.go`.
- **Utilities**: Use the `common` package for shared logic like M3U formatting or player execution.
- **Configuration**: Do not hardcode proxy or player settings; rely on the `config` package.
