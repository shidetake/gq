# gq - GPX Query Tool

A command-line tool for analyzing GPX files, inspired by `jq`. Calculates elevation gain and loss per distance segment.

## Features

- Parse GPX files (routes and tracks)
- Calculate elevation gain/loss per configurable distance segments
- Multiple output formats: JSON, CSV, compact JSON
- Support for stdin input (pipeline-friendly)
- High-precision distance calculations using Haversine formula

## Installation

### Using Go (Recommended)

```bash
go install github.com/shidetake/gq@latest
```

### Download Binary

Download the latest binary from [GitHub Releases](https://github.com/shidetake/gq/releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/shidetake/gq/releases/latest/download/gq-darwin-arm64 -o gq
chmod +x gq
sudo mv gq /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/shidetake/gq/releases/latest/download/gq-darwin-amd64 -o gq
chmod +x gq
sudo mv gq /usr/local/bin/

# Linux (x64)
curl -L https://github.com/shidetake/gq/releases/latest/download/gq-linux-amd64 -o gq
chmod +x gq
sudo mv gq /usr/local/bin/

# Windows (x64)
# Download gq-windows-amd64.exe from releases page
```

### Build from Source

```bash
git clone https://github.com/shidetake/gq.git
cd gq
make build
```

## Usage

### Basic Examples

```bash
# Analyze with 1km segments (JSON output)
gq route.gpx

# Output in CSV format
gq --format csv route.gpx

# Use 500m segments
gq -d 0.5 --format csv route.gpx

# Read from stdin
cat route.gpx | gq --format csv
```

### Command Line Options

```
Usage: gq [options] [file]

Options:
  -f, --format FORMAT Output format: json, csv (default: json)
  -d, --distance NUM  Segment distance in km (default: 1.0)
  -h, --help          Show help
```

## Output Formats

### JSON Output (default)

```json
{
  "metadata": {
    "total_distance_km": 29.59,
    "total_points": 1342,
    "segment_distance_km": 1.0,
    "total_elevation_gain_m": 2881.4,
    "total_elevation_loss_m": -2881.0,
    "min_elevation_m": 171.1,
    "max_elevation_m": 1098.9
  },
  "segments": [
    {
      "segment": 1,
      "start_km": 0.0,
      "end_km": 1.0,
      "distance_km": 1.0,
      "elevation_gain_m": 70.7,
      "elevation_loss_m": -19.2,
      "net_elevation_m": 51.6,
      "start_elevation_m": 251.3,
      "end_elevation_m": 302.9,
      "point_count": 38
    }
  ]
}
```

### CSV Output

```csv
segment,start_km,end_km,distance_km,elevation_gain_m,elevation_loss_m,net_elevation_m,start_elevation_m,end_elevation_m,point_count
1,0.000,1.000,1.000,70.7,-19.2,51.6,251.3,302.9,38
2,1.000,2.000,1.000,271.3,-28.1,243.3,302.6,545.8,60
```

## Development

### Build

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make clean          # Clean build artifacts
```

### Test

```bash
make test           # Run tests
make demo           # Run demo with sample data
```

### Code Quality

```bash
make fmt            # Format code
make lint           # Lint code (requires golangci-lint)
```

## Algorithm

1. **Parse GPX**: Extract route points or track points with coordinates and elevation
2. **Calculate Distances**: Use Haversine formula for accurate geographic distance calculation
3. **Segment Division**: Divide route into equal distance segments (default 1km)
4. **Elevation Analysis**: Calculate elevation gain and loss within each segment
5. **Output**: Format results in JSON or CSV

## Supported GPX Elements

- Route points (`<rtept>`)
- Track points (`<trkpt>`)
- Elevation data (`<ele>`)
- Time data (`<time>`) - parsed but not used in analysis

## License

MIT