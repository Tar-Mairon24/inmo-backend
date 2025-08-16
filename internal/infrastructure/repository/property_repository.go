package repository

import (
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
)

type PropertyRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func NewPropertyRepository(db *sql.DB) ports.PropertyRepository {
	return &PropertyRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

func (r *PropertyRepository) GetAll() ([]models.PropertyResponse, error) {
	query := r.qb.Select("*").
		From("properties").
		Where(squirrel.Expr("deleted_at IS NULL"))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for getting all properties")
		return nil, err
	}
	rows, err := r.db.Query(sqlStr, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for getting all properties")
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close rows after getting all properties")
		}
	}()

	var properties []models.PropertyResponse
	for rows.Next() {
		var property models.Property

        if err := rows.Scan(
            &property.ID,
            &property.Title,
            &property.ListingDate,
            &property.Address,
            &property.Neighborhood,
            &property.City,
            &property.Zone,
            &property.Reference,
            &property.Price,
            &property.ConstructionM2,
            &property.LandM2,
            &property.IsOccupied,
            &property.IsFurnished,
            &property.Floors,
            &property.Bedrooms,
            &property.Bathrooms,
            &property.GarageSize,
            &property.GardenM2,
            &property.GasTypes,
            &property.Amenities,
            &property.Extras,
            &property.Utilities,
            &property.Notes,
            &property.OwnerID,
            &property.UserID,
            &property.PropertyType,
            &property.TransactionType,
            &property.Status,
            &property.CreatedAt,
            &property.UpdatedAt,
            &property.DeletedAt,
		); err != nil {
			logrus.WithError(err).Error("Failed to scan property row")
			return nil, err
		}

		properties = append(properties, *property.ToResponse())

		if err := rows.Err(); err != nil {
			logrus.WithError(err).Error("Error occurred while iterating over property rows")
			return nil, err
		}
	}

	if len(properties) == 0 {
		logrus.Warn("No properties found in the database")
		return []models.PropertyResponse{}, nil
	}

	return properties, nil
}

func (r *PropertyRepository) GetByID(id uint) (*models.PropertyResponse, error) {
	query := r.qb.Select("*").
		From("properties").
		Where(squirrel.And{
			squirrel.Eq{"id": id},
			squirrel.Expr("deleted_at IS NULL"),
		})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for getting property by ID")
		return nil, err
	}

	var property models.Property
	err = r.db.QueryRow(sqlStr, args...).Scan(
		&property.ID,
		&property.Title,
		&property.ListingDate,
		&property.Address,
		&property.Neighborhood,
		&property.City,
		&property.Zone,
		&property.Reference,
		&property.Price,
		&property.ConstructionM2,
		&property.LandM2,
		&property.IsOccupied,
		&property.IsFurnished,
		&property.Floors,
		&property.Bedrooms,
		&property.Bathrooms,
		&property.GarageSize,
		&property.GardenM2,
		&property.GasTypes,
		&property.Amenities,
		&property.Extras,
		&property.Utilities,
		&property.Notes,
		&property.OwnerID,
		&property.UserID,
		&property.PropertyType,
		&property.TransactionType,
		&property.Status,
		&property.CreatedAt,
		&property.UpdatedAt,
		&property.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithError(err).Warnf("No property found with ID %d", id)
			return nil, err
		}
		logrus.WithError(err).Error("Failed to execute query for getting property by ID")
		return nil, err
	}

	return property.ToResponse(), nil
}

func (r *PropertyRepository) Create(property *models.Property) (*models.PropertyResponse, error) {
    query := r.qb.Insert("properties").
        Columns(
            "title", "listing_date", "address", "neighborhood", "city",
            "zone", "reference", "price", "construction_m2", "land_m2",
            "is_occupied", "is_furnished", "floors", "bedrooms", "bathrooms",
            "garage_size", "garden_m2", "gas_types", "amenities", "extras",
            "utilities", "notes", "owner_id", "user_id", "property_type",
            "transaction_type", "status",
        ).
        Values(
            property.Title, property.ListingDate, property.Address, property.Neighborhood, property.City,
            property.Zone, property.Reference, property.Price, property.ConstructionM2, property.LandM2,
            property.IsOccupied, property.IsFurnished, property.Floors, property.Bedrooms, property.Bathrooms,
            property.GarageSize, property.GardenM2, property.GasTypes, property.Amenities, property.Extras,
            property.Utilities, property.Notes, property.OwnerID, property.UserID, property.PropertyType,
            property.TransactionType, property.Status,
        )

    sqlStr, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for creating a new property")
        return nil, err
    }

    result, err := r.db.Exec(sqlStr, args...)
    if err != nil {
        logrus.WithError(err).Error("Failed to execute query for creating a new property")
        return nil, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        logrus.WithError(err).Error("Failed to get last insert ID")
        return nil, err
    }

    property.ID = uint(id)
    logrus.Infof("Property created successfully with ID: %d", property.ID)
    return property.ToResponse(), nil
}

func (r *PropertyRepository) Update(property *models.Property) (*models.PropertyResponse, error) {
	query := r.qb.Update("properties").
		Set("title", property.Title).
		Set("listing_date", property.ListingDate).
		Set("address", property.Address).
		Set("neighborhood", property.Neighborhood).
		Set("city", property.City).
		Set("zone", property.Zone).
		Set("reference", property.Reference).
		Set("price", property.Price).
		Set("construction_m2", property.ConstructionM2).
		Set("land_m2", property.LandM2).
		Set("is_occupied", property.IsOccupied).
		Set("is_furnished", property.IsFurnished).
		Set("floors", property.Floors).
		Set("bedrooms", property.Bedrooms).
		Set("bathrooms", property.Bathrooms).
		Set("garage_size", property.GarageSize).
		Set("garden_m2", property.GardenM2).
		Set("gas_types", property.GasTypes).
		Set("amenities", property.Amenities).
		Set("extras", property.Extras).
		Set("utilities", property.Utilities).
		Set("notes", property.Notes).
		Set("owner_id", property.OwnerID).
		Set("user_id", property.UserID).
		Set("property_type", property.PropertyType).
		Set("transaction_type", property.TransactionType).
		Set("status", property.Status).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": property.ID}).
		Where(squirrel.Expr("deleted_at IS NULL"))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for updating a property")
		return nil, err
	}

	result, err := r.db.Exec(sqlStr, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for updating a property")
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Error("Failed to get rows affected after updating a property")
		return nil, err
	}

	if rowsAffected == 0 {
		logrus.Warnf("No rows updated for Property ID %d. It may not exist or has already been deleted.", property.ID)
		return nil, errors.New("property not found or already deleted")
	}

	logrus.Infof("Property with ID %d updated successfully", property.ID)
	return property.ToResponse(), nil
}

func (r *PropertyRepository) Delete(id uint) error {
	query := r.qb.Update("properties").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Expr("deleted_at IS NULL"))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for deleting a property")
		return err
	}

	result, err := r.db.Exec(sqlStr, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for deleting a property")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Error("Failed to get rows affected after deleting a property")
		return err
	}

	if rowsAffected == 0 {
		logrus.Warnf("No property found with ID %d or already deleted", id)
		return errors.New("property not found or already deleted")
	}

	return nil
} 