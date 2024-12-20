# DnD 5e API (Golang)

This project provides a RESTful API to manage DnD 5e entities like Races, Traits, Subraces, Languages, and Proficiencies. It follows Domain-Driven Design (DDD), SOLID principles, and standard Go project conventions to ensure clean architecture, scalability, and maintainability.

## Key Features

- **Entities**: Defines core DnD 5e entities (Race, Trait, Subrace, Language, Proficiency) with structured relationships.
- **Domain-Driven Design**: Separates domain logic from infrastructure and interfaces for a clear and extensible architecture.
- **SOLID Principles**: Encourages code maintainability, testability, and clear separation of concerns.
- **Design Patterns**: Employs repository interfaces, use cases, and services to decouple business logic from persistence details.
- **GORM/SQL Integration**: Uses GORM with a PostgreSQL backend for data persistence and automatic handling of relationships.
- **HTTP API**: Exposes RESTful endpoints for creating, retrieving, updating, and deleting DnD 5e domain objects.
- **Scalable & Extensible**: Easily add new features, entities, or change the storage layer without impacting the domain logic.

## Project Structure (Overview)

- **`cmd/`**: Entry point for the application executable.
- **`internal/domain/race/`**: Entities, repositories (interfaces), services, and use cases related to race domain logic.
- **`internal/infrastructure/persistence/gorm/`**: Concrete repository implementations using GORM and PostgreSQL.
- **`internal/interfaces/api/`**: HTTP handlers, routers, and middleware.
- **`internal/infrastructure/db/`**: Database connection and initialization.
- **`pkg/`**: Shared utilities or libraries.

## Getting Started

- Clone the repository
- Set up PostgreSQL and configure environment variables
- Run `go build` and `./your_project` to start the server
- Access the API endpoints (e.g., `GET /api/races`) to manage DnD 5e data.

## Next Steps

- Implement additional endpoints or domain logic.
- Add authentication and authorization if required.
- Integrate with a front-end (e.g., Angular) for a complete application experience.
