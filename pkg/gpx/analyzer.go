package gpx

import (
	"github.com/shidetake/gq/pkg/geo"
)

type Segment struct {
	Number           int     `json:"segment"`
	StartKm          float64 `json:"start_km"`
	EndKm            float64 `json:"end_km"`
	DistanceKm       float64 `json:"distance_km"`
	ElevationGainM   float64 `json:"elevation_gain_m"`
	ElevationLossM   float64 `json:"elevation_loss_m"`
	NetElevationM    float64 `json:"net_elevation_m"`
	StartElevationM  float64 `json:"start_elevation_m"`
	EndElevationM    float64 `json:"end_elevation_m"`
	PointCount       int     `json:"point_count"`
}

type AnalysisResult struct {
	Metadata Metadata  `json:"metadata"`
	Segments []Segment `json:"segments"`
}

type Metadata struct {
	TotalDistanceKm   float64 `json:"total_distance_km"`
	TotalPoints       int     `json:"total_points"`
	SegmentDistanceKm float64 `json:"segment_distance_km"`
	TotalElevationGainM float64 `json:"total_elevation_gain_m"`
	TotalElevationLossM float64 `json:"total_elevation_loss_m"`
	MinElevationM     float64 `json:"min_elevation_m"`
	MaxElevationM     float64 `json:"max_elevation_m"`
}

func AnalyzeElevation(points []Point, segmentDistanceKm float64) *AnalysisResult {
	if len(points) < 2 {
		return &AnalysisResult{
			Metadata: Metadata{SegmentDistanceKm: segmentDistanceKm},
			Segments: []Segment{},
		}
	}

	// Calculate cumulative distances
	distances := make([]float64, len(points))
	distances[0] = 0

	for i := 1; i < len(points); i++ {
		dist := geo.HaversineDistance(
			points[i-1].Lat, points[i-1].Lon,
			points[i].Lat, points[i].Lon,
		)
		distances[i] = distances[i-1] + dist
	}

	totalDistance := distances[len(distances)-1]

	// Calculate metadata
	metadata := calculateMetadata(points, totalDistance, segmentDistanceKm)

	// Divide into segments
	segments := divideIntoSegments(points, distances, segmentDistanceKm)

	return &AnalysisResult{
		Metadata: metadata,
		Segments: segments,
	}
}

func calculateMetadata(points []Point, totalDistance, segmentDistance float64) Metadata {
	if len(points) == 0 {
		return Metadata{
			TotalDistanceKm:   totalDistance,
			TotalPoints:       len(points),
			SegmentDistanceKm: segmentDistance,
		}
	}

	var totalGain, totalLoss float64
	minElev, maxElev := points[0].Elevation, points[0].Elevation

	for i := 1; i < len(points); i++ {
		elevDiff := points[i].Elevation - points[i-1].Elevation
		if elevDiff > 0 {
			totalGain += elevDiff
		} else {
			totalLoss += -elevDiff
		}

		if points[i].Elevation < minElev {
			minElev = points[i].Elevation
		}
		if points[i].Elevation > maxElev {
			maxElev = points[i].Elevation
		}
	}

	return Metadata{
		TotalDistanceKm:     totalDistance,
		TotalPoints:         len(points),
		SegmentDistanceKm:   segmentDistance,
		TotalElevationGainM: totalGain,
		TotalElevationLossM: totalLoss,
		MinElevationM:       minElev,
		MaxElevationM:       maxElev,
	}
}

func divideIntoSegments(points []Point, distances []float64, segmentDistanceKm float64) []Segment {
	if len(points) < 2 {
		return []Segment{}
	}

	var segments []Segment
	segmentNum := 1
	currentSegmentStart := 0.0
	currentSegmentEnd := segmentDistanceKm

	for currentSegmentStart < distances[len(distances)-1] {
		segment := calculateSegmentElevation(
			points, distances,
			currentSegmentStart, currentSegmentEnd,
			segmentNum,
		)

		if segment.PointCount > 0 {
			segments = append(segments, segment)
		}

		segmentNum++
		currentSegmentStart = currentSegmentEnd
		currentSegmentEnd += segmentDistanceKm
	}

	return segments
}

