package gpx

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type GPX struct {
	XMLName xml.Name `xml:"gpx"`
	Routes  []Route  `xml:"rte"`
	Tracks  []Track  `xml:"trk"`
}

type Route struct {
	Name   string      `xml:"name"`
	Points []RoutePoint `xml:"rtept"`
}

type Track struct {
	Name     string        `xml:"name"`
	Segments []TrackSegment `xml:"trkseg"`
}

type TrackSegment struct {
	Points []TrackPoint `xml:"trkpt"`
}

type RoutePoint struct {
	Lat       float64   `xml:"lat,attr"`
	Lon       float64   `xml:"lon,attr"`
	Elevation *float64  `xml:"ele"`
	Time      *time.Time `xml:"time"`
	Name      string    `xml:"name"`
}

type TrackPoint struct {
	Lat       float64   `xml:"lat,attr"`
	Lon       float64   `xml:"lon,attr"`
	Elevation *float64  `xml:"ele"`
	Time      *time.Time `xml:"time"`
}

type Point struct {
	Lat       float64
	Lon       float64
	Elevation float64
	Time      *time.Time
	Name      string
}

func ParseGPX(reader io.Reader) (*GPX, error) {
	var gpx GPX
	decoder := xml.NewDecoder(reader)

	if err := decoder.Decode(&gpx); err != nil {
		return nil, fmt.Errorf("failed to parse GPX: %w", err)
	}

	return &gpx, nil
}

func (g *GPX) GetAllPoints() []Point {
	var points []Point

	// Extract points from routes
	for _, route := range g.Routes {
		for _, rpt := range route.Points {
			elevation := 0.0
			if rpt.Elevation != nil {
				elevation = *rpt.Elevation
			}
			points = append(points, Point{
				Lat:       rpt.Lat,
				Lon:       rpt.Lon,
				Elevation: elevation,
				Time:      rpt.Time,
				Name:      rpt.Name,
			})
		}
	}

	// Extract points from tracks
	for _, track := range g.Tracks {
		for _, segment := range track.Segments {
			for _, tpt := range segment.Points {
				elevation := 0.0
				if tpt.Elevation != nil {
					elevation = *tpt.Elevation
				}
				points = append(points, Point{
					Lat:       tpt.Lat,
					Lon:       tpt.Lon,
					Elevation: elevation,
					Time:      tpt.Time,
				})
			}
		}
	}

	return points
}

func (p Point) String() string {
	return fmt.Sprintf("Point{Lat: %.6f, Lon: %.6f, Elevation: %.2f}",
		p.Lat, p.Lon, p.Elevation)
}