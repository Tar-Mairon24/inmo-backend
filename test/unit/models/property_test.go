package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"inmo-backend/internal/domain/models"
)

func TestStringArrayScan_NilValue(t *testing.T) {
	var sa models.StringArray
	err := sa.Scan(nil)
	assert.NoError(t, err)
	assert.Equal(t, models.StringArray{}, sa)
}

func TestStringArrayScan_ByteSlice(t *testing.T) {
	var sa models.StringArray
	input := []byte(`["gas","electric"]`)
	err := sa.Scan(input)
	assert.NoError(t, err)
	assert.Equal(t, models.StringArray{"gas", "electric"}, sa)
}

func TestStringArrayScan_String(t *testing.T) {
	var sa models.StringArray
	input := `["water","internet"]`
	err := sa.Scan(input)
	assert.NoError(t, err)
	assert.Equal(t, models.StringArray{"water", "internet"}, sa)
}

func TestStringArrayScan_InvalidType(t *testing.T) {
	var sa models.StringArray
	err := sa.Scan(123)
	assert.Error(t, err)
	assert.Equal(t, "cannot scan StringArray", err.Error())
}

func TestStringArrayScan_InvalidJSON(t *testing.T) {
	var sa models.StringArray
	input := []byte(`not a json`)
	err := sa.Scan(input)
	assert.Error(t, err)
	_, ok := err.(*json.SyntaxError)
	assert.True(t, ok)
}
func TestStringArrayValue_Empty(t *testing.T) {
	sa := models.StringArray{}
	val, err := sa.Value()
	assert.NoError(t, err)
	assert.Equal(t, "[]", val)
}

func TestStringArrayValue_NonEmpty(t *testing.T) {
	sa := models.StringArray{"gas", "electric"}
	val, err := sa.Value()
	assert.NoError(t, err)
	// Should be a valid JSON array string
	expected, _ := json.Marshal(sa)
	assert.Equal(t, string(expected), string(val.([]byte)))
}
func TestProperty_ToResponse_NoUser(t *testing.T) {
	p := &models.Property{
		ID:             1,
		Title:          "Test Property",
		Address:        "123 Main St",
		Neighborhood:   "Downtown",
		City:           "Metropolis",
		Zone:           "Central",
		Price:          100000.0,
		ConstructionM2: 120,
		LandM2:         200,
		IsOccupied:     true,
		IsFurnished:    false,
		Floors:         2,
		Bedrooms:       3,
		Bathrooms:      2,
		GarageSize:     1,
		GardenM2:       50,
		GasTypes:       models.StringArray{"gas"},
		Amenities:      models.StringArray{"pool"},
		Extras:         models.StringArray{"alarm"},
		Utilities:      models.StringArray{"water"},
		PropertyType:   models.TypeHouse,
		TransactionType: models.TransactionSale,
		Status:         models.StatusAvailable,
	}
	resp := p.ToResponse()
	assert.Equal(t, p.ID, resp.ID)
	assert.Equal(t, p.Title, resp.Title)
	assert.Equal(t, p.Address, resp.Address)
	assert.Equal(t, p.Neighborhood, resp.Neighborhood)
	assert.Equal(t, p.City, resp.City)
	assert.Equal(t, p.Zone, resp.Zone)
	assert.Equal(t, p.Price, resp.Price)
	assert.Equal(t, p.ConstructionM2, resp.ConstructionM2)
	assert.Equal(t, p.LandM2, resp.LandM2)
	assert.Equal(t, p.IsOccupied, resp.IsOccupied)
	assert.Equal(t, p.IsFurnished, resp.IsFurnished)
	assert.Equal(t, p.Floors, resp.Floors)
	assert.Equal(t, p.Bedrooms, resp.Bedrooms)
	assert.Equal(t, p.Bathrooms, resp.Bathrooms)
	assert.Equal(t, p.GarageSize, resp.GarageSize)
	assert.Equal(t, p.GardenM2, resp.GardenM2)
	assert.Equal(t, p.GasTypes, resp.GasTypes)
	assert.Equal(t, p.Amenities, resp.Amenities)
	assert.Equal(t, p.Extras, resp.Extras)
	assert.Equal(t, p.Utilities, resp.Utilities)
	assert.Equal(t, p.PropertyType, resp.PropertyType)
	assert.Equal(t, p.TransactionType, resp.TransactionType)
	assert.Equal(t, p.Status, resp.Status)
	assert.Nil(t, resp.Agent)
}

func TestProperty_ToResponse_WithUser(t *testing.T) {
	user := &models.User{
		ID:   42,
		Username: "Agent Smith",
	}
	p := &models.Property{
		ID:    2,
		Title: "Another Property",
		User:  user,
	}
	resp := p.ToResponse()
	assert.NotNil(t, resp.Agent)
	assert.Equal(t, user.ID, resp.Agent.ID)
	assert.Equal(t, user.Username, resp.Agent.Username)
}

func TestProperty_ToResponse_UserZeroID(t *testing.T) {
	user := &models.User{
		ID:   0,
		Username: "No Agent",
	}
	p := &models.Property{
		ID:    3,
		Title: "No Agent Property",
		User:  user,
	}
	resp := p.ToResponse()
	assert.Nil(t, resp.Agent)
}
func TestProperty_ToCard_FullFields(t *testing.T) {
	p := &models.Property{
		ID:              10,
		Title:           "Luxury Villa",
		Price:           500000.0,
		Bedrooms:        5,
		Bathrooms:       4,
		ConstructionM2:  350,
		City:            "Springfield",
		Neighborhood:    "Green Acres",
		PropertyType:    models.TypeHouse,
		TransactionType: models.TransactionSale,
		Status:          models.StatusAvailable,
		CreatedAt:       time.Time{}, // zero value for simplicity
	}
	card := p.ToCard()
	assert.Equal(t, p.ID, card.ID)
	assert.Equal(t, p.Title, card.Title)
	assert.Equal(t, p.Price, card.Price)
	assert.Equal(t, p.Bedrooms, card.Bedrooms)
	assert.Equal(t, p.Bathrooms, card.Bathrooms)
	assert.Equal(t, p.ConstructionM2, card.ConstructionM2)
	assert.Equal(t, p.City, card.City)
	assert.Equal(t, p.Neighborhood, card.Neighborhood)
	assert.Equal(t, p.PropertyType, card.PropertyType)
	assert.Equal(t, p.TransactionType, card.TransactionType)
	assert.Equal(t, p.Status, card.Status)
	assert.Equal(t, p.CreatedAt, card.CreatedAt)
}

func TestProperty_ToCard_ZeroValues(t *testing.T) {
	p := &models.Property{}
	card := p.ToCard()
	assert.Equal(t, uint(0), card.ID)
	assert.Equal(t, "", card.Title)
	assert.Equal(t, float64(0), card.Price)
	assert.Equal(t, 0, card.Bedrooms)
	assert.Equal(t, 0, card.Bathrooms)
	assert.Equal(t, 0, card.ConstructionM2)
	assert.Equal(t, "", card.City)
	assert.Equal(t, "", card.Neighborhood)
	assert.Equal(t, models.PropertyType(""), card.PropertyType)
	assert.Equal(t, models.TransactionType(""), card.TransactionType)
	assert.Equal(t, models.PropertyStatus(""), card.Status)
	assert.True(t, card.CreatedAt.IsZero())
}


