# TarSplitter

## Usage

```
Split or join tar (or tar.gz) files on file boundaries

Usage:
  tarsplitter [command]

Available Commands:
  help        Help about any command
  join        Join multiple tar or tar.gz files into one
  split       Split tar or tar.gz files on file boundaries

Flags:
  -h, --help   help for tarsplitter

Use "tarsplitter [command] --help" for more information about a command.
```

```
Split tar or tar.gz files on file boundaries

Usage:
  tarsplitter split <tar file> <destination direction> [flags]

Flags:
  -h, --help             help for split
  -s, --split-size int   max size of split in MiB (default 1000)
```

```
Join multiple tar or tar.gz files into one

Usage:
  tarsplitter join <tar files ...> [flags]

Flags:
  -h, --help   help for join
```

## Build

```bash
make
```

## Install

```bash
make install
```
