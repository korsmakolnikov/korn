# kornvimgen

A modern CLI tool written in Go to easily manage and run multiple Neovim
configurations.

## About
`kornvimgen` is a command-line interface designed to streamline your Neovim
workflow. If you're a developer, a system administrator, or a hobbyist who
frequently switches between different projects and requires different Neovim
setups or that needs a different home/work configuration, this tool is for you.
With `kornvimgen`, you can create, manage, and run isolated Neovim
configurations. This prevents conflicts between plugins and settings, making it
easy to test new setups or switch environments without messing up your main
configuration.
`kornvimgen` also push the user to separate the different concerns of its
configuration. Instead of maintaining a large Neovim configuration, having
separated configurations for each need will make emerge what they have in
common and push you to isolate that feature in a separated plugin.
Another benefit of `kornvimgen` is installing all the plugins in the build
directory. If this duplicates the plugin codes between configurations, it will
make also easier to debug them and to experiment with your own plugins in
development.

## Features

- **Isolated Builds**: Create "builds" with their own `lazy.nvim` directories and configurations, preventing conflicts.
- **Easy Management**: Use simple commands to create, run, and delete your builds.
- **Customizable**: Set up custom configurations and plugin lists for any project.

## Installation

### Prerequisites

- [Go](https://go.dev/doc/install) (version 1.18 or later)
- [Neovim](https://neovim.io/doc/user/starting.html#nvim-install) (version 0.9 or later)

### From Source

1. Clone the repository:
   ```bash
   git clone [https://github.com/korsmakolnikov/kornvimgen.git](https://github.com/korsmakolnikov/kornvimgen.git)
2. Build the executable
   ```bash
   go build
   ```
