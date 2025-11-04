# cref - C Refresher

A lightweight hot-reload tool for C development that allows you to quickly recompile and run your code without leaving the terminal.

## Overview

cref is a simple command-line utility that watches for keyboard input and recompiles your C programs on demand. It's designed to speed up the development workflow by eliminating the need to manually recompile and run your code after every change.

## Features

- Instant recompilation with CTRL+R
- Automatic rerun prompt after program execution
- Clean temporary file handling
- Terminal-based interface with no external dependencies except clang
- Minimal output for distraction-free coding

## Installation

### From Source

Requirements:
- Go 1.21 or later
- clang compiler

```bash
git clone https://github.com/yourusername/cref.git
cd cref
go build -o cref
sudo mv cref /usr/local/bin/
```

### From Release

Download the latest binary for your platform from the [releases page](https://github.com/yourusername/cref/releases).

```bash
chmod +x cref-linux-amd64
sudo mv cref-linux-amd64 /usr/local/bin/cref
```

## Usage

Run your C file with cref:

```bash
cref run program.c
```

### Keyboard Shortcuts

- **CTRL+R** - Recompile and run the program
- **CTRL+C** - Exit cref
- **y** - Rerun the program (when prompted)
- **n** - Exit after program finishes (when prompted)

## Example

```bash
$ cref run hello.c
=== C Refresher ===
Press CTRL+R to recompile

Hello, World!
Rerun? (y/n): y
Hello, World!
Rerun? (y/n): n
```

Edit your source file, press CTRL+R in the terminal, and see your changes immediately.

## How It Works

cref compiles your C code using clang and stores the executable in a temporary directory. When you press CTRL+R, it recompiles the source file and runs the new version. All temporary files are automatically cleaned up when you exit.

## Requirements

- clang compiler must be installed and available in PATH
- Linux operating system (ARM64 and x86_64 supported)
- Terminal with raw mode support

## License

MIT License - See LICENSE file for details

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
