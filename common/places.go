package common

import (
	// "log"

	"github.com/golang/geo/s2"
)

//
// Is current position (e) is in same cell that given place (sp) ?
// As advised, it uses level 16 cells.
//
func IsNearSingle(sp *SignPlace, e *SignPlace) bool {
	sp_latlon := s2.LatLngFromDegrees(sp.GetLatitude(), sp.GetLongitude())
	sp_cell := s2.CellFromLatLng(sp_latlon)

	parent_sp_cell_id := sp_cell.ID().Parent(16)
	parent_sp_cell := s2.CellFromCellID(parent_sp_cell_id)

	latlon := s2.LatLngFromDegrees(e.GetLatitude(), e.GetLongitude())
	p1 := s2.PointFromLatLng(latlon)

	return parent_sp_cell.ContainsPoint(p1)
}

//
// Check a position matches a user's signification places
//
func IsNear(sps []*SignPlace, e *SignPlace) bool {
	for _, sp := range sps {
		r := IsNearSingle(sp, e)
		if r {
			return true
		}
	}

	return false
}

//
// Check a user's signification places match on of other's user
//
func IsNearMultiple(sps []*SignPlace, es []*SignPlace) bool {
	for _, e := range es {
		r := IsNear(sps, e)
		if r {
			return true
		}
	}

	return false
}
