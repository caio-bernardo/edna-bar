package handlers

// EmbeddedSwaggerSpec contains the OpenAPI 3.0 specification as a string
// This serves as a fallback when the swagger.yaml file is not available
const EmbeddedSwaggerSpec = `
openapi: 3.0.3
info:
  title: Edna Bar Book Printing API
  description: |
    A comprehensive RESTful API for managing book printing operations.
    
    This API provides complete functionality for:
    - Managing books, authors, and publishers
    - Handling printing companies (private and contracted)
    - Creating and tracking printing contracts
    - Scheduling and monitoring printing jobs
    - Real-time delivery tracking
    
    Built with Go following clean architecture principles.
  version: 1.0.0
  contact:
    name: API Support
    url: https://github.com/edna-bar/api
    email: support@edna-bar.com
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: /api
    description: Current server

tags:
  - name: books
    description: Book management operations
  - name: authors
    description: Author management operations
  - name: publishers
    description: Publisher management operations
  - name: printing-companies
    description: Printing company management operations
  - name: contracts
    description: Contract management operations
  - name: printing-jobs
    description: Printing job management operations
  - name: system
    description: System health and information endpoints

paths:
  /livros:
    get:
      tags: [books]
      summary: List all books
      description: Retrieve a paginated list of all books in the system
      parameters:
        - name: limit
          in: query
          description: Maximum number of books to return
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
        - name: offset
          in: query
          description: Number of books to skip
          schema:
            type: integer
            minimum: 0
            default: 0
        - name: title
          in: query
          description: Filter books by title (partial match)
          schema:
            type: string
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Book'
                  total:
                    type: integer
                  limit:
                    type: integer
                  offset:
                    type: integer
        '500':
          $ref: '#/components/responses/InternalServerError'
    post:
      tags: [books]
      summary: Create a new book
      description: Add a new book to the system
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateBookRequest'
      responses:
        '201':
          description: Book created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Book'
        '400':
          $ref: '#/components/responses/BadRequest'
        '500':
          $ref: '#/components/responses/InternalServerError'

  /livros/{isbn}:
    get:
      tags: [books]
      summary: Get book by ISBN
      parameters:
        - name: isbn
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Book found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Book'
        '404':
          $ref: '#/components/responses/NotFound'
    put:
      tags: [books]
      summary: Update book
      parameters:
        - name: isbn
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateBookRequest'
      responses:
        '200':
          description: Book updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Book'
        '400':
          $ref: '#/components/responses/BadRequest'
        '404':
          $ref: '#/components/responses/NotFound'
    delete:
      tags: [books]
      summary: Delete book
      parameters:
        - name: isbn
          in: path
          required: true
          schema:
            type: string
      responses:
        '204':
          description: Book deleted successfully
        '404':
          $ref: '#/components/responses/NotFound'

  /livros/{isbn}/authors:
    get:
      tags: [books]
      summary: Get book authors
      parameters:
        - name: isbn
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Authors found
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Author'
    post:
      tags: [books]
      summary: Add author to book
      parameters:
        - name: isbn
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                rg:
                  type: string
              required: [rg]
      responses:
        '201':
          description: Author added to book successfully

  /autores:
    get:
      tags: [authors]
      summary: List all authors
      parameters:
        - name: name
          in: query
          description: Filter authors by name
          schema:
            type: string
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Author'
    post:
      tags: [authors]
      summary: Create a new author
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateAuthorRequest'
      responses:
        '201':
          description: Author created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Author'

  /autores/{rg}:
    get:
      tags: [authors]
      summary: Get author by RG
      parameters:
        - name: rg
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Author found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Author'
        '404':
          $ref: '#/components/responses/NotFound'
    put:
      tags: [authors]
      summary: Update author
      parameters:
        - name: rg
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateAuthorRequest'
      responses:
        '200':
          description: Author updated successfully
    delete:
      tags: [authors]
      summary: Delete author
      parameters:
        - name: rg
          in: path
          required: true
          schema:
            type: string
      responses:
        '204':
          description: Author deleted successfully

  /autores/{rg}/books:
    get:
      tags: [authors]
      summary: Get author's books
      parameters:
        - name: rg
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Books found
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Book'

  /editoras:
    get:
      tags: [publishers]
      summary: List all publishers
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Publisher'
    post:
      tags: [publishers]
      summary: Create a new publisher
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreatePublisherRequest'
      responses:
        '201':
          description: Publisher created successfully

  /editoras/{id}:
    get:
      tags: [publishers]
      summary: Get publisher by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Publisher found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Publisher'
    put:
      tags: [publishers]
      summary: Update publisher
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdatePublisherRequest'
      responses:
        '200':
          description: Publisher updated successfully
    delete:
      tags: [publishers]
      summary: Delete publisher
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '204':
          description: Publisher deleted successfully

  /graficas:
    get:
      tags: [printing-companies]
      summary: List all printing companies
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/PrintingCompany'
    post:
      tags: [printing-companies]
      summary: Create a new printing company
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreatePrintingCompanyRequest'
      responses:
        '201':
          description: Printing company created successfully

  /graficas/{id}:
    get:
      tags: [printing-companies]
      summary: Get printing company by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Printing company found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PrintingCompany'

  /contratos:
    get:
      tags: [contracts]
      summary: List all contracts
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Contract'
    post:
      tags: [contracts]
      summary: Create a new contract
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateContractRequest'
      responses:
        '201':
          description: Contract created successfully

  /contratos/{id}:
    get:
      tags: [contracts]
      summary: Get contract by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: Contract found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Contract'

  /printing-jobs:
    get:
      tags: [printing-jobs]
      summary: List all printing jobs
      parameters:
        - name: status
          in: query
          description: Filter by job status
          schema:
            type: string
            enum: [pending, in_progress, completed, overdue]
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/PrintingJob'
    post:
      tags: [printing-jobs]
      summary: Create a new printing job
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreatePrintingJobRequest'
      responses:
        '201':
          description: Printing job created successfully

  /printing-jobs/overdue:
    get:
      tags: [printing-jobs]
      summary: Get overdue printing jobs
      responses:
        '200':
          description: Overdue jobs found
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/PrintingJob'

  /:
    get:
      tags: [system]
      summary: API information
      responses:
        '200':
          description: API information
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/APIInfo'

components:
  schemas:
    Book:
      type: object
      properties:
        isbn:
          type: string
          description: Book ISBN
          example: "978-0307474728"
        titulo:
          type: string
          description: Book title
          example: "Cem Anos de Solidão"
        data_de_publicacao:
          type: string
          format: date
          description: Publication date
          example: "1967-06-05"
        editora_id:
          type: integer
          description: Publisher ID
          example: 1
      required: [isbn, titulo, data_de_publicacao, editora_id]

    CreateBookRequest:
      type: object
      properties:
        isbn:
          type: string
          description: Book ISBN
          example: "978-0307474728"
        titulo:
          type: string
          description: Book title
          example: "Cem Anos de Solidão"
        data_de_publicacao:
          type: string
          format: date
          description: Publication date
          example: "1967-06-05"
        editora_id:
          type: integer
          description: Publisher ID
          example: 1
      required: [isbn, titulo, data_de_publicacao, editora_id]

    UpdateBookRequest:
      type: object
      properties:
        titulo:
          type: string
          description: Book title
          example: "Cem Anos de Solidão - Edição Especial"
        data_de_publicacao:
          type: string
          format: date
          description: Publication date
          example: "1967-06-05"
        editora_id:
          type: integer
          description: Publisher ID
          example: 1

    Author:
      type: object
      properties:
        rg:
          type: string
          description: Author RG
          example: "12345678901"
        nome:
          type: string
          description: Author name
          example: "Gabriel García Márquez"
        endereco:
          type: string
          description: Author address
          example: "América Latina"
      required: [rg, nome, endereco]

    CreateAuthorRequest:
      type: object
      properties:
        rg:
          type: string
          description: Author RG
          example: "12345678901"
        nome:
          type: string
          description: Author name
          example: "Gabriel García Márquez"
        endereco:
          type: string
          description: Author address
          example: "América Latina"
      required: [rg, nome, endereco]

    UpdateAuthorRequest:
      type: object
      properties:
        nome:
          type: string
          description: Author name
          example: "Gabriel García Márquez"
        endereco:
          type: string
          description: Author address
          example: "América Latina"

    Publisher:
      type: object
      properties:
        id:
          type: integer
          description: Publisher ID
          example: 1
        nome:
          type: string
          description: Publisher name
          example: "Penguin Random House Brasil"
        endereco:
          type: string
          description: Publisher address
          example: "São Paulo, SP, Brasil"
      required: [id, nome, endereco]

    CreatePublisherRequest:
      type: object
      properties:
        nome:
          type: string
          description: Publisher name
          example: "Penguin Random House Brasil"
        endereco:
          type: string
          description: Publisher address
          example: "São Paulo, SP, Brasil"
      required: [nome, endereco]

    UpdatePublisherRequest:
      type: object
      properties:
        nome:
          type: string
          description: Publisher name
          example: "Penguin Random House Brasil"
        endereco:
          type: string
          description: Publisher address
          example: "São Paulo, SP, Brasil"

    PrintingCompany:
      type: object
      properties:
        id:
          type: integer
          description: Printing company ID
          example: 1
        nome:
          type: string
          description: Printing company name
          example: "PrintTech Solutions"
        type:
          type: string
          enum: [particular, contratada]
          description: Type of printing company
          example: "contratada"
        endereco:
          type: string
          description: Address (for contracted companies)
          example: "Distrito Industrial, São Paulo, SP"
      required: [id, nome, type]

    CreatePrintingCompanyRequest:
      type: object
      properties:
        nome:
          type: string
          description: Printing company name
          example: "PrintTech Solutions"
        type:
          type: string
          enum: [particular, contratada]
          description: Type of printing company
          example: "contratada"
        endereco:
          type: string
          description: Address (required for contracted companies)
          example: "Distrito Industrial, São Paulo, SP"
      required: [nome, type]

    Contract:
      type: object
      properties:
        id:
          type: integer
          description: Contract ID
          example: 1
        valor:
          type: number
          format: float
          description: Contract value
          example: 15000.50
        nome_responsavel:
          type: string
          description: Responsible person name
          example: "Maria Silva"
        grafica_cont_id:
          type: integer
          description: Contracted printing company ID
          example: 1
      required: [id, valor, nome_responsavel, grafica_cont_id]

    CreateContractRequest:
      type: object
      properties:
        valor:
          type: number
          format: float
          description: Contract value
          example: 15000.50
        nome_responsavel:
          type: string
          description: Responsible person name
          example: "Maria Silva"
        grafica_cont_id:
          type: integer
          description: Contracted printing company ID
          example: 1
      required: [valor, nome_responsavel, grafica_cont_id]

    PrintingJob:
      type: object
      properties:
        isbn:
          type: string
          description: Book ISBN
          example: "978-0307474728"
        grafica_id:
          type: integer
          description: Printing company ID
          example: 1
        nto_copias:
          type: integer
          description: Number of copies to print
          example: 5000
        data_entrega:
          type: string
          format: date
          description: Expected delivery date
          example: "2024-02-15"
        status:
          type: string
          enum: [pending, in_progress, completed, overdue]
          description: Job status
          example: "pending"
      required: [isbn, grafica_id, nto_copias, data_entrega]

    CreatePrintingJobRequest:
      type: object
      properties:
        isbn:
          type: string
          description: Book ISBN
          example: "978-0307474728"
        grafica_id:
          type: integer
          description: Printing company ID
          example: 1
        nto_copias:
          type: integer
          description: Number of copies to print
          example: 5000
        data_entrega:
          type: string
          format: date
          description: Expected delivery date
          example: "2024-02-15"
      required: [isbn, grafica_id, nto_copias, data_entrega]

    APIInfo:
      type: object
      properties:
        name:
          type: string
          example: "Edna Bar Book Printing API"
        version:
          type: string
          example: "1.0.0"
        description:
          type: string
          example: "RESTful API for managing book printing operations"

    Error:
      type: object
      properties:
        error:
          type: string
          description: Error message
          example: "Resource not found"
        code:
          type: string
          description: Error code
          example: "NOT_FOUND"
        details:
          type: string
          description: Additional error details
          example: "Book with ISBN 978-0123456789 was not found"
        timestamp:
          type: string
          format: date-time
          description: Error timestamp
          example: "2024-01-15T10:30:00Z"
      required: [error, code]

  responses:
    BadRequest:
      description: Bad request - invalid input
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'

    NotFound:
      description: Resource not found
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'

    InternalServerError:
      description: Internal server error
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/Error'
`
