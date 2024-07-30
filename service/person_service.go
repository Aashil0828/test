package service

import (
	"database/sql"
	"infilon_task/models"
	"infilon_task/utils"
)

// PersonService struct that holds the database dependency
type PersonService struct {
	DB *sql.DB
}

// NewPersonService creates a new PersonService
func NewPersonService(db *sql.DB) *PersonService {
	return &PersonService{
		DB: db,
	}
}

// CheckPersonExists checks if a person with the given name already exists in the database
func (ps *PersonService) CheckPersonExists(name string) (bool, error) {
	var count int
	err := ps.DB.QueryRow(`SELECT COUNT(*) FROM person WHERE name = ?`, name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreatePerson creates a new person and their related phone and address records and inserts them into the database
func (ps *PersonService) CreatePerson(personRequest models.PersonRequest) (models.PersonRequest, error) {

	// Check if a person with the same name already exists in case it is needed, or we can omit this check
	exists, err := ps.CheckPersonExists(personRequest.Name)
	if err != nil {
		return models.PersonRequest{}, err
	}
	if exists {
		return models.PersonRequest{}, utils.ErrAlreadyExists
	}

	tx, err := ps.DB.Begin()
	if err != nil {
		return models.PersonRequest{}, err
	}

	// Insert new person
	result, err := tx.Exec(`
        INSERT INTO person (name, age)
        VALUES (?, ?)`, personRequest.Name, personRequest.Age)
	if err != nil {
		tx.Rollback()
		return models.PersonRequest{}, err
	}

	personID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return models.PersonRequest{}, err
	}

	// Insert new address
	AddressResult, err := tx.Exec(`
        INSERT INTO address (city, state, street1, street2, zip_code)
        VALUES (?, ?, ?, ?, ?)`, personRequest.City, personRequest.State, personRequest.Street1, personRequest.Street2, personRequest.ZipCode)
	if err != nil {
		tx.Rollback()
		return models.PersonRequest{}, err
	}

	// Get the address ID
	addressID, err := AddressResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		return models.PersonRequest{}, err
	}

	// Insert address join
	_, err = tx.Exec(`
        INSERT INTO address_join (person_id, address_id)
        VALUES (?, ?)`, personID, addressID)
	if err != nil {
		tx.Rollback()
		return models.PersonRequest{}, err
	}

	// Insert phone number
	_, err = tx.Exec(`
        INSERT INTO phone (number, person_id)
        VALUES (?, ?)`, personRequest.PhoneNumber, personID)
	if err != nil {
		tx.Rollback()
		return models.PersonRequest{}, err
	}

	if err = tx.Commit(); err != nil {
		return models.PersonRequest{}, err
	}
	personRequest.Id = personID
	return personRequest, nil
}

// GetPersonInfo retrieves the details of a person by their ID
func (ps *PersonService) GetPersonInfo(personID int) (models.PersonRequest, error) {
	var person models.PersonRequest

	// Fetch person details, phone number, and address
	query := `
        SELECT p.id, p.name, p.age, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code
        FROM person p
        LEFT JOIN phone ph ON p.id = ph.person_id
        LEFT JOIN address_join aj ON p.id = aj.person_id
        LEFT JOIN address a ON aj.address_id = a.id
        WHERE p.id = ?`

	row := ps.DB.QueryRow(query, personID)

	err := row.Scan(
		&person.Id,
		&person.Name,
		&person.Age,
		&person.PhoneNumber,
		&person.City,
		&person.State,
		&person.Street1,
		&person.Street2,
		&person.ZipCode,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return person, nil
		}
		return person, err
	}

	return person, nil
}
