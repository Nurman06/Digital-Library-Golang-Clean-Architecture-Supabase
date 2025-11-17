# Technical Design Document

## Context

This document outlines the technical design decisions for the Digital Library Management System. The system is a greenfield implementation using Go, Clean Architecture, and Supabase as the backend infrastructure. This design aims to create a scalable, maintainable, and secure library management system.

### Stakeholders
- **Library Staff**: Librarians and administrators who manage books and users
- **Library Members**: Users who borrow and search for books
- **System Administrators**: Technical team maintaining the infrastructure
- **Developers**: Team implementing and maintaining the codebase

### Constraints
- Must use Supabase for data persistence and authentication
- Go 1.21+ for backend implementation
- RESTful API design principles
- Response time < 500ms for 95th percentile
- Support horizontal scaling

## Goals / Non-Goals

### Goals
- ✅ Implement Clean Architecture with clear layer boundaries
- ✅ Provide secure authentication and role-based authorization
- ✅ Enable efficient book search and discovery
- ✅ Track borrowing operations with complete audit trail
- ✅ Support real-time availability tracking
- ✅ Achieve >80% test coverage for business logic
- ✅ Design for horizontal scalability
- ✅ Implement comprehensive error handling

### Non-Goals
- ❌ Mobile native applications (API-first, clients can be built separately)
- ❌ Real-time notifications (can be added in future iteration)
- ❌ Payment processing integration (late fees tracked but not processed)
- ❌ Multi-language support in initial version
- ❌ Advanced analytics and reporting (basic statistics only)

## Technical Decisions

### 1. Architecture Pattern: Clean Architecture

**Decision**: Use Clean Architecture with four distinct layers

**Rationale**:
- Clear separation of concerns
- Framework independence (easy to swap Supabase if needed)
- Testability (business logic independent of infrastructure)
- Maintainability (changes isolated to specific layers)

**Layer Structure**:
```
internal/
├── entity/              # Domain models (Book, User, BorrowRecord)
├── usecase/             # Business logic (BookUseCase, BorrowingUseCase)
├── adapter/
│   ├── handler/         # HTTP handlers
│   └── repository/      # Data access implementations
└── infrastructure/
    ├── supabase/        # Supabase client
    ├── server/          # HTTP server
    └── config/          # Configuration management
```

**Alternatives Considered**:
- MVC Pattern: Rejected due to tight coupling between layers
- Hexagonal Architecture: Similar benefits but Clean Architecture more familiar to Go community

### 2. Database Schema Design

**Decision**: Normalized relational schema with proper foreign keys and indexes

**Schema Overview**:

```sql
-- Books table
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(13) UNIQUE NOT NULL,
    category VARCHAR(100),
    publication_year INTEGER,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL  -- Soft delete
);

-- Book copies (physical/digital copies)
CREATE TABLE book_copies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    copy_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL,  -- available, borrowed, reserved, lost, damaged
    location VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Users table (managed by Supabase Auth, extended here)
CREATE TABLE users (
    id UUID PRIMARY KEY REFERENCES auth.users(id),
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,  -- admin, librarian, member
    status VARCHAR(20) NOT NULL,  -- active, suspended, expired
    borrowing_limit INTEGER DEFAULT 5,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Borrow records
CREATE TABLE borrow_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE RESTRICT,
    book_copy_id UUID REFERENCES book_copies(id) ON DELETE RESTRICT,
    checkout_date TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP NULL,
    renewal_count INTEGER DEFAULT 0,
    late_fee DECIMAL(10,2) DEFAULT 0.00,
    status VARCHAR(20) NOT NULL,  -- active, returned, overdue
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Reservations
CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    reservation_date TIMESTAMP DEFAULT NOW(),
    expiry_date TIMESTAMP,
    status VARCHAR(20) NOT NULL,  -- pending, fulfilled, expired, cancelled
    queue_position INTEGER NOT NULL
);
```

**Indexes for Performance**:
```sql
-- Book search optimization
CREATE INDEX idx_books_title ON books USING GIN (to_tsvector('english', title));
CREATE INDEX idx_books_author ON books USING GIN (to_tsvector('english', author));
CREATE INDEX idx_books_isbn ON books(isbn);
CREATE INDEX idx_books_category ON books(category);

-- Borrowing operations
CREATE INDEX idx_borrow_records_user_id ON borrow_records(user_id);
CREATE INDEX idx_borrow_records_status ON borrow_records(status);
CREATE INDEX idx_borrow_records_due_date ON borrow_records(due_date) WHERE status = 'active';

-- Availability tracking
CREATE INDEX idx_book_copies_book_id ON book_copies(book_id);
CREATE INDEX idx_book_copies_status ON book_copies(status);
```

