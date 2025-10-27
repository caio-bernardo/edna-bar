package domain

import (
	"time"
)

// DomainEvent represents a domain event that occurred in the system
type DomainEvent interface {
	EventType() string
	AggregateID() string
	OccurredAt() time.Time
	Version() int
}

// BaseDomainEvent provides common fields for all domain events
type BaseDomainEvent struct {
	eventType   string
	aggregateID string
	occurredAt  time.Time
	version     int
}

// EventType returns the type of the event
func (e *BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the ID of the aggregate that generated the event
func (e *BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

// OccurredAt returns when the event occurred
func (e *BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// Version returns the event version
func (e *BaseDomainEvent) Version() int {
	return e.version
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType, aggregateID string, version int) BaseDomainEvent {
	return BaseDomainEvent{
		eventType:   eventType,
		aggregateID: aggregateID,
		occurredAt:  time.Now(),
		version:     version,
	}
}

// Book-related events

// BookCreatedEvent represents a book creation event
type BookCreatedEvent struct {
	BaseDomainEvent
	ISBN             string    `json:"isbn"`
	Title            string    `json:"title"`
	PublicationDate  time.Time `json:"publication_date"`
	PublisherID      int       `json:"publisher_id"`
}

// NewBookCreatedEvent creates a new book created event
func NewBookCreatedEvent(isbn, title string, publicationDate time.Time, publisherID int) *BookCreatedEvent {
	return &BookCreatedEvent{
		BaseDomainEvent:  NewBaseDomainEvent("BookCreated", isbn, 1),
		ISBN:             isbn,
		Title:            title,
		PublicationDate:  publicationDate,
		PublisherID:      publisherID,
	}
}

// BookUpdatedEvent represents a book update event
type BookUpdatedEvent struct {
	BaseDomainEvent
	ISBN            string    `json:"isbn"`
	Title           string    `json:"title"`
	PublicationDate time.Time `json:"publication_date"`
	PublisherID     int       `json:"publisher_id"`
	Changes         map[string]interface{} `json:"changes"`
}

// NewBookUpdatedEvent creates a new book updated event
func NewBookUpdatedEvent(isbn string, changes map[string]interface{}) *BookUpdatedEvent {
	return &BookUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("BookUpdated", isbn, 1),
		ISBN:            isbn,
		Changes:         changes,
	}
}

// BookDeletedEvent represents a book deletion event
type BookDeletedEvent struct {
	BaseDomainEvent
	ISBN  string `json:"isbn"`
	Title string `json:"title"`
}

// NewBookDeletedEvent creates a new book deleted event
func NewBookDeletedEvent(isbn, title string) *BookDeletedEvent {
	return &BookDeletedEvent{
		BaseDomainEvent: NewBaseDomainEvent("BookDeleted", isbn, 1),
		ISBN:            isbn,
		Title:           title,
	}
}

// Author-related events

// AuthorCreatedEvent represents an author creation event
type AuthorCreatedEvent struct {
	BaseDomainEvent
	RG      string `json:"rg"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// NewAuthorCreatedEvent creates a new author created event
func NewAuthorCreatedEvent(rg, name, address string) *AuthorCreatedEvent {
	return &AuthorCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("AuthorCreated", rg, 1),
		RG:              rg,
		Name:            name,
		Address:         address,
	}
}

// AuthorUpdatedEvent represents an author update event
type AuthorUpdatedEvent struct {
	BaseDomainEvent
	RG      string                 `json:"rg"`
	Changes map[string]interface{} `json:"changes"`
}

// NewAuthorUpdatedEvent creates a new author updated event
func NewAuthorUpdatedEvent(rg string, changes map[string]interface{}) *AuthorUpdatedEvent {
	return &AuthorUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("AuthorUpdated", rg, 1),
		RG:              rg,
		Changes:         changes,
	}
}

// AuthorDeletedEvent represents an author deletion event
type AuthorDeletedEvent struct {
	BaseDomainEvent
	RG   string `json:"rg"`
	Name string `json:"name"`
}

// NewAuthorDeletedEvent creates a new author deleted event
func NewAuthorDeletedEvent(rg, name string) *AuthorDeletedEvent {
	return &AuthorDeletedEvent{
		BaseDomainEvent: NewBaseDomainEvent("AuthorDeleted", rg, 1),
		RG:              rg,
		Name:            name,
	}
}

// AuthorBookRelationship events

// AuthorAddedToBookEvent represents adding an author to a book
type AuthorAddedToBookEvent struct {
	BaseDomainEvent
	ISBN       string `json:"isbn"`
	BookTitle  string `json:"book_title"`
	AuthorRG   string `json:"author_rg"`
	AuthorName string `json:"author_name"`
}

// NewAuthorAddedToBookEvent creates a new author added to book event
func NewAuthorAddedToBookEvent(isbn, bookTitle, authorRG, authorName string) *AuthorAddedToBookEvent {
	return &AuthorAddedToBookEvent{
		BaseDomainEvent: NewBaseDomainEvent("AuthorAddedToBook", isbn+"-"+authorRG, 1),
		ISBN:            isbn,
		BookTitle:       bookTitle,
		AuthorRG:        authorRG,
		AuthorName:      authorName,
	}
}

// AuthorRemovedFromBookEvent represents removing an author from a book
type AuthorRemovedFromBookEvent struct {
	BaseDomainEvent
	ISBN       string `json:"isbn"`
	BookTitle  string `json:"book_title"`
	AuthorRG   string `json:"author_rg"`
	AuthorName string `json:"author_name"`
}

// NewAuthorRemovedFromBookEvent creates a new author removed from book event
func NewAuthorRemovedFromBookEvent(isbn, bookTitle, authorRG, authorName string) *AuthorRemovedFromBookEvent {
	return &AuthorRemovedFromBookEvent{
		BaseDomainEvent: NewBaseDomainEvent("AuthorRemovedFromBook", isbn+"-"+authorRG, 1),
		ISBN:            isbn,
		BookTitle:       bookTitle,
		AuthorRG:        authorRG,
		AuthorName:      authorName,
	}
}

// Publisher-related events

// PublisherCreatedEvent represents a publisher creation event
type PublisherCreatedEvent struct {
	BaseDomainEvent
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// NewPublisherCreatedEvent creates a new publisher created event
func NewPublisherCreatedEvent(id int, name, address string) *PublisherCreatedEvent {
	return &PublisherCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("PublisherCreated", string(rune(id)), 1),
		ID:              id,
		Name:            name,
		Address:         address,
	}
}

// PublisherUpdatedEvent represents a publisher update event
type PublisherUpdatedEvent struct {
	BaseDomainEvent
	ID      int                    `json:"id"`
	Changes map[string]interface{} `json:"changes"`
}

// NewPublisherUpdatedEvent creates a new publisher updated event
func NewPublisherUpdatedEvent(id int, changes map[string]interface{}) *PublisherUpdatedEvent {
	return &PublisherUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("PublisherUpdated", string(rune(id)), 1),
		ID:              id,
		Changes:         changes,
	}
}

// Printing company events

// PrintingCompanyCreatedEvent represents a printing company creation event
type PrintingCompanyCreatedEvent struct {
	BaseDomainEvent
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // "particular" or "contratada"
	Address     string `json:"address,omitempty"`
}

// NewPrintingCompanyCreatedEvent creates a new printing company created event
func NewPrintingCompanyCreatedEvent(id int, name, companyType, address string) *PrintingCompanyCreatedEvent {
	return &PrintingCompanyCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("PrintingCompanyCreated", string(rune(id)), 1),
		ID:              id,
		Name:            name,
		Type:            companyType,
		Address:         address,
	}
}

// PrintingCompanyUpdatedEvent represents a printing company update event
type PrintingCompanyUpdatedEvent struct {
	BaseDomainEvent
	ID      int                    `json:"id"`
	Changes map[string]interface{} `json:"changes"`
}

// NewPrintingCompanyUpdatedEvent creates a new printing company updated event
func NewPrintingCompanyUpdatedEvent(id int, changes map[string]interface{}) *PrintingCompanyUpdatedEvent {
	return &PrintingCompanyUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("PrintingCompanyUpdated", string(rune(id)), 1),
		ID:              id,
		Changes:         changes,
	}
}

// Contract events

// ContractCreatedEvent represents a contract creation event
type ContractCreatedEvent struct {
	BaseDomainEvent
	ID                 int     `json:"id"`
	Value              float64 `json:"value"`
	ResponsiblePerson  string  `json:"responsible_person"`
	PrintingCompanyID  int     `json:"printing_company_id"`
}

// NewContractCreatedEvent creates a new contract created event
func NewContractCreatedEvent(id int, value float64, responsiblePerson string, printingCompanyID int) *ContractCreatedEvent {
	return &ContractCreatedEvent{
		BaseDomainEvent:   NewBaseDomainEvent("ContractCreated", string(rune(id)), 1),
		ID:                id,
		Value:             value,
		ResponsiblePerson: responsiblePerson,
		PrintingCompanyID: printingCompanyID,
	}
}

// ContractUpdatedEvent represents a contract update event
type ContractUpdatedEvent struct {
	BaseDomainEvent
	ID      int                    `json:"id"`
	Changes map[string]interface{} `json:"changes"`
}

// NewContractUpdatedEvent creates a new contract updated event
func NewContractUpdatedEvent(id int, changes map[string]interface{}) *ContractUpdatedEvent {
	return &ContractUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ContractUpdated", string(rune(id)), 1),
		ID:              id,
		Changes:         changes,
	}
}

// Printing job events

// PrintingJobScheduledEvent represents a printing job scheduling event
type PrintingJobScheduledEvent struct {
	BaseDomainEvent
	ISBN              string    `json:"isbn"`
	BookTitle         string    `json:"book_title"`
	PrintingCompanyID int       `json:"printing_company_id"`
	CompanyName       string    `json:"company_name"`
	Copies            int       `json:"copies"`
	DeliveryDate      time.Time `json:"delivery_date"`
	ScheduledAt       time.Time `json:"scheduled_at"`
}

// NewPrintingJobScheduledEvent creates a new printing job scheduled event
func NewPrintingJobScheduledEvent(isbn, bookTitle string, printingCompanyID int, companyName string, copies int, deliveryDate time.Time) *PrintingJobScheduledEvent {
	return &PrintingJobScheduledEvent{
		BaseDomainEvent:   NewBaseDomainEvent("PrintingJobScheduled", isbn+"-"+string(rune(printingCompanyID)), 1),
		ISBN:              isbn,
		BookTitle:         bookTitle,
		PrintingCompanyID: printingCompanyID,
		CompanyName:       companyName,
		Copies:            copies,
		DeliveryDate:      deliveryDate,
		ScheduledAt:       time.Now(),
	}
}

// PrintingJobUpdatedEvent represents a printing job update event
type PrintingJobUpdatedEvent struct {
	BaseDomainEvent
	ISBN              string                 `json:"isbn"`
	PrintingCompanyID int                    `json:"printing_company_id"`
	Changes           map[string]interface{} `json:"changes"`
}

// NewPrintingJobUpdatedEvent creates a new printing job updated event
func NewPrintingJobUpdatedEvent(isbn string, printingCompanyID int, changes map[string]interface{}) *PrintingJobUpdatedEvent {
	return &PrintingJobUpdatedEvent{
		BaseDomainEvent:   NewBaseDomainEvent("PrintingJobUpdated", isbn+"-"+string(rune(printingCompanyID)), 1),
		ISBN:              isbn,
		PrintingCompanyID: printingCompanyID,
		Changes:           changes,
	}
}

// PrintingJobCompletedEvent represents a printing job completion event
type PrintingJobCompletedEvent struct {
	BaseDomainEvent
	ISBN              string    `json:"isbn"`
	BookTitle         string    `json:"book_title"`
	PrintingCompanyID int       `json:"printing_company_id"`
	CompanyName       string    `json:"company_name"`
	Copies            int       `json:"copies"`
	DeliveryDate      time.Time `json:"delivery_date"`
	CompletedAt       time.Time `json:"completed_at"`
	OnTime            bool      `json:"on_time"`
}

// NewPrintingJobCompletedEvent creates a new printing job completed event
func NewPrintingJobCompletedEvent(isbn, bookTitle string, printingCompanyID int, companyName string, copies int, deliveryDate time.Time, onTime bool) *PrintingJobCompletedEvent {
	return &PrintingJobCompletedEvent{
		BaseDomainEvent:   NewBaseDomainEvent("PrintingJobCompleted", isbn+"-"+string(rune(printingCompanyID)), 1),
		ISBN:              isbn,
		BookTitle:         bookTitle,
		PrintingCompanyID: printingCompanyID,
		CompanyName:       companyName,
		Copies:            copies,
		DeliveryDate:      deliveryDate,
		CompletedAt:       time.Now(),
		OnTime:            onTime,
	}
}

// PrintingJobOverdueEvent represents a printing job overdue event
type PrintingJobOverdueEvent struct {
	BaseDomainEvent
	ISBN              string    `json:"isbn"`
	BookTitle         string    `json:"book_title"`
	PrintingCompanyID int       `json:"printing_company_id"`
	CompanyName       string    `json:"company_name"`
	Copies            int       `json:"copies"`
	DeliveryDate      time.Time `json:"delivery_date"`
	DaysOverdue       int       `json:"days_overdue"`
}

// NewPrintingJobOverdueEvent creates a new printing job overdue event
func NewPrintingJobOverdueEvent(isbn, bookTitle string, printingCompanyID int, companyName string, copies int, deliveryDate time.Time, daysOverdue int) *PrintingJobOverdueEvent {
	return &PrintingJobOverdueEvent{
		BaseDomainEvent:   NewBaseDomainEvent("PrintingJobOverdue", isbn+"-"+string(rune(printingCompanyID)), 1),
		ISBN:              isbn,
		BookTitle:         bookTitle,
		PrintingCompanyID: printingCompanyID,
		CompanyName:       companyName,
		Copies:            copies,
		DeliveryDate:      deliveryDate,
		DaysOverdue:       daysOverdue,
	}
}

// EventDispatcher defines the interface for dispatching domain events
type EventDispatcher interface {
	Dispatch(event DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
}

// EventHandler defines the interface for handling domain events
type EventHandler interface {
	Handle(event DomainEvent) error
	CanHandle(eventType string) bool
}

// EventStore defines the interface for storing and retrieving domain events
type EventStore interface {
	Save(event DomainEvent) error
	GetEvents(aggregateID string) ([]DomainEvent, error)
	GetEventsByType(eventType string) ([]DomainEvent, error)
	GetEventsSince(since time.Time) ([]DomainEvent, error)
}

// AggregateRoot represents an aggregate root that can generate domain events
type AggregateRoot interface {
	GetUncommittedEvents() []DomainEvent
	CommitEvents()
	LoadFromHistory(events []DomainEvent) error
}

// BaseAggregateRoot provides common functionality for aggregate roots
type BaseAggregateRoot struct {
	uncommittedEvents []DomainEvent
	version           int
}

// AddEvent adds a domain event to the uncommitted events
func (ar *BaseAggregateRoot) AddEvent(event DomainEvent) {
	ar.uncommittedEvents = append(ar.uncommittedEvents, event)
}

// GetUncommittedEvents returns all uncommitted events
func (ar *BaseAggregateRoot) GetUncommittedEvents() []DomainEvent {
	return ar.uncommittedEvents
}

// CommitEvents clears the uncommitted events
func (ar *BaseAggregateRoot) CommitEvents() {
	ar.uncommittedEvents = nil
}

// GetVersion returns the current version
func (ar *BaseAggregateRoot) GetVersion() int {
	return ar.version
}

// IncrementVersion increments the version
func (ar *BaseAggregateRoot) IncrementVersion() {
	ar.version++
}