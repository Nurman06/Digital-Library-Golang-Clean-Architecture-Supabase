# Digital Library Management System

A modern digital library management system built with Go, Clean Architecture, and Supabase.

## Features

- **Book Catalog Management**: Complete CRUD operations for books with metadata
- **User Authentication & Authorization**: Secure JWT-based authentication with role-based access control
- **Borrowing/Lending Operations**: Book checkout, return, and renewal functionality
- **Borrowing History**: Complete audit trail of all borrowing transactions
- **Search & Discovery**: Advanced search functionality for finding books
- **Availability Tracking**: Real-time tracking of book availability status

## Architecture

This project follows Clean Architecture principles with clear separation of concerns:

```
internal/
├── entity/              # Domain models (Book, User, BorrowRecord)
├── usecase/             # Business logic
├── adapter/
│   ├── handler/         # HTTP handlers
│   └── repository/      # Data access implementations
└── infrastructure/
    ├── supabase/        # Supabase client
    ├── server/          # HTTP server
    ├── config/          # Configuration management
    └── logger/          # Logging infrastructure
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+ (via Supabase)
- Supabase account

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase.git
cd Digital-Library-Golang-Clean-Architecture-Supabase
```

2. Copy the environment file and configure it:
```bash
cp .env.example .env
```

3. Update the `.env` file with your Supabase credentials:
```env
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your_anon_key_here
SUPABASE_SERVICE_KEY=your_service_key_here
JWT_SECRET=your_jwt_secret_key_here
```

4. Install dependencies:
```bash
go mod download
```

5. Run database migrations:
```bash
# Connect to your Supabase database and run the migration
psql -h your-db-host -U postgres -d library_db -f migrations/001_create_tables.sql
```

## Running the Application

### Development Mode

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080` by default.

### Production Build

```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

## API Endpoints

### Health Check
- `GET /health` - Check API health status

### API v1
- `GET /api/v1/` - API information

### Books (Coming Soon)
- `GET /api/v1/books` - List all books
- `GET /api/v1/books/:id` - Get book details
- `POST /api/v1/books` - Create new book (Admin/Librarian)
- `PUT /api/v1/books/:id` - Update book (Admin/Librarian)
- `DELETE /api/v1/books/:id` - Delete book (Admin)

### Authentication (Coming Soon)
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login

### Borrowing (Coming Soon)
- `POST /api/v1/borrowing/checkout` - Borrow a book
- `POST /api/v1/borrowing/return` - Return a book
- `POST /api/v1/borrowing/renew` - Renew a book
- `GET /api/v1/borrowing/history` - Get borrowing history

## Configuration

The application can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | `8080` |
| `SERVER_HOST` | Server host | `0.0.0.0` |
| `ENV` | Environment (development/production) | `development` |
| `SUPABASE_URL` | Supabase project URL | Required |
| `SUPABASE_ANON_KEY` | Supabase anonymous key | Required |
| `SUPABASE_SERVICE_KEY` | Supabase service key | Required |
| `JWT_SECRET` | JWT signing secret | Required in production |
| `JWT_EXPIRY_HOURS` | JWT token expiry time | `24` |
| `DB_MAX_OPEN_CONNS` | Maximum open database connections | `25` |
| `DB_MAX_IDLE_CONNS` | Maximum idle database connections | `5` |

## Database Schema

The application uses PostgreSQL with the following main tables:

- **books**: Book catalog information
- **book_copies**: Individual book copies
- **users**: User accounts
- **borrow_records**: Borrowing transactions
- **reservations**: Book reservation queue

See `migrations/001_create_tables.sql` for the complete schema.

## Development

### Project Structure

```
.
├── cmd/
│   └── api/                # Application entry point
├── internal/
│   ├── entity/             # Domain entities
│   ├── usecase/            # Business logic
│   ├── adapter/
│   │   ├── handler/        # HTTP handlers
│   │   └── repository/     # Data repositories
│   └── infrastructure/
│       ├── supabase/       # Supabase client
│       ├── server/         # HTTP server
│       ├── config/         # Configuration
│       └── logger/         # Logging
├── migrations/             # Database migrations
└── openspec/              # OpenSpec documentation
```

### Running Tests

```bash
go test ./...
```

### Code Quality

```bash
# Run linter
golangci-lint run

# Format code
go fmt ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Contact

Nurman - [@Nurman06](https://github.com/Nurman06)

Project Link: [https://github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase](https://github.com/Nurman06/Digital-Library-Golang-Clean-Architecture-Supabase)
