# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building
```bash
make build          # Build for current platform (creates ./gq binary)
make build-all      # Build for all platforms (outputs to build/ directory)
go build ./cmd/gq   # Direct go build
```

### Testing
```bash
make test           # Run all tests
go test -v ./...    # Run tests with verbose output
make demo           # Run demo with sample data
```

### Code Quality
```bash
make fmt            # Format code with go fmt
make lint           # Lint code (requires golangci-lint)
go mod tidy         # Clean up dependencies
```

### Development Workflow
```bash
make deps           # Install/update dependencies
make clean          # Clean build artifacts
```

## Architecture Overview

This is a command-line GPX (GPS Exchange Format) analysis tool that processes GPS track/route data and calculates elevation metrics per distance segments.

### Core Data Flow
1. **Input**: GPX file (XML) containing GPS coordinates with elevation data
2. **Parsing**: XML → structured Point data with lat/lon/elevation
3. **Distance Calculation**: Haversine formula for accurate geographic distances
4. **Segmentation**: Divide route into equal distance segments (default 1km)
5. **Analysis**: Calculate elevation gain/loss per segment
6. **Output**: JSON/CSV formatted results

### Package Architecture

#### `pkg/gpx/` - Core GPX Processing
- **`parser.go`**: XML parsing of GPX files
  - Handles both routes (`<rte>`) and tracks (`<trk>`) with segments
  - Normalizes RoutePoint and TrackPoint into unified Point structure
  - Supports optional elevation and time data
- **`analyzer.go`**: Elevation analysis engine
  - Segments routes by distance using cumulative distance calculation
  - Calculates elevation gain/loss per segment with interpolation for segment boundaries
  - Generates comprehensive metadata (totals, min/max elevations)

#### `pkg/geo/` - Geographic Calculations
- **`distance.go`**: Haversine distance formula implementation
  - High-precision great-circle distance calculations
  - Earth radius constant: 6371.0 km

#### `pkg/output/` - Output Formatting
- **`formatter.go`**: Multiple output format support
  - JSON (pretty and compact)
  - CSV with standard headers
  - Configurable precision for different data types

#### `cmd/gq/` - CLI Interface
- **`main.go`**: Command-line argument parsing and application orchestration
  - Supports stdin input for pipeline operations
  - Flag parsing for format and segment distance options

### Key Algorithms

#### Segment Division Strategy
The analyzer divides GPS tracks into equal distance segments (e.g., 1km each) regardless of point density. When segment boundaries don't align with actual GPS points, the system interpolates coordinate and elevation data to create precise segment start/end points.

#### Elevation Calculations
- **Gain**: Sum of positive elevation changes between consecutive points
- **Loss**: Sum of negative elevation changes (converted to positive values)
- **Net**: Gain minus loss for each segment

### Testing and Sample Data
- Sample GPX file: `sample/ikpht-long-new.gpx`
- Demo command runs analysis on sample data in CSV format
- CI pipeline tests with different segment distances (1km and 0.5km)

### Build Configuration
- **Target Go Version**: 1.25.1 (go.mod), CI uses 1.21
- **Main Entry Point**: `cmd/gq/main.go`
- **Binary Name**: `gq`
- **Cross-Platform Builds**: Linux, macOS (Intel/ARM), Windows