**Rationale**:
- UUID for primary keys (distributed system friendly)
- Soft deletes for audit trail
- Proper foreign keys for data integrity
- GIN indexes for full-text search
- Composite indexes for common query patterns

### 3. API Design and Versioning

**Decision**: RESTful API with URL-based versioning

**API Structure**:
```
/api/v1/
├── /auth
│   ├── POST /register
│   ├── POST /login
│   └── POST /refresh
├── /books
│   ├── GET    /books
│   ├── GET    /books/:id
│   ├── POST   /books
│   ├── PUT    /books/:id
│   ├── DELETE /books/:id
│   └── GET    /books/search
├── /borrowing
│   ├── POST /borrowing/checkout
│   ├── POST /borrowing/return
│   ├── POST /borrowing/renew
│   └── GET  /borrowing/history
├── /users
│   ├── GET    /users
│   ├── GET    /users/:id
│   ├── PUT    /users/:id
│   └── DELETE /users/:id
└── /availability
    ├── GET /availability/:bookId
    └── GET /availability/search
```

**Response Format**:
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "page": 1,
    "pageSize": 20,
    "total": 100
  },
  "timestamp": "2025-01-01T00:00:00Z"
}
```

**Error Format**:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "isbn",
        "message": "ISBN format is invalid"
      }
    ]
  },
  "timestamp": "2025-01-01T00:00:00Z"
}
```

**Rationale**:
- URL versioning (simple, cache-friendly)
- Consistent response structure
- Detailed error messages for debugging
- Pagination metadata for large datasets

### 4. Caching Strategy

**Decision**: Multi-layer caching with Redis (optional) and in-memory cache

**Caching Layers**:

1. **In-Memory Cache** (for frequently accessed data):
   - Book catalog (TTL: 5 minutes)
   - User roles and permissions (TTL: 10 minutes)
   - Available book count (TTL: 1 minute)

2. **Redis Cache** (optional, for distributed caching):
   - Search results (TTL: 2 minutes)
   - User session data (TTL: 24 hours)
   - Popular books list (TTL: 1 hour)

**Cache Invalidation Strategy**:
- Write-through for critical data (availability status)
- Time-based expiration for read-heavy data
- Event-based invalidation for updates

**Implementation**:
```go
type CacheService interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    Clear() error
}
```

**Rationale**:
- Reduce database load for read-heavy operations
- Improve response times for common queries
- Optional Redis for horizontal scaling

### 5. Authentication and Authorization

**Decision**: JWT-based authentication with Supabase Auth integration

**JWT Token Structure**:
```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "role": "member",
  "exp": 1735689600,
  "iat": 1735603200
}
```

**Authorization Flow**:
1. User logs in via Supabase Auth
2. System generates JWT with user role
3. Each request includes JWT in Authorization header
4. Middleware validates JWT and extracts user context
5. Handler checks role-based permissions

**Role Permissions Matrix**:

| Operation | Admin | Librarian | Member |
|-----------|-------|-----------|--------|
| Create Book | ✅ | ✅ | ❌ |
| Update Book | ✅ | ✅ | ❌ |
| Delete Book | ✅ | ❌ | ❌ |
| View Books | ✅ | ✅ | ✅ |
| Borrow Book | ✅ | ✅ | ✅ |
| Manage Users | ✅ | ❌ | ❌ |
| Override Limits | ✅ | ✅ | ❌ |

**Rationale**:
- Supabase Auth handles password security
- JWT stateless (no session storage needed)
- Role-based permissions easy to extend

### 6. Error Handling Strategy

**Decision**: Layered error handling with custom error types

**Error Types**:
```go
type DomainError struct {
    Code    string
    Message string
    Details map[string]interface{}
}

// Error categories
const (
    ErrValidation    = "VALIDATION_ERROR"
    ErrNotFound      = "NOT_FOUND"
    ErrConflict      = "CONFLICT"
    ErrUnauthorized  = "UNAUTHORIZED"
    ErrForbidden     = "FORBIDDEN"
    ErrInternal      = "INTERNAL_ERROR"
)
```

