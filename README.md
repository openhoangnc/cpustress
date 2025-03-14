# CPU Stress Test Tool (AKA The Fan Symphony Conductor)

A lightweight command-line utility written in Go that turns your peaceful computer into a jet engine simulator. Perfect for those cold winter days when you need an impromptu space heater or when you simply miss the sound of airplane takeoffs.

## Features

- Transform your quiet computer into a roaring beast
- Make your laptop hover slightly above your desk
- Check if your cooling system actually works (spoiler: you'll find out quickly)
- Graceful shutdown with Ctrl+C (in case your room gets too toasty)
- Cross-platform compatible: heat up any macOS, Linux, or Windows machine

## Installation

### Prerequisites

- Go 1.18 or higher
- A working fan (you'll need it)
- Fire extinguisher (just kidding... mostly)

### Option 1: Run directly from GitHub

```bash
go run github.com/openhoangnc/cpustress@latest
```

### Option 2: Clone and run locally

```bash
git clone https://github.com/openhoangnc/cpustress.git
cd cpustress
go run .
```

### Option 3: Install using Go

```bash
go install github.com/openhoangnc/cpustress@latest
```

Then run:

```bash
cpustress
```

## Usage

### Basic usage

Launch with default settings and listen to the beautiful sound of fans reaching for the stars:

```bash
cpustress
```

### Command line options

- `-w`: Number of worker goroutines (default: all cores, because why not use everything you paid for?)
- `-t`: Duration in minutes (default: 0, runs until interrupted or your laptop achieves liftoff)

### Examples

Run with 4 worker goroutines (for a gentle breeze):

```bash
cpustress -w 4
```

Run for 5 minutes (perfect for warming up leftover coffee):

```bash
cpustress -t 5
```

Run with 8 worker goroutines for 10 minutes (recommended if you're trying to simulate a desktop space heater):

```bash
cpustress -w 8 -t 10
```

## Cross-Platform CPU Usage Monitoring

How we check if your machine is properly suffering:

- macOS: Uses the built-in `top` command (we could have used "How hot does this MacBook feel to touch?" but that's less scientific)
- Linux: Uses the built-in `top` command with batch mode (penguin-approved method)
- Windows: Uses `wmic` to query CPU load percentage (because Task Manager is too easy)

## Important Notes

- High CPU usage will increase your system's temperature. Your laptop may become suitable for frying eggs
- Running this on battery is like trying to drain a swimming pool with a coffee mug – fast and effective
- If your computer starts making new, exciting noises, that's just its way of asking for a break
- Perfect for testing if your thermal paste application was actually "good enough"

## License

MIT (Making Intense Temperatures)