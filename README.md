# makecli

## it just works

if it doesnt it tells you exactly what doesnt work and we can fix it

## introduction
This package streamlines the process of adding Go-based CLI tools to your system. With minimal configuration, you can compile and add your Go programs to your system's PATH, making them easily accessible as command-line utilities. This tool is particularly useful for developers who frequently create and use custom CLI tools.

## Features
- **Easy CLI Tool Addition**: Simplify the process of adding new Go tools to your system.
- **API Integration**: Each CLI tool can access APIs, functioning as an independent service.
- **Automatic PATH Handling**: Automatically adds the tool to the user's PATH.
- **Customizable Installation**: Choose the binary name and target directory for your tools.

## Requirements
- Go programming language installed.
- git installed.
- Basic understanding of CLI operations and Go language.

## Installation

dont do this if you have no clue what youre doing:

```sh
curl https://github.com/m-c-frank/makecli | sh
```

## Usage

as long as your own go file does something you can use this to make it work everywhere

1. `makecli -name "name of your tool" -source your-go-file.go

2. If the destination is not specified, the tool defaults to `$HOME/tools`.

3. The script will compile your Go source code, creating an executable in the specified directory.

4. The tool's directory is automatically added to your PATH, making the tool available user-environment-wide.

## Example
Add a new tool named "mytool":
```sh
mkcli -name note -source main.go
```

## Notes
- Ensure Go is properly installed and configured on your system.
- Restart your terminal or source your profile (e.g., `source ~/.bashrc`) to apply PATH changes.

## Conclusion
This package offers an efficient way to integrate Go-based tools into your CLI environment. It's a versatile solution for developers looking to enhance their productivity and streamline tool management.
