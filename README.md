# Flashcard Learning Backend API

## 📚 Project Overview
A REST API backend for a digital flashcard learning platform built with Go. Users can create personal collections of flashcards, manage their study materials, and practice with randomized card sequences for effective learning.

## 🎯 Core Functionality
- **User Authentication**: JWT-based registration and login system
- **Collection Management**: Create, read, update, delete flashcard collections
- **Card Management**: Add, edit, remove individual flashcards within collections
- **Practice Mode**: Random card shuffling for study sessions
- **Profile Management**: Update username and password

## 🏗️ Architecture
**Clean Architecture Pattern**: Handler → Service → Repository → Database
- **Handler Layer**: HTTP request/response handling with Chi router
- **Service Layer**: Business logic implementation
- **Repository Layer**: Data access abstraction with GORM
- **Database**: PostgreSQL with cascade delete relationships

## 📊 Data Model
```
Users (1:N) Collections (1:N) Cards
- id, username, password_hash    - id, name, description, user_id    - id, question, answer, collection_id
```

## 🔧 Technology Stack
- **Language**: Go 1.22.1
- **Framework**: Chi router, GORM ORM
- **Database**: PostgreSQL
- **Authentication**: JWT tokens with bcrypt password hashing
- **Validation**: go-playground/validator
- **Containerization**: Docker & Docker Compose
- **Testing**: testify framework

## 🚀 API Endpoints
- `POST /auth/signup` - User registration
- `POST /auth/login` - User authentication
- `GET /collections` - Get user collections
- `POST /collections` - Create new collection
- `POST /collections/{id}/cards` - Add card to collection
- `GET /collections/{id}/train` - Start practice session (randomized cards)
- `PUT/DELETE /cards/{id}` - Update/delete specific cards

## 🎲 Key Features
- **Random Practice**: Cards are shuffled randomly for each practice session
- **Secure Authentication**: JWT tokens with configurable expiration
- **Data Validation**: Request validation at multiple layers
- **Error Handling**: Comprehensive error responses with proper HTTP status codes
- **Cascade Operations**: Deleting collections automatically removes associated cards

## 💡 Use Cases
Perfect for students, language learners, or anyone who wants to create digital flashcards for memorization and spaced repetition learning.