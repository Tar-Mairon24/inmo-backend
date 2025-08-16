package models

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
    "time"
)

type Property struct {
    ID              uint               `gorm:"primaryKey" json:"id"`
    Title           string             `gorm:"not null;size:255" json:"title"`
    ListingDate     *time.Time         `gorm:"autoCreateTime" json:"listing_date"`
    Address         string             `gorm:"not null;size:500" json:"address"`
    Neighborhood    string             `gorm:"size:255" json:"neighborhood"`
    City            string             `gorm:"not null;size:255" json:"city"`
	Zone			string             `gorm:"size:255" json:"zone"`
    Reference       string             `gorm:"size:500" json:"reference"`
    Price           float64            `gorm:"not null" json:"price"`
    ConstructionM2  int                `gorm:"default:0" json:"construction_m2"`
    LandM2          int                `gorm:"default:0" json:"land_m2"`
    IsOccupied      bool               `gorm:"default:false" json:"is_occupied"`
    IsFurnished     bool               `gorm:"default:false" json:"is_furnished"`
    Floors          int                `gorm:"default:1" json:"floors"`
    Bedrooms        int                `gorm:"default:0" json:"bedrooms"`
    Bathrooms       int                `gorm:"default:0" json:"bathrooms"`
    GarageSize      int                `gorm:"default:0" json:"garage_size"` // Number of cars
    GardenM2        int                `gorm:"default:0" json:"garden_m2"`
    GasTypes        StringArray        `gorm:"type:json" json:"gas_types"`
    Amenities       StringArray        `gorm:"type:json" json:"amenities"`
    Extras          StringArray        `gorm:"type:json" json:"extras"`
    Utilities       StringArray        `gorm:"type:json" json:"utilities"`
    Notes           string             `gorm:"type:text" json:"notes"`
    OwnerID         uint               `gorm:"not null" json:"owner_id"`
    UserID          uint               `gorm:"not null" json:"user_id"`
	PropertyType    PropertyType       `gorm:"not null" json:"property_type"`
    TransactionType TransactionType    `gorm:"not null" json:"transaction_type"`
    Status          PropertyStatus     `gorm:"default:'available'" json:"status"`
    CreatedAt       time.Time          `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt       time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt       *time.Time         `gorm:"index" json:"-"`
	Owner           *User              `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	User            *User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// PropertyResponse represents the public view of a property
type PropertyResponse struct {
    ID              uint            `json:"id"`
    Title           string          `json:"title"`
    Address         string          `json:"address"`
    Neighborhood    string          `json:"neighborhood"`
    City            string          `json:"city"`
	Zone            string          `json:"zone"`
    Price           float64         `json:"price"`
    ConstructionM2  int             `json:"construction_m2"`
    LandM2          int             `json:"land_m2"`
    IsOccupied      bool            `json:"is_occupied"`
    IsFurnished     bool            `json:"is_furnished"`
    Floors          int             `json:"floors"`
    Bedrooms        int             `json:"bedrooms"`
    Bathrooms       int             `json:"bathrooms"`
    GarageSize      int             `json:"garage_size"`
    GardenM2        int             `json:"garden_m2"`
    GasTypes        StringArray     `json:"gas_types"`
    Amenities       StringArray     `json:"amenities"`
    Extras          StringArray     `json:"extras"`
    Utilities       StringArray     `json:"utilities"`
	PropertyType   	PropertyType    `json:"property_type"`
    TransactionType TransactionType `json:"transaction_type"`
    Status          PropertyStatus  `json:"status"`
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
    Agent          	*UserResponse   `json:"agent,omitempty"` // Agent handling the property
}

// PropertyCard represents a simplified property view for listings
type PropertyCard struct {
    ID              uint            `json:"id"`
    Title           string          `json:"title"`
    Price           float64         `json:"price"`
    Bedrooms        int             `json:"bedrooms"`
    Bathrooms       int             `json:"bathrooms"`
    ConstructionM2  int             `json:"construction_m2"`
    City            string          `json:"city"`
    Neighborhood    string          `json:"neighborhood"`
	PropertyType    PropertyType    `json:"property_type"`
    TransactionType TransactionType `json:"transaction_type"`
    Status          PropertyStatus  `json:"status"`
    CreatedAt       time.Time       `json:"created_at"`
}

// Enums for better type safety
type TransactionType string

const (
    TransactionSale   TransactionType = "sale"
    TransactionRental TransactionType = "rental"
)

type PropertyStatus string

const (
    StatusAvailable PropertyStatus = "available"
    StatusSold      PropertyStatus = "sold"
    StatusRented    PropertyStatus = "rented"
    StatusReserved  PropertyStatus = "reserved"
)

type PropertyType string 

const (
	TypeHouse     	PropertyType = "house"
	TypeApartment 	PropertyType = "apartment"
	TypeLand      	PropertyType = "land"
	TypeCommercial 	PropertyType = "commercial"
	TypeStorehouse 	PropertyType = "storehouse"
	TypeOffice     	PropertyType = "office"
	TypeIndustrial  PropertyType = "industrial"
	TypeOther      	PropertyType = "other"
)

type StringArray []string

func (sa *StringArray) Scan(value any) error {
    if value == nil {
        *sa = []string{}
        return nil
    }

    switch v := value.(type) {
    case []byte:
        return json.Unmarshal(v, sa)
    case string:
        return json.Unmarshal([]byte(v), sa)
    default:
        return errors.New("cannot scan StringArray")
    }
}

// Value implements the Valuer interface for database writing
func (sa StringArray) Value() (driver.Value, error) {
    if len(sa) == 0 {
        return "[]", nil
    }
    return json.Marshal(sa)
}

func (p *Property) ToResponse() *PropertyResponse {
    response := &PropertyResponse{
        ID:              p.ID,
        Title:           p.Title,
        Address:         p.Address,
        Neighborhood:    p.Neighborhood,
        City:            p.City,
        Zone:            p.Zone,
        Price:           p.Price,
        ConstructionM2:  p.ConstructionM2,
        LandM2:          p.LandM2,
        IsOccupied:      p.IsOccupied,
        IsFurnished:     p.IsFurnished,
        Floors:          p.Floors,
        Bedrooms:        p.Bedrooms,
        Bathrooms:       p.Bathrooms,
        GarageSize:      p.GarageSize,
        GardenM2:        p.GardenM2,
        GasTypes:        p.GasTypes,
        Amenities:       p.Amenities,
        Extras:          p.Extras,
        Utilities:       p.Utilities,
        PropertyType:    p.PropertyType,
        TransactionType: p.TransactionType,
        Status:          p.Status,
        CreatedAt:       p.CreatedAt,
        UpdatedAt:       p.UpdatedAt,
    }
    
    if p.User != nil && p.User.ID != 0 {
        response.Agent = p.User.ToUserResponse()  // Convert User to UserResponse
    }
    
    return response
}

func (p *Property) ToCard() *PropertyCard {
    return &PropertyCard{
        ID:              p.ID,
        Title:           p.Title,
        Price:           p.Price,
        Bedrooms:        p.Bedrooms,
        Bathrooms:       p.Bathrooms,
        ConstructionM2:  p.ConstructionM2,
        City:            p.City,
        Neighborhood:    p.Neighborhood,
        PropertyType:    p.PropertyType,
        TransactionType: p.TransactionType,
        Status:          p.Status,
        CreatedAt:       p.CreatedAt,
    }
}