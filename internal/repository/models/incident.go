package models

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain"
)

type Incident struct {
	IncidentID string     `gorm:"primaryKey;column:incident_id"`
	OperatorID string     `gorm:"column:operator_id;not null"`
	Location   string     `gorm:"column:location;type:geometry(POINT,4326);not null"`
	Radius     uint16     `gorm:"column:radius;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;default:null"`
	UpdatedAt  *time.Time `gorm:"column:updated_at;default:null"`
}
type UpdateIncidentInput struct {
	OperatorID *string  `json:"operator_id"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude"`
	Radius     *uint16  `json:"radius"`
}

func IncidentModelToDomain(m *Incident) (*domain.Incident, error) {
	lon, lat, err := ParseEWKBPoint(m.Location)
	if err != nil {
		return nil, err
	}
	return &domain.Incident{
		IncidentID: m.IncidentID,
		OperatorID: m.OperatorID,
		Longitude:  lon,
		Latitude:   lat,
		Radius:     m.Radius,
		CreatedAt:  m.CreatedAt,
		DeletedAt:  m.DeletedAt,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}

func IncidentDomainToModel(i *domain.Incident) *Incident {
	now := time.Now().UTC()
	return &Incident{
		IncidentID: i.IncidentID,
		OperatorID: i.OperatorID,
		Location:   PointWKT(i.Longitude, i.Latitude),
		Radius:     i.Radius,
		CreatedAt:  now,
		UpdatedAt:  &now,
	}
}
func PointWKT(lon, lat float64) string {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", lon, lat)
}
func ParseEWKBPoint(hexStr string) (lon, lat float64, err error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return 0, 0, err
	}

	if len(data) < 9 {
		return 0, 0, errors.New("EWKB too short")
	}

	byteOrder := data[0]
	var bo binary.ByteOrder
	if byteOrder == 0 {
		bo = binary.BigEndian
	} else if byteOrder == 1 {
		bo = binary.LittleEndian
	} else {
		return 0, 0, errors.New("unknown byte order")
	}

	geomTypeWithFlags := bo.Uint32(data[1:5])
	const ewkbHasSRID = 0x20000000

	hasSRID := (geomTypeWithFlags & ewkbHasSRID) != 0
	geomType := geomTypeWithFlags &^ ewkbHasSRID

	if geomType != 1 {
		return 0, 0, errors.New("not a POINT")
	}

	pos := 5

	if hasSRID {
		if len(data) < pos+4 {
			return 0, 0, errors.New("EWKB too short for SRID")
		}

		pos += 4
	}

	if len(data) < pos+16 {
		return 0, 0, errors.New("EWKB too short for coordinates")
	}

	lon = math.Float64frombits(bo.Uint64(data[pos : pos+8]))
	lat = math.Float64frombits(bo.Uint64(data[pos+8 : pos+16]))

	return lon, lat, nil
}
