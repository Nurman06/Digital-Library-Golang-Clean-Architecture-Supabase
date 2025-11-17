# Implementation Tasks

## 1. Project Setup and Infrastructure
- [x] 1.1 Initialize Go module and project structure
- [x] 1.2 Set up directory structure following Clean Architecture (entity, usecase, adapter, infrastructure)
- [x] 1.3 Configure Supabase client and connection
- [x] 1.4 Create environment configuration management
- [x] 1.5 Set up logging infrastructure
- [x] 1.6 Configure HTTP server and routing

## 2. Database Schema and Migrations
- [x] 2.1 Design and create books table schema
- [x] 2.2 Design and create users table schema
- [x] 2.3 Design and create borrow_records table schema
- [x] 2.4 Design and create book_copies table schema
- [x] 2.5 Create database indexes for performance
- [ ] 2.6 Set up Supabase Row Level Security (RLS) policies (NOT REQUESTED - deferred)
- [x] 2.7 Create database migration scripts (including rollback scripts and documentation)

## 3. Entity Layer (Domain Models)
- [x] 3.1 Implement Book entity with validation
- [x] 3.2 Implement User entity with role definitions
- [x] 3.3 Implement BorrowRecord entity
- [x] 3.4 Implement BookCopy entity
- [ ] 3.5 Define domain error types
- [x] 3.6 Implement value objects (ISBN, Email, etc.)

## 4. Repository Interfaces
- [ ] 4.1 Define BookRepository interface
- [ ] 4.2 Define UserRepository interface
- [ ] 4.3 Define BorrowRecordRepository interface
- [ ] 4.4 Define BookCopyRepository interface

## 5. Use Case Layer (Business Logic)
- [ ] 5.1 Implement BookUseCase (CRUD operations)
- [ ] 5.2 Implement UserUseCase (registration, profile management)
- [ ] 5.3 Implement AuthUseCase (login, token validation)
- [ ] 5.4 Implement BorrowingUseCase (checkout, return, renewal)
- [ ] 5.5 Implement SearchUseCase (book search by various criteria)
- [ ] 5.6 Implement AvailabilityUseCase (real-time availability tracking)
- [ ] 5.7 Implement late fee calculation logic (logic sudah ada di entity, perlu di use case)
- [ ] 5.8 Implement borrowing limit validation (logic sudah ada di entity, perlu di use case)
- [ ] 5.9 Implement due date calculation logic (konstanta sudah ada di entity, perlu di use case)

## 6. Repository Implementations
- [ ] 6.1 Implement PostgresBookRepository
- [ ] 6.2 Implement PostgresUserRepository
- [ ] 6.3 Implement PostgresBorrowRecordRepository
- [ ] 6.4 Implement PostgresBookCopyRepository
- [ ] 6.5 Add error handling and transaction management

## 7. HTTP Handlers (API Endpoints)
- [ ] 7.1 Implement book CRUD endpoints (POST, GET, PUT, DELETE /books)
- [ ] 7.2 Implement user registration endpoint (POST /auth/register)
- [ ] 7.3 Implement login endpoint (POST /auth/login)
- [ ] 7.4 Implement user profile endpoints (GET, PUT /users/{id})
- [ ] 7.5 Implement borrow book endpoint (POST /borrowing/checkout)
- [ ] 7.6 Implement return book endpoint (POST /borrowing/return)
- [ ] 7.7 Implement renew book endpoint (POST /borrowing/renew)
- [ ] 7.8 Implement borrowing history endpoint (GET /borrowing/history)
- [ ] 7.9 Implement book search endpoint (GET /books/search)
- [ ] 7.10 Implement availability check endpoint (GET /books/{id}/availability)
- [ ] 7.11 Add request validation middleware
- [ ] 7.12 Add authentication middleware
- [ ] 7.13 Add authorization middleware (role-based)

## 8. Authentication & Authorization
- [ ] 8.1 Integrate Supabase Auth for JWT token generation
- [ ] 8.2 Implement JWT token validation middleware
- [ ] 8.3 Implement role-based access control (RBAC)
- [ ] 8.4 Add password hashing and validation
- [ ] 8.5 Implement token refresh mechanism

