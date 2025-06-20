# yt-dlp-web

A simple web interface for [yt-dlp](https://github.com/yt-dlp/yt-dlp), built in vanilla Go

## Pre-requisites

- [yt-dlp](https://github.com/yt-dlp/yt-dlp) must be installed and available in your system's PATH.
- [Go](https://golang.org/dl/) (if building from source)

## Download and Run

### Windows

Download the latest release for 64-bit AMD64 Windows by visiting the [releases](https://github.com/jere-mie/yt-dlp-web/releases) page, or by running the following command in your shell:

```batch
irm -Uri https://github.com/jere-mie/yt-dlp-web/releases/latest/download/yt_dlp_web_windows_amd64.exe -O yt-dlp-web.exe
```

### Linux

Download the latest release for 64-bit AMD64 Linux from the [releases](https://github.com/jere-mie/yt-dlp-web/releases) page, or run the following command in your terminal:

```sh
curl -L https://github.com/jere-mie/yt-dlp-web/releases/latest/download/yt_dlp_web_linux_amd64 -o yt-dlp-web
chmod +x yt-dlp-web
```

## Configuration

yt-dlp-web can be configured using environment variables or a `.env` file in the project directory. This allows you to customize settings such as the server port, host, and admin password.

### Using Environment Variables

You can set environment variables directly in your shell before starting the application. For example, to change the port and download directory:

**Windows (Command Prompt):**
```batch
set PORT=8081
set PASSWORD=password
set HOST=127.0.0.1
yt-dlp-web.exe
```

**Linux/macOS:**
```sh
export PORT=8081
export PASSWORD=password
export HOST=127.0.0.1
./yt-dlp-web
```

### Using a `.env` File

Alternatively, you can create a `.env` file in the same directory as the executable. Each line should contain a key-value pair, for example:

```
PORT=8081
PASSWORD=password
HOST=127.0.0.1
```

When you start yt-dlp-web, it will automatically load configuration from this file if present.

### example.env

An `example.env` file is included in the repository to demonstrate the available configuration options. You can copy this file to `.env` and modify the values as needed:

```
cp example.env .env
```

Then edit `.env` to suit your environment.