**Error Propagation**:
- Entity layer: Domain validation errors
- Use case layer: Business logic errors
- Repository layer: Data access errors
- Handler layer: HTTP status code mapping

**Rationale**:
- Clear error categorization
- Consistent error responses
- Easy debugging and monitoring

### 7. Performance Optimization

**Decision**: Database query optimization and connection pooling

**Optimization Techniques**:

1. **Database Connection Pool**:
   ```go
   MaxOpenConns: 25
   MaxIdleConns: 5
   ConnMaxLifetime: 5 minutes
   ```

2. **Query Optimization**:
   - Use prepared statements
   - Batch inserts for bulk operations
   - Selective column retrieval
   - Proper JOIN usage

3. **Pagination**:
   - Cursor-based for large datasets
   - Limit + Offset for smaller datasets

4. **Async Operations**:
   - Background jobs for notifications
   - Async logging

**Performance Targets**:
- API response time: < 200ms (p95)
- Database query time: < 50ms (p95)
- Concurrent users: 1000+
- Throughput: 500 req/sec

## Risks / Trade-offs

### Risk 1: Supabase Vendor Lock-in
**Mitigation**: 
- Abstract database access behind repository interfaces
- Use standard PostgreSQL features
- Maintain migration scripts

### Risk 2: Performance Degradation with Scale
**Mitigation**:
- Implement caching strategy
- Database query optimization
- Horizontal scaling capability
- Regular performance testing

### Risk 3: Data Consistency in Distributed Cache
**Mitigation**:
- Write-through caching for critical data
- Event-based cache invalidation
- Short TTL for volatile data

### Risk 4: Security Vulnerabilities
**Mitigation**:
- Input validation at all layers
- Parameterized queries (SQL injection prevention)
- Rate limiting
- Regular security audits
- HTTPS only in production

## Migration Plan

**Note**: This is a greenfield implementation, so no data migration is required. However, for future reference:

### Database Migrations
- Use migration tool (e.g., golang-migrate)
- Version-controlled migration scripts
- Rollback capability for each migration
- Test migrations in staging before production

### Zero-Downtime Deployment
1. Deploy new version alongside old version
2. Route percentage of traffic to new version
3. Monitor error rates and performance
4. Gradually increase traffic to new version
5. Decommission old version

## Open Questions

1. **Email Notifications**: Should we implement email notifications for due dates and overdue books in this iteration or defer to v2?
   - **Recommendation**: Defer to v2, focus on core functionality first

2. **Book Cover Images**: Should we support book cover image uploads via Supabase Storage?
   - **Recommendation**: Include basic support, optional field

3. **Multi-tenancy**: Should the system support multiple library branches from day one?
   - **Recommendation**: Design for single library, but structure allows future multi-tenancy

4. **API Rate Limiting**: What are the appropriate rate limits per user role?
   - **Recommendation**: Member: 100 req/min, Librarian: 200 req/min, Admin: unlimited

5. **Data Retention**: How long should we retain borrowing history and audit logs?
   - **Recommendation**: 7 years for compliance, with archival strategy

## Performance Benchmarks

### Target Metrics

**API Response Times** (95th percentile):
- Book search: < 150ms
- Book detail: < 100ms
- Checkout operation: < 200ms
- User authentication: < 150ms
- List operations: < 200ms

**Database Query Times** (95th percentile):
- Simple SELECT: < 20ms
- Complex JOIN: < 50ms
- Full-text search: < 100ms
- Write operations: < 30ms

**System Capacity**:
- Concurrent users: 1000+
- Requests per second: 500+
- Database connections: 25 max
- Memory usage: < 512MB per instance

**Availability**:
- Uptime: 99.9% (8.76 hours downtime/year)
- Error rate: < 0.1%
- Success rate: > 99.9%

### Monitoring and Alerting

**Key Metrics to Monitor**:
- API response times (p50, p95, p99)
- Error rates by endpoint
- Database query performance
- Cache hit/miss ratios
- Active user sessions
- Borrowing transaction rates

**Alert Thresholds**:
- Response time > 500ms for 5 minutes
- Error rate > 1% for 5 minutes
- Database connection pool > 80% for 10 minutes
- Disk space > 85%
- Memory usage > 90%