## 9. Testing
- [ ] 9.1 Write unit tests for Book entity
- [ ] 9.2 Write unit tests for User entity
- [ ] 9.3 Write unit tests for BorrowRecord entity
- [ ] 9.4 Write unit tests for BookUseCase
- [ ] 9.5 Write unit tests for BorrowingUseCase
- [ ] 9.6 Write unit tests for SearchUseCase
- [ ] 9.7 Write integration tests for repository implementations
- [ ] 9.8 Write API endpoint tests
- [ ] 9.9 Set up test database and fixtures
- [ ] 9.10 Achieve >80% code coverage for business logic

## 10. Documentation
- [ ] 10.1 Write API documentation (endpoints, request/response formats)
- [x] 10.2 Document database schema
- [x] 10.3 Create setup and installation guide
- [x] 10.4 Document environment variables and configuration
- [x] 10.5 Add code comments and package documentation
- [ ] 10.6 Create example requests and responses

## 11. Error Handling and Validation
- [ ] 11.1 Implement consistent error response format
- [ ] 11.2 Add input validation for all endpoints
- [ ] 11.3 Implement proper HTTP status codes
- [ ] 11.4 Add error logging and monitoring

## 12. Performance Optimization
- [x] 12.1 Implement database connection pooling (MaxOpenConns: 25, MaxIdleConns: 5)
- [ ] 12.2 Add database query optimization with prepared statements
- [x] 12.3 Create database indexes for common query patterns
- [ ] 12.4 Implement caching layer (in-memory cache with TTL)
- [ ] 12.5 Add pagination for large result sets (cursor-based and offset-based)
- [x] 12.6 Optimize full-text search queries with GIN indexes
- [ ] 12.7 Implement batch operations for bulk inserts
- [ ] 12.8 Add query performance monitoring and logging
- [ ] 12.9 Set up performance benchmarks and load testing
- [ ] 12.10 Verify API response times meet targets (< 200ms p95)

## 13. Security Implementation
- [ ] 13.1 Add input sanitization to prevent injection attacks
- [ ] 13.2 Implement API rate limiting (Member: 100 req/min, Librarian: 200 req/min)
- [ ] 13.3 Configure CORS policies for production
- [ ] 13.4 Add request/response logging with sensitive data masking
- [ ] 13.5 Implement HTTPS-only in production
- [ ] 13.6 Add security headers (HSTS, CSP, X-Frame-Options)
- [x] 13.7 Implement SQL injection prevention with parameterized queries
- [ ] 13.8 Add XSS protection in API responses
- [x] 13.9 Implement graceful shutdown for zero-downtime deployments
- [ ] 13.10 Set up security audit logging

## 14. Monitoring and Observability
- [x] 14.1 Implement health check endpoint (GET /health)
- [ ] 14.2 Add readiness and liveness probes
- [x] 14.3 Set up structured logging with log levels
- [ ] 14.4 Implement metrics collection (Prometheus format)
- [ ] 14.5 Add performance metrics (response time, throughput, error rate)
- [ ] 14.6 Set up database query performance tracking
- [ ] 14.7 Implement distributed tracing (optional)
- [ ] 14.8 Configure alerting thresholds (response time > 500ms, error rate > 1%)
- [ ] 14.9 Create monitoring dashboard
- [ ] 14.10 Set up log aggregation and analysis

## 15. Deployment Preparation
- [x] 15.1 Create Dockerfile with multi-stage build
- [ ] 15.2 Set up CI/CD pipeline (GitHub Actions)
- [ ] 15.3 Configure production environment variables
- [x] 15.4 Create deployment documentation
- [x] 15.5 Set up database migration scripts
- [ ] 15.6 Configure production database (connection pooling, backups)
- [ ] 15.7 Set up staging environment for testing
- [ ] 15.8 Implement blue-green deployment strategy
- [ ] 15.9 Create rollback procedures
- [ ] 15.10 Verify all performance targets are met before production

## 16. Performance Verification
- [ ] 16.1 Run load tests with 1000+ concurrent users
- [ ] 16.2 Verify API response times (< 200ms p95 for all endpoints)
- [ ] 16.3 Verify database query times (< 50ms p95 for complex queries)
- [ ] 16.4 Test throughput capacity (500+ req/sec)
- [ ] 16.5 Verify cache hit rates (> 80% for cached endpoints)
- [ ] 16.6 Test error handling under load
- [ ] 16.7 Verify graceful degradation under high load
- [ ] 16.8 Test database connection pool behavior
- [ ] 16.9 Verify memory usage stays within limits (< 512MB per instance)
- [ ] 16.10 Document performance test results and benchmarks