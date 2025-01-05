# 🐎 Horse Tracking Management System

## Overview
A comprehensive Go-based application for horse breeding, health tracking, and management.

## Features
- 🚀 High-performance caching mechanism
- 🔒 Robust input sanitization
- 📊 Detailed horse health and pregnancy tracking
- 🛡️ Secure database interactions

## Tech Stack
- Language: Go (Golang)
- Database: PostgreSQL
- Caching: In-memory cache
- Web Framework: Gin
- ORM: GORM
- Validation: Custom sanitization

## Setup

### Prerequisites
- Go 1.20+
- PostgreSQL
- Git

### Installation
1. Clone the repository
2. Set up database
3. Configure environment variables
4. Run `go mod tidy`
5. Start the application

## Testing
- Run tests: `go test ./...`
- Test database setup: Use `scripts/setup_test_db.sql`

## Performance
- Implemented generic caching interface
- Thread-safe in-memory cache
- Configurable cache duration

## Security
- Input sanitization
- SQL injection prevention
- Strict input validation

## Roadmap
See [TODO.md](TODO.md) for upcoming features

## License
© 2025 hulta equestrian. All Rights Reserved.

PROPRIETARY AND CONFIDENTIAL

THIS SOFTWARE IS THE EXCLUSIVE PROPERTY OF hulta equestrian. 

UNAUTHORIZED COPYING, REPRODUCTION, DISTRIBUTION, OR USE OF THIS SOFTWARE, IN WHOLE OR IN PART, BY ANY MEANS, IS STRICTLY PROHIBITED. 

THIS SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED. 

VIOLATION OF THESE TERMS WILL RESULT IN LEGAL ACTION.

NO LICENSE IS GRANTED TO USE, COPY, MODIFY, OR DISTRIBUTE THIS SOFTWARE UNDER ANY CIRCUMSTANCES.

## Contributing
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
