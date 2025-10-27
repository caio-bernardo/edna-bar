package domain

import "fmt"

// Domain error types
var (
	// Book related errors
	ErrBookNotFound         = fmt.Errorf("book not found")
	ErrBookAlreadyExists    = fmt.Errorf("book already exists")
	ErrInvalidISBN          = fmt.Errorf("invalid ISBN")
	ErrInvalidPublicationDate = fmt.Errorf("invalid publication date")
	ErrBookTitleRequired    = fmt.Errorf("book title is required")

	// Author related errors
	ErrAuthorNotFound       = fmt.Errorf("author not found")
	ErrAuthorAlreadyExists  = fmt.Errorf("author already exists")
	ErrInvalidRG            = fmt.Errorf("invalid RG")
	ErrAuthorNameRequired   = fmt.Errorf("author name is required")
	ErrAuthorAddressRequired = fmt.Errorf("author address is required")

	// Publisher related errors
	ErrPublisherNotFound      = fmt.Errorf("publisher not found")
	ErrPublisherAlreadyExists = fmt.Errorf("publisher already exists")
	ErrPublisherNameRequired  = fmt.Errorf("publisher name is required")
	ErrPublisherAddressRequired = fmt.Errorf("publisher address is required")

	// Printing company related errors
	ErrGraficaNotFound       = fmt.Errorf("printing company not found")
	ErrGraficaAlreadyExists  = fmt.Errorf("printing company already exists")
	ErrGraficaNameRequired   = fmt.Errorf("printing company name is required")
	ErrInvalidGraficaType    = fmt.Errorf("invalid printing company type")

	// Contract related errors
	ErrContractNotFound      = fmt.Errorf("contract not found")
	ErrContractAlreadyExists = fmt.Errorf("contract already exists")
	ErrInvalidContractValue  = fmt.Errorf("invalid contract value")
	ErrContractResponsableRequired = fmt.Errorf("contract responsible person is required")
	ErrOnlyContractedGraficaCanHaveContract = fmt.Errorf("only contracted printing companies can have contracts")

	// Printing job related errors
	ErrPrintingJobNotFound    = fmt.Errorf("printing job not found")
	ErrPrintingJobAlreadyExists = fmt.Errorf("printing job already exists")
	ErrInvalidCopiesNumber    = fmt.Errorf("invalid number of copies")
	ErrInvalidDeliveryDate    = fmt.Errorf("invalid delivery date")
	ErrDeliveryDateInPast     = fmt.Errorf("delivery date cannot be in the past")

	// Relationship errors
	ErrRelationshipNotFound    = fmt.Errorf("relationship not found")
	ErrRelationshipAlreadyExists = fmt.Errorf("relationship already exists")
	ErrCannotRemoveLastAuthor  = fmt.Errorf("cannot remove the last author from a book")

	// Business rule errors
	ErrBookMustHaveAtLeastOneAuthor = fmt.Errorf("book must have at least one author")
	ErrContractedGraficaMustHaveAddress = fmt.Errorf("contracted printing company must have an address")
	ErrParticularGraficaCannotHaveContract = fmt.Errorf("private printing company cannot have contracts")

	// Validation errors
	ErrFieldTooLong    = fmt.Errorf("field value is too long")
	ErrFieldRequired   = fmt.Errorf("field is required")
	ErrInvalidFormat   = fmt.Errorf("invalid format")
	ErrNegativeValue   = fmt.Errorf("value cannot be negative")
	ErrZeroValue       = fmt.Errorf("value cannot be zero")
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Field   string
	Value   any
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("[%s] %s (field: %s, value: %v)", e.Code, e.Message, e.Field, e.Value)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// NewFieldError creates a new domain error for a specific field
func NewFieldError(code, message, field string, value interface{}) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Field:   field,
		Value:   value,
	}
}