func calculateSegmentElevation(points []Point, distances []float64, startKm, endKm float64, segmentNum int) Segment {
	var segmentPoints []Point
	var segmentDistances []float64

	// Find points within the segment
	for i, distance := range distances {
		if distance >= startKm && distance <= endKm {
			segmentPoints = append(segmentPoints, points[i])
			segmentDistances = append(segmentDistances, distance)
		}
	}

	// If no points in segment, interpolate
	if len(segmentPoints) == 0 {
		// Find surrounding points for interpolation
		startPoint, endPoint := interpolateSegmentBoundaries(points, distances, startKm, endKm)
		if startPoint != nil && endPoint != nil {
			segmentPoints = []Point{*startPoint, *endPoint}
		}
	}

	if len(segmentPoints) < 2 {
		actualEndKm := endKm
		if endKm > distances[len(distances)-1] {
			actualEndKm = distances[len(distances)-1]
		}

		return Segment{
			Number:      segmentNum,
			StartKm:     startKm,
			EndKm:       actualEndKm,
			DistanceKm:  actualEndKm - startKm,
			PointCount:  len(segmentPoints),
		}
	}

	// Calculate elevation changes
	var elevationGain, elevationLoss float64
	startElevation := segmentPoints[0].Elevation
	endElevation := segmentPoints[len(segmentPoints)-1].Elevation

	for i := 1; i < len(segmentPoints); i++ {
		elevDiff := segmentPoints[i].Elevation - segmentPoints[i-1].Elevation
		if elevDiff > 0 {
			elevationGain += elevDiff
		} else {
			elevationLoss += -elevDiff
		}
	}

	actualEndKm := endKm
	if endKm > distances[len(distances)-1] {
		actualEndKm = distances[len(distances)-1]
	}

	return Segment{
		Number:          segmentNum,
		StartKm:         startKm,
		EndKm:           actualEndKm,
		DistanceKm:      actualEndKm - startKm,
		ElevationGainM:  elevationGain,
		ElevationLossM:  elevationLoss,
		NetElevationM:   elevationGain - elevationLoss,
		StartElevationM: startElevation,
		EndElevationM:   endElevation,
		PointCount:      len(segmentPoints),
	}
}

func interpolateSegmentBoundaries(points []Point, distances []float64, startKm, endKm float64) (*Point, *Point) {
	if len(points) < 2 {
		return nil, nil
	}

	var startPoint, endPoint *Point

	// Find or interpolate start point
	for i := 1; i < len(distances); i++ {
		if distances[i-1] <= startKm && distances[i] >= startKm {
			if distances[i-1] == startKm {
				startPoint = &points[i-1]
			} else if distances[i] == startKm {
				startPoint = &points[i]
			} else {
				// Interpolate
				ratio := (startKm - distances[i-1]) / (distances[i] - distances[i-1])
				interpolated := Point{
					Lat:       points[i-1].Lat + ratio*(points[i].Lat-points[i-1].Lat),
					Lon:       points[i-1].Lon + ratio*(points[i].Lon-points[i-1].Lon),
					Elevation: points[i-1].Elevation + ratio*(points[i].Elevation-points[i-1].Elevation),
				}
				startPoint = &interpolated
			}
			break
		}
	}

	// Find or interpolate end point
	for i := 1; i < len(distances); i++ {
		if distances[i-1] <= endKm && distances[i] >= endKm {
			if distances[i-1] == endKm {
				endPoint = &points[i-1]
			} else if distances[i] == endKm {
				endPoint = &points[i]
			} else {
				// Interpolate
				ratio := (endKm - distances[i-1]) / (distances[i] - distances[i-1])
				interpolated := Point{
					Lat:       points[i-1].Lat + ratio*(points[i].Lat-points[i-1].Lat),
					Lon:       points[i-1].Lon + ratio*(points[i].Lon-points[i-1].Lon),
					Elevation: points[i-1].Elevation + ratio*(points[i].Elevation-points[i-1].Elevation),
				}
				endPoint = &interpolated
			}
			break
		}
	}

	return startPoint, endPoint
}