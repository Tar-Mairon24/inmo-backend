package usecase

import (
	"errors"

	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
)

type PropertyUseCase struct {
	propertyRepo ports.PropertyRepository
}

func NewPropertyUseCase(propertyRepo ports.PropertyRepository) *PropertyUseCase {
	return &PropertyUseCase{
		propertyRepo: propertyRepo,
	}
}

func (p *PropertyUseCase) GetAllProperties() ([]models.PropertyResponse, error) {
	properties, err := p.propertyRepo.GetAll()
	if err != nil {
		return nil, err
	}
	logrus.Infof("Retrieved %d properties", len(properties))
	return properties, nil
}

func (p *PropertyUseCase) GetPropertyByID(id uint) (*models.PropertyResponse, error) {
	property, err := p.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if property == nil {
		logrus.Warnf("No property found with ID %d", id)
		return nil, errors.New("property not found")
	}
	logrus.Infof("Retrieved property with ID %d", id)
	return property, nil
}

func (p *PropertyUseCase) CreateProperty(property *models.Property) (*models.PropertyResponse, error) {
	if property == nil {
		logrus.Error("Property cannot be nil")
		return nil, errors.New("property cannot be nil")
	}
	if property.Address == "" {
		logrus.Error("Address cannot be empty")
		return nil, errors.New("address cannot be empty")
	}
	if property.Price <= 0 {
		logrus.Error("Price must be greater than zero")
		return nil, errors.New("price must be greater than zero")
	}

	createdProperty, err := p.propertyRepo.Create(property)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Created property with ID %d", createdProperty.ID)
	return createdProperty, nil
}

func (p *PropertyUseCase) UpdateProperty(property *models.Property) (*models.PropertyResponse, error) {
	if property == nil {
		logrus.Error("Property cannot be nil")
		return nil, errors.New("property cannot be nil")
	}
	if property.ID == 0 {
		logrus.Error("Property ID must be provided")
		return nil, errors.New("property ID must be provided")
	}
	if property.Address == "" {
		logrus.Error("Address cannot be empty")
		return nil, errors.New("address cannot be empty")
	}
	if property.Price <= 0 {
		logrus.Error("Price must be greater than zero")
		return nil, errors.New("price must be greater than zero")
	}

	updatedProperty, err := p.propertyRepo.Update(property)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Updated property with ID %d", updatedProperty.ID)
	return updatedProperty, nil
}

func (p *PropertyUseCase) DeleteProperty(id uint) error {
	if id == 0 {
		logrus.Error("Property ID must be provided")
		return errors.New("property ID must be provided")
	}

	err := p.propertyRepo.Delete(id)
	if err != nil {
		return err
	}
	logrus.Infof("Deleted property with ID %d", id)
	return nil
}