// Error codes
const (
	// Book error codes
	CodeBookNotFound         = "BOOK_NOT_FOUND"
	CodeBookAlreadyExists    = "BOOK_ALREADY_EXISTS"
	CodeInvalidISBN          = "INVALID_ISBN"
	CodeInvalidPublicationDate = "INVALID_PUBLICATION_DATE"
	CodeBookTitleRequired    = "BOOK_TITLE_REQUIRED"

	// Author error codes
	CodeAuthorNotFound       = "AUTHOR_NOT_FOUND"
	CodeAuthorAlreadyExists  = "AUTHOR_ALREADY_EXISTS"
	CodeInvalidRG            = "INVALID_RG"
	CodeAuthorNameRequired   = "AUTHOR_NAME_REQUIRED"
	CodeAuthorAddressRequired = "AUTHOR_ADDRESS_REQUIRED"

	// Publisher error codes
	CodePublisherNotFound      = "PUBLISHER_NOT_FOUND"
	CodePublisherAlreadyExists = "PUBLISHER_ALREADY_EXISTS"
	CodePublisherNameRequired  = "PUBLISHER_NAME_REQUIRED"
	CodePublisherAddressRequired = "PUBLISHER_ADDRESS_REQUIRED"

	// Printing company error codes
	CodeGraficaNotFound       = "GRAFICA_NOT_FOUND"
	CodeGraficaAlreadyExists  = "GRAFICA_ALREADY_EXISTS"
	CodeGraficaNameRequired   = "GRAFICA_NAME_REQUIRED"
	CodeInvalidGraficaType    = "INVALID_GRAFICA_TYPE"

	// Contract error codes
	CodeContractNotFound      = "CONTRACT_NOT_FOUND"
	CodeContractAlreadyExists = "CONTRACT_ALREADY_EXISTS"
	CodeInvalidContractValue  = "INVALID_CONTRACT_VALUE"
	CodeContractResponsableRequired = "CONTRACT_RESPONSABLE_REQUIRED"
	CodeOnlyContractedGraficaCanHaveContract = "ONLY_CONTRACTED_GRAFICA_CAN_HAVE_CONTRACT"

	// Printing job error codes
	CodePrintingJobNotFound    = "PRINTING_JOB_NOT_FOUND"
	CodePrintingJobAlreadyExists = "PRINTING_JOB_ALREADY_EXISTS"
	CodeInvalidCopiesNumber    = "INVALID_COPIES_NUMBER"
	CodeInvalidDeliveryDate    = "INVALID_DELIVERY_DATE"
	CodeDeliveryDateInPast     = "DELIVERY_DATE_IN_PAST"

	// Relationship error codes
	CodeRelationshipNotFound    = "RELATIONSHIP_NOT_FOUND"
	CodeRelationshipAlreadyExists = "RELATIONSHIP_ALREADY_EXISTS"
	CodeCannotRemoveLastAuthor  = "CANNOT_REMOVE_LAST_AUTHOR"

	// Business rule error codes
	CodeBookMustHaveAtLeastOneAuthor = "BOOK_MUST_HAVE_AT_LEAST_ONE_AUTHOR"
	CodeContractedGraficaMustHaveAddress = "CONTRACTED_GRAFICA_MUST_HAVE_ADDRESS"
	CodeParticularGraficaCannotHaveContract = "PARTICULAR_GRAFICA_CANNOT_HAVE_CONTRACT"

	// Validation error codes
	CodeFieldTooLong    = "FIELD_TOO_LONG"
	CodeFieldRequired   = "FIELD_REQUIRED"
	CodeInvalidFormat   = "INVALID_FORMAT"
	CodeNegativeValue   = "NEGATIVE_VALUE"
	CodeZeroValue       = "ZERO_VALUE"
)

// Validation error constructors
func NewBookNotFoundError(isbn string) *DomainError {
	return NewFieldError(CodeBookNotFound, "Book not found", "isbn", isbn)
}

func NewBookAlreadyExistsError(isbn string) *DomainError {
	return NewFieldError(CodeBookAlreadyExists, "Book already exists", "isbn", isbn)
}

func NewAuthorNotFoundError(rg string) *DomainError {
	return NewFieldError(CodeAuthorNotFound, "Author not found", "rg", rg)
}

func NewAuthorAlreadyExistsError(rg string) *DomainError {
	return NewFieldError(CodeAuthorAlreadyExists, "Author already exists", "rg", rg)
}

func NewPublisherNotFoundError(id int) *DomainError {
	return NewFieldError(CodePublisherNotFound, "Publisher not found", "id", id)
}

func NewPublisherAlreadyExistsError(name string) *DomainError {
	return NewFieldError(CodePublisherAlreadyExists, "Publisher already exists", "name", name)
}

func NewGraficaNotFoundError(id int) *DomainError {
	return NewFieldError(CodeGraficaNotFound, "Printing company not found", "id", id)
}

func NewGraficaAlreadyExistsError(name string) *DomainError {
	return NewFieldError(CodeGraficaAlreadyExists, "Printing company already exists", "name", name)
}

func NewContractNotFoundError(id int) *DomainError {
	return NewFieldError(CodeContractNotFound, "Contract not found", "id", id)
}

func NewPrintingJobNotFoundError(isbn string, graficaID int) *DomainError {
	return NewFieldError(CodePrintingJobNotFound, "Printing job not found", "isbn_grafica",
		fmt.Sprintf("%s-%d", isbn, graficaID))
}

func NewPrintingJobAlreadyExistsError(isbn string, graficaID int) *DomainError {
	return NewFieldError(CodePrintingJobAlreadyExists, "Printing job already exists", "isbn_grafica",
		fmt.Sprintf("%s-%d", isbn, graficaID))
}

func NewRelationshipNotFoundError(isbn, rg string) *DomainError {
	return NewFieldError(CodeRelationshipNotFound, "Author-book relationship not found", "isbn_rg",
		fmt.Sprintf("%s-%s", isbn, rg))
}

func NewRelationshipAlreadyExistsError(isbn, rg string) *DomainError {
	return NewFieldError(CodeRelationshipAlreadyExists, "Author-book relationship already exists", "isbn_rg",
		fmt.Sprintf("%s-%s", isbn, rg))
}

// Validation helper functions
func NewFieldRequiredError(field string) *DomainError {
	return NewFieldError(CodeFieldRequired, "Field is required", field, nil)
}

func NewFieldTooLongError(field string, maxLength int, actualLength int) *DomainError {
	return NewFieldError(CodeFieldTooLong,
		fmt.Sprintf("Field exceeds maximum length of %d characters", maxLength),
		field, actualLength)
}

func NewInvalidFormatError(field string, value interface{}) *DomainError {
	return NewFieldError(CodeInvalidFormat, "Field has invalid format", field, value)
}

func NewNegativeValueError(field string, value interface{}) *DomainError {
	return NewFieldError(CodeNegativeValue, "Field cannot have negative value", field, value)
}

func NewZeroValueError(field string) *DomainError {
	return NewFieldError(CodeZeroValue, "Field cannot be zero", field, nil)
}
