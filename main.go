package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// CPUInfo stores CPU usage data
type CPUInfo struct {
	Usage float64
}

func main() {
	// Parse command line flags
	numWorkers := flag.Int("w", runtime.NumCPU(), "number of worker goroutines")
	duration := flag.Int("t", 0, "duration in minutes (0 means run until interrupted)")
	flag.Parse()

	fmt.Printf("Starting CPU stress test with %d workers\n", *numWorkers)
	if *duration > 0 {
		fmt.Printf("Test will run for %d minutes\n", *duration)
	} else {
		fmt.Println("Test will run until interrupted (press Ctrl+C to stop)")
	}

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// If duration is specified, set up a timer to cancel the context
	if *duration > 0 {
		go func() {
			timer := time.NewTimer(time.Duration(*duration) * time.Minute)
			select {
			case <-timer.C:
				fmt.Println("\nTest duration completed")
				cancel()
			case <-ctx.Done():
				timer.Stop()
			}
		}()
	}

	// Go routine to handle interrupt signal
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Shutting down...")
		cancel()
	}()

	// Launch monitoring routine
	go monitorCPU(ctx)

	// Launch worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i)
	}

	// Wait for all workers to complete
	wg.Wait()
	fmt.Println("CPU stress test completed")
}

// worker performs CPU-intensive calculations to stress the CPU
func worker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	fmt.Printf("Worker %d started on CPU core\n", id)

	// Random number generator for each worker
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d shutting down\n", id)
			return
		default:
			// Perform CPU-intensive calculations
			for i := 0; i < 1000000; i++ {
				// Complex math operations to stress CPU
				_ = rng.Float64() * rng.Float64() / (rng.Float64() + 0.1)
				_ = rng.Int63n(1000000) * rng.Int63n(1000000)
			}
		}
	}
}

// monitorCPU periodically checks and displays CPU information
func monitorCPU(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			info, err := getCPUInfo()
			if err != nil {
				fmt.Printf("\rError getting CPU info: %v                      ", err)
				continue
			}

			// Clear line and print new info
			fmt.Printf("\rCPU Usage: %.1f%%          ", info.Usage)
		}
	}
}

// getCPUInfo retrieves CPU usage
func getCPUInfo() (CPUInfo, error) {
	var info CPUInfo

	// Get CPU usage
	var err error
	info.Usage, err = getCPUUsage()
	if err != nil {
		return info, fmt.Errorf("failed to get CPU usage: %v", err)
	}

	return info, nil
}

// getCPUUsage retrieves the CPU usage percentage
func getCPUUsage() (float64, error) {
	switch runtime.GOOS {
	case "darwin":
		return getMacOSCPUUsage()
	case "linux":
		return getLinuxCPUUsage()
	case "windows":
		return getWindowsCPUUsage()
	default:
		return 0, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// getMacOSCPUUsage gets CPU usage on macOS
func getMacOSCPUUsage() (float64, error) {
	cmd := exec.Command("top", "-l", "1", "-n", "0", "-stats", "cpu")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("top command failed: %v", err)
	}

	// Parse the output to get CPU usage
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "CPU usage") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				userStr := strings.TrimSuffix(fields[2], "%")
				sysStr := strings.TrimSuffix(fields[4], "%")
				userUsage, _ := strconv.ParseFloat(userStr, 64)
				sysUsage, _ := strconv.ParseFloat(sysStr, 64)
				return userUsage + sysUsage, nil
			}
		}
	}
	return 0, fmt.Errorf("failed to parse CPU usage from top output")
}

// getLinuxCPUUsage gets CPU usage on Linux
func getLinuxCPUUsage() (float64, error) {
	cmd := exec.Command("top", "-bn", "1")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("top command failed: %v", err)
	}

	// Parse the output to get CPU usage
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "%Cpu(s)") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				usageStr := fields[1]
				usage, err := strconv.ParseFloat(usageStr, 64)
				if err == nil {
					return usage, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("failed to parse CPU usage from top output")
}

// getWindowsCPUUsage gets CPU usage on Windows
func getWindowsCPUUsage() (float64, error) {
	cmd := exec.Command("wmic", "cpu", "get", "LoadPercentage")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("wmic command failed: %v", err)
	}

	// Parse the output to get CPU usage
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("unexpected wmic output format")
	}

	// The usage is typically in the second line
	usageStr := strings.TrimSpace(lines[1])
	if usageStr == "" && len(lines) > 2 {
		usageStr = strings.TrimSpace(lines[2])
	}
	usage, err := strconv.ParseFloat(usageStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CPU usage: %v", err)
	}
	return usage, nil
}
