package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/shidetake/gq/pkg/gpx"
	"github.com/shidetake/gq/pkg/output"
)

const (
	DefaultSegmentDistance = 1.0
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]

	if len(args) == 0 {
		return fmt.Errorf("usage: gq [options] [file]\n\nOptions:\n  -f, --format FORMAT Output format: json, csv (default: json)\n  -d, --distance NUM  Segment distance in km (default: %.1f)\n  -h, --help          Show help", DefaultSegmentDistance)
	}

	// Parse arguments
	var (
		filename        string
		format          = output.FormatJSON
		segmentDistance = DefaultSegmentDistance
		showHelp        = false
	)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			showHelp = true
		case "-c", "--csv":
			// Deprecated: use --format csv instead
			fmt.Fprintf(os.Stderr, "Warning: -c/--csv is deprecated, use --format csv instead\n")
			format = output.FormatCSV
		case "-f", "--format":
			if i+1 < len(args) {
				switch args[i+1] {
				case "json":
					format = output.FormatJSON
				case "csv":
					format = output.FormatCSV
				default:
					return fmt.Errorf("unsupported format: %s (supported: json, csv)", args[i+1])
				}
				i++ // Skip the next argument since it's the format value
			} else {
				return fmt.Errorf("format flag requires a value")
			}
		case "-d", "--distance":
			if i+1 < len(args) {
				if dist, err := strconv.ParseFloat(args[i+1], 64); err == nil && dist > 0 {
					segmentDistance = dist
					i++ // Skip the next argument since it's the distance value
				} else {
					return fmt.Errorf("invalid distance value: %s", args[i+1])
				}
			} else {
				return fmt.Errorf("distance flag requires a value")
			}
		default:
			if len(arg) == 0 || arg[0] != '-' {
				if filename == "" {
					filename = arg
				}
			}
		}
	}

	if showHelp {
		fmt.Printf("gq - GPX Query Tool\n\n")
		fmt.Printf("Usage: gq [options] [file]\n\n")
		fmt.Printf("Options:\n")
		fmt.Printf("  -f, --format FORMAT Output format: json, csv (default: json)\n")
		fmt.Printf("  -d, --distance NUM  Segment distance in km (default: %.1f)\n", DefaultSegmentDistance)
		fmt.Printf("  -h, --help          Show help\n\n")
		fmt.Printf("Examples:\n")
		fmt.Printf("  gq route.gpx                    # Analyze with 1km segments (JSON)\n")
		fmt.Printf("  gq --format csv route.gpx       # Output in CSV format\n")
		fmt.Printf("  gq -f csv route.gpx             # Output in CSV format (short)\n")
		fmt.Printf("  gq -d 0.5 -f csv route.gpx     # 500m segments in CSV\n")
		fmt.Printf("  cat route.gpx | gq -f csv       # Read from stdin\n")
		return nil
	}

	// Determine input source
	var reader io.Reader
	if filename == "" || filename == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", filename, err)
		}
		defer file.Close()
		reader = file
	}

	// Parse GPX
	gpxData, err := gpx.ParseGPX(reader)
	if err != nil {
		return fmt.Errorf("failed to parse GPX: %w", err)
	}

	// Extract points
	points := gpxData.GetAllPoints()
	if len(points) == 0 {
		return fmt.Errorf("no points found in GPX file")
	}

	// Analyze elevation
	result := gpx.AnalyzeElevation(points, segmentDistance)

	// Output result
	if err := output.FormatResult(result, format, os.Stdout); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}