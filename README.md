# PomoLite

<div align="center">

```
██████╗  ██████╗ ███╗   ███╗ ██████╗ ██╗     ██╗████████╗███████╗
██╔══██╗██╔═══██╗████╗ ████║██╔═══██╗██║     ██║╚══██╔══╝██╔════╝
██████╔╝██║   ██║██╔████╔██║██║   ██║██║     ██║   ██║   █████╗  
██╔═══╝ ██║   ██║██║╚██╔╝██║██║   ██║██║     ██║   ██║   ██╔══╝  
██║     ╚██████╔╝██║ ╚═╝ ██║╚██████╔╝███████╗██║   ██║   ███████╗
╚═╝      ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚══════╝
```

**A lightweight CLI Pomodoro application designed for students and professionals to efficiently accomplish tasks and maximize their learning potential.**

</div>

---

## Features

- **Simple Pomodoro Timer**: Start work and break sessions right from your terminal.
- **Session Tracking**: Automatically saves your completed sessions to a local SQLite database.
- **Productivity Stats**: View detailed statistics of your work sessions filtered by timeframes (today, week, month, year, or all-time).
- **Interactive Controls**: Pause, resume, or quit the timer using keyboard shortcuts.
- **Customizable Sessions**: Set custom durations for work and break periods and add labels to your sessions.
- **Desktop Notifications**: Get notified when a session or break is complete.

---

## Installation

### Using the Installation Script

For a seamless installation, you can use the provided installation script. This will build the binary and attempt to install it.

```sh
chmod +x scripts/install.sh
./scripts/install.sh
```

### Using `go install`

You can install PomoLite using `go install`. This will build the binary and place it in your `$GOPATH/bin` directory.

```sh
go install github.com/Dima-salang/pomolite/cmd/pomo@latest
```
The executable will be named `pomo`.

### Building from Source

Alternatively, you can clone the repository and build the application from source.

```sh
git clone https://github.com/Dima-salang/pomolite.git
cd pomolite
go build -o pomolite cmd/pomo/main.go
```

---

## Usage

PomoLite provides three main commands: `start`, `sessions`, and `stat`.

### `start`

Starts a new Pomodoro timer.

```sh
pomo start [flags]
```

**Flags:**
- `-m`, `--minutes`: The duration of the work session in minutes (default: 30).
- `-b`, `--break`: The duration of the break in minutes (default: 5).
- `-l`, `--label`: A descriptive label for the work session (default: "Work").

**Example:**
```sh
# Start a 25-minute timer with a 5-minute break and the label "Coding"
pomo start -m 25 -b 5 -l "Coding"
```

**Interactive Controls:**
While the timer is running, you can use the following keys:
- `p`: Pause the timer.
- `r`: Resume the timer.
- `q`: Quit the timer and save the session progress.

### `sessions`

Lists your past Pomodoro sessions.

```sh
pomo sessions [flags]
```

**Flags:**
- `-l`, `--limit`: The number of recent sessions to display. If not specified, all sessions are shown.

**Example:**
```sh
# List the last 10 sessions
pomo sessions -l 10
```

### `stat`

Displays statistics about your Pomodoro sessions.

```sh
pomo stat [flags]
```

**Flags:**
- `-t`, `--timeframe`: The timeframe for the statistics. Possible values are `all`, `today`, `week`, `month`, `year` (default: "all").

**Example:**
```sh
# Show statistics for the current week
pomo stat -t week
```

---

## License

This project is licensed under the terms of the LICENSE file.

---

Developed by **PUTAN LUIS GABRIELLE** <luisgabrielle1026@gmail.com>