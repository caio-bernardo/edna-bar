package domain

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ISBN represents a book's International Standard Book Number
type ISBN struct {
	value string
}

// NewISBN creates a new ISBN with validation
func NewISBN(value string) (*ISBN, error) {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(value, "-", ""), " ", "")
	
	if len(cleaned) != 10 && len(cleaned) != 13 {
		return nil, fmt.Errorf("ISBN must be 10 or 13 digits long")
	}
	
	// Basic format validation
	if !regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
		return nil, fmt.Errorf("ISBN must contain only digits")
	}
	
	return &ISBN{value: cleaned}, nil
}

// Value returns the ISBN value
func (isbn *ISBN) Value() string {
	return isbn.value
}

// String returns formatted ISBN
func (isbn *ISBN) String() string {
	if len(isbn.value) == 13 {
		return fmt.Sprintf("%s-%s-%s-%s-%s", 
			isbn.value[0:3], 
			isbn.value[3:4], 
			isbn.value[4:9], 
			isbn.value[9:12], 
			isbn.value[12:13])
	}
	return fmt.Sprintf("%s-%s-%s-%s", 
		isbn.value[0:1], 
		isbn.value[1:6], 
		isbn.value[6:9], 
		isbn.value[9:10])
}

// IsISBN13 checks if the ISBN is 13 digits
func (isbn *ISBN) IsISBN13() bool {
	return len(isbn.value) == 13
}

// RG represents a Brazilian identity document number
type RG struct {
	value string
}

// NewRG creates a new RG with validation
func NewRG(value string) (*RG, error) {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(value, ".", ""), "-", "")
	cleaned = strings.TrimSpace(cleaned)
	
	if len(cleaned) < 7 || len(cleaned) > 20 {
		return nil, fmt.Errorf("RG must be between 7 and 20 characters")
	}
	
	return &RG{value: cleaned}, nil
}

// Value returns the RG value
func (rg *RG) Value() string {
	return rg.value
}

// String returns the RG value
func (rg *RG) String() string {
	return rg.value
}

// Money represents a monetary value
type Money struct {
	amount   int64 // stored in cents to avoid floating point issues
	currency string
}

// NewMoney creates a new Money value object
func NewMoney(amount float64, currency string) (*Money, error) {
	if amount < 0 {
		return nil, fmt.Errorf("money amount cannot be negative")
	}
	
	if currency == "" {
		currency = "BRL" // Default to Brazilian Real
	}
	
	// Convert to cents to avoid floating point precision issues
	cents := int64(amount * 100)
	
	return &Money{
		amount:   cents,
		currency: currency,
	}, nil
}

// Amount returns the amount as float64
func (m *Money) Amount() float64 {
	return float64(m.amount) / 100
}

// AmountInCents returns the amount in cents
func (m *Money) AmountInCents() int64 {
	return m.amount
}

// Currency returns the currency code
func (m *Money) Currency() string {
	return m.currency
}

// String returns formatted money
func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Amount(), m.currency)
}

// Add adds another money value (must be same currency)
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("cannot add different currencies")
	}
	
	return &Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract subtracts another money value (must be same currency)
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, fmt.Errorf("cannot subtract different currencies")
	}
	
	if m.amount < other.amount {
		return nil, fmt.Errorf("insufficient funds")
	}
	
	return &Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

// Address represents a physical address
type Address struct {
	street     string
	number     string
	complement string
	district   string
	city       string
	state      string
	zipCode    string
	country    string
}

// NewAddress creates a new Address value object
func NewAddress(street, number, district, city, state, zipCode string) (*Address, error) {
	if street == "" || city == "" || state == "" {
		return nil, fmt.Errorf("street, city, and state are required")
	}
	
	return &Address{
		street:   strings.TrimSpace(street),
		number:   strings.TrimSpace(number),
		district: strings.TrimSpace(district),
		city:     strings.TrimSpace(city),
		state:    strings.TrimSpace(state),
		zipCode:  strings.TrimSpace(zipCode),
		country:  "Brazil", // Default to Brazil
	}, nil
}

// Street returns the street name
func (a *Address) Street() string {
	return a.street
}

// Number returns the street number
func (a *Address) Number() string {
	return a.number
}

// District returns the district
func (a *Address) District() string {
	return a.district
}

// City returns the city
func (a *Address) City() string {
	return a.city
}

// State returns the state
func (a *Address) State() string {
	return a.state
}

// ZipCode returns the zip code
func (a *Address) ZipCode() string {
	return a.zipCode
}

// Country returns the country
func (a *Address) Country() string {
	return a.country
}

// String returns formatted address
func (a *Address) String() string {
	parts := []string{}
	
	if a.street != "" {
		streetPart := a.street
		if a.number != "" {
			streetPart += ", " + a.number
		}
		parts = append(parts, streetPart)
	}
	
	if a.district != "" {
		parts = append(parts, a.district)
	}
	
	if a.city != "" {
		cityPart := a.city
		if a.state != "" {
			cityPart += " - " + a.state
		}
		parts = append(parts, cityPart)
	}
	
	if a.zipCode != "" {
		parts = append(parts, "CEP: "+a.zipCode)
	}
	
	return strings.Join(parts, ", ")
}

// PublicationDate represents a book's publication date
type PublicationDate struct {
	date time.Time
}

// NewPublicationDate creates a new PublicationDate
func NewPublicationDate(date time.Time) (*PublicationDate, error) {
	if date.IsZero() {
		return nil, fmt.Errorf("publication date cannot be zero")
	}
	
	if date.After(time.Now()) {
		return nil, fmt.Errorf("publication date cannot be in the future")
	}
	
	return &PublicationDate{date: date}, nil
}

// Date returns the publication date
func (pd *PublicationDate) Date() time.Time {
	return pd.date
}

// Year returns the publication year
func (pd *PublicationDate) Year() int {
	return pd.date.Year()
}

// IsClassic checks if the book is considered a classic (over 50 years old)
func (pd *PublicationDate) IsClassic() bool {
	return time.Since(pd.date).Hours() > (50 * 365 * 24)
}

// Age returns the age of the publication in years
func (pd *PublicationDate) Age() int {
	return int(time.Since(pd.date).Hours() / 24 / 365)
}

// String returns formatted publication date
func (pd *PublicationDate) String() string {
	return pd.date.Format("2006-01-02")
}

// Quantity represents a quantity of items
type Quantity struct {
	value int
}

// NewQuantity creates a new Quantity
func NewQuantity(value int) (*Quantity, error) {
	if value < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}
	
	return &Quantity{value: value}, nil
}

// Value returns the quantity value
func (q *Quantity) Value() int {
	return q.value
}

// IsZero checks if quantity is zero
func (q *Quantity) IsZero() bool {
	return q.value == 0
}

// Add adds to the quantity
func (q *Quantity) Add(amount int) (*Quantity, error) {
	if amount < 0 {
		return nil, fmt.Errorf("cannot add negative amount")
	}
	
	return NewQuantity(q.value + amount)
}

// Subtract subtracts from the quantity
func (q *Quantity) Subtract(amount int) (*Quantity, error) {
	if amount < 0 {
		return nil, fmt.Errorf("cannot subtract negative amount")
	}
	
	if q.value < amount {
		return nil, fmt.Errorf("insufficient quantity")
	}
	
	return NewQuantity(q.value - amount)
}

// String returns the quantity as string
func (q *Quantity) String() string {
	return strconv.Itoa(q.value)
}