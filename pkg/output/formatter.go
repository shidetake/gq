package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/shidetake/gq/pkg/gpx"
)

type Format int

const (
	FormatJSON Format = iota
	FormatCSV
)

func FormatResult(result *gpx.AnalysisResult, format Format, writer io.Writer) error {
	switch format {
	case FormatJSON:
		return formatJSON(result, writer)
	case FormatCSV:
		return formatCSV(result, writer)
	default:
		return fmt.Errorf("unsupported format")
	}
}

func formatJSON(result *gpx.AnalysisResult, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func formatCSV(result *gpx.AnalysisResult, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{
		"segment",
		"start_km",
		"end_km",
		"distance_km",
		"elevation_gain_m",
		"elevation_loss_m",
		"net_elevation_m",
		"start_elevation_m",
		"end_elevation_m",
		"point_count",
	}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, segment := range result.Segments {
		row := []string{
			strconv.Itoa(segment.Number),
			fmt.Sprintf("%.3f", segment.StartKm),
			fmt.Sprintf("%.3f", segment.EndKm),
			fmt.Sprintf("%.3f", segment.DistanceKm),
			fmt.Sprintf("%.1f", segment.ElevationGainM),
			fmt.Sprintf("%.1f", segment.ElevationLossM),
			fmt.Sprintf("%.1f", segment.NetElevationM),
			fmt.Sprintf("%.1f", segment.StartElevationM),
			fmt.Sprintf("%.1f", segment.EndElevationM),
			strconv.Itoa(segment.PointCount),
		}
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

func FormatMetadata(metadata *gpx.Metadata, writer io.Writer) error {
	fmt.Fprintf(writer, "Total Distance: %.2f km\n", metadata.TotalDistanceKm)
	fmt.Fprintf(writer, "Total Points: %d\n", metadata.TotalPoints)
	fmt.Fprintf(writer, "Segment Distance: %.1f km\n", metadata.SegmentDistanceKm)
	fmt.Fprintf(writer, "Total Elevation Gain: %.1f m\n", metadata.TotalElevationGainM)
	fmt.Fprintf(writer, "Total Elevation Loss: %.1f m\n", metadata.TotalElevationLossM)
	fmt.Fprintf(writer, "Elevation Range: %.1f - %.1f m\n", metadata.MinElevationM, metadata.MaxElevationM)
	return nil
}