# LearnRoyale - Bacalaureat Learning Platform

## Table of Contents

1. [Product Vision](#product-vision)
2. [Product Features](#product-features)
3. [User Stories](#user-stories)
   - [Authentication and Profile](#authentication-and-profile)
   - [Courses and Lessons](#courses-and-lessons)
   - [Quizzes and Practice](#quizzes-and-practice)
   - [Multimedia Content](#multimedia-content)
   - [Gamification and Motivation](#gamification-and-motivation)
   - [AI Tutor and Support](#ai-tutor-and-support)
   - [Community and Communication](#community-and-communication)
   - [Reviews and Feedback](#reviews-and-feedback)
4. [User Scenarios (Usage Examples)](#user-scenarios-usage-examples)
5. [UML Diagram](#uml-diagram)
6. [Context Diagram](#context-diagram)
7. [Container Diagram](#container-diagram)
8. [Component Diagram (API Application)](#component-diagram-api-application)
9. [Architectural Description](#architectural-description)
   - [Product Summary](#product-summary)
   - [Comparison with Intermediate Deliverable](#comparison-with-intermediate-deliverable)
   - [Architectural Decisions](#architectural-decisions)
   - [C4 Diagrams Description](#c4-diagrams-description)
   - [Non-Functional Requirements and Architectural Solutions](#non-functional-requirements-and-architectural-solutions)
10. [QA (Quality Assurance)](#qa-quality-assurance)
    - [Building a Test Plan](#building-a-test-plan)
    - [Testing Objectives](#testing-objectives)
    - [Testing Process in SDLC](#testing-process-in-sdlc)
    - [Testing Methods](#testing-methods)
    - [Testing Results](#testing-results)
11. [Security Analysis](#security-analysis)
    - [Main Security Risks Analysis](#main-security-risks-analysis)
    - [Tactics for Addressing Risks](#tactics-for-addressing-risks)
12. [CI/CD](#cicd)
    - [Environment Description](#environment-description)
    - [Environment Differences](#environment-differences)
    - [Specific Configurations](#specific-configurations)
13. [Setup Local Development](#setup-local-development)
14. [Features](#features)

---

# Product Vision

**FOR** students and lifelong learners preparing for important exams or seeking to strengthen their knowledge through guided practice

**WHO** want an engaging, efficient, and personalized way to learn theory, practice problems, and track their progress

**THE** LearnRoyale App IS an intelligent learning ecosystem, available on mobile and web,

**THAT** combines structured courses, adaptive quizzes, personalized tutoring sessions, and AI powered feedback to help users build confidence, improve understanding, and achieve their academic goals.
**UNLIKE** traditional tutoring, static textbooks, or scattered online resources,

**OUR PRODUCT** offers a unified learning experience: interactive lessons aligned with national curricula, personalized study recommendations, real-time progress tracking, and motivational gamification, therefore making learning smarter, faster, and more enjoyable.

---

# Product Features

1. All BAC Subjects
LearnRoyale covers every subject included in the BAC exam. Students can easily select their field (Real or Uman) and dive into subject-specific content.

2. Past BAC Questions (Categorized by Topic)
The app includes real questions from previous BAC exams, organized by subject and topic for structured studying.

3. Mathematics Section (Difficulty Levels: M1, M2, M3)
Students can practice math questions based on their profile level:

M1: Advanced level (Science and Technical profiles)

M2: Intermediate level (Social Sciences)

M3: Basic level (Arts and Vocational profiles)

4. Romanian Language Section
Contains questions about literary works, grammar, synonyms, essay structures, and typical BAC topics.

5. Psihotrop Music Section
A special section featuring Psihotrop songs that creatively explain literary works and BAC materials through music.

6. Pauza de Mate Video Links
Curated YouTube videos from Pauza de Mate, helping students visualize and better understand math concepts.

7. Past BAC Tests
Full official BAC exams from previous years available for timed practice sessions.

8. General Knowledge Questions
Includes humanities-based general knowledge questions to help students broaden their understanding beyond textbooks.

9. Formula Practice Mode
A unique section where students must write formulas from memory, reinforcing key mathematical and scientific knowledge.

10. Timed Quizzes with Lives System
Each quiz is timed. If a student answers incorrectly or runs out of time, the app shows the correct answer and removes one life — similar to Duolingo's life system.

11. Learn Mode
A study-first mode that allows users to go through materials, notes, and explanations before attempting quizzes.

12. Friend Challenges & Competitions
Students can challenge their friends to 10-question matches from the same subject, competing for high scores and rankings.

13. User Forum
A community forum where users can ask questions, share study tips, discuss BAC topics, and motivate each other.

14. Learning Rewards
Users earn points, badges, and streaks for completing lessons, quizzes, and daily goals. These can unlock bonuses or cosmetic rewards.

15. Competitions & Leaderboards
National or school-wide competitions encourage students to test their knowledge and climb the leaderboard for prizes or recognition.

---

# User Stories

## Authentication and Profile
1. **As a** new user, **I want** to create an account with my email and password **so that** I can save my progress and access personalized content.

2. **As a** registered user, **I want** to log in to my account **so that** I can continue my learning from where I left off.

3. **As a** student, **I want** to view my profile with XP, streak, and completed lessons **so that** I can track my overall progress.

## Courses and Lessons
4. **As a** BAC student, **I want** to browse courses by subject (Romanian, Math, Informatics) **so that** I can find relevant study material.

5. **As a** Math student, **I want** to select my difficulty level (M1, M2, or M3) **so that** I can practice problems appropriate for my profile.

6. **As a** student, **I want** to access past BAC questions categorized by topic **so that** I can practice specific areas where I need improvement.

7. **As a** learner, **I want** to complete lessons in Learn Mode before taking quizzes **so that** I can understand the material first.

## Quizzes and Practice
8. **As a** student, **I want** to take timed quizzes with a lives system **so that** I can simulate exam conditions and learn from my mistakes.

9. **As a** student, **I want** to practice writing formulas from memory **so that** I can reinforce key mathematical concepts.

10. **As a** user, **I want** to see the correct answer when I make a mistake **so that** I can learn and avoid repeating errors.

11. **As a** BAC candidate, **I want** to take full past BAC exams in timed practice sessions **so that** I can prepare for the real exam format.

## Multimedia Content
12. **As a** Romanian Literature student, **I want** to listen to Psihotrop songs about literary works **so that** I can memorize key ideas more easily.

13. **As a** Math student, **I want** to watch Pauza de Mate videos **so that** I can better visualize and understand complex concepts.

## Gamification and Motivation
14. **As a** user, **I want** to earn points, badges, and maintain streaks **so that** I stay motivated to study daily.

15. **As a** competitive student, **I want** to view leaderboards **so that** I can compare my performance with other students.

16. **As a** student, **I want** to challenge my friends to 10-question matches **so that** we can learn together in a fun, competitive way.

17. **As a** participant, **I want** to join national or school-wide competitions **so that** I can test my knowledge and win recognition.

## AI Tutor and Support
18. **As a** student, **I want** to ask questions to an AI tutor **so that** I can get instant clarifications without waiting for a teacher.

19. **As a** learner, **I want** to receive personalized study recommendations **so that** I can focus on areas where I need the most improvement.

## Community and Communication
20. **As a** user, **I want** to participate in a forum **so that** I can ask questions, share tips, and connect with other students.

21. **As a** user, **I want** to subscribe to a newsletter **so that** I receive updates and study tips via email.

## Reviews and Feedback
22. **As a** student, **I want** to leave reviews for courses **so that** I can help other students choose quality content.

23. **As a** potential user, **I want** to read course reviews **so that** I can decide which courses are worth my time.

---

# User Scenarios (Usage Examples)

Scenario 1: Daily Practice
> User opens LearnRoyale before school and takes a short 5-minute Romanian quiz. He answers most questions correctly, earns a few points, and keeps his 10-day streak going.

Scenario 2: Quick Revision
> User uses LearnRoyale on the bus to review formulas for Math. User writes a few from memory in Formula Mode and sees instant feedback when User gets one wrong.

Scenario 3: Friend Challenge
> User invites his friend to a short Literature challenge in LearnRoyale. They both answer 10 questions, and User wins by just one point. They laugh about it and decide to rematch later.

Scenario 4: Learning Through Music
> User listens to Psihotrop songs in the Romanian section while reading the related literary notes. User remembers the main ideas more easily when studying later.

Scenario 5: Competition Day
> During the weekend, LearnRoyale launches a mini BAC competition. Dozens of students from the same school join in, and everyone tries to reach the leaderboard.

Scenario 6: AI Tutor Assistance
> User is stuck on a complex math problem. User opens the AI Chat and asks for help. The AI tutor explains the concept step-by-step, and User finally understands how to solve similar problems.

Scenario 7: Pre-Exam Full Test
> Two weeks before BAC, User takes a complete past exam in timed mode. User identifies weak areas and spends the remaining days focusing on those topics.

Scenario 8: Newsletter Discovery
> User subscribes to the LearnRoyale newsletter. Every week, User receives study tips and notifications about new courses, helping User stay on track.

---

# UML Diagram

## UML Diagram Explanation (Use Case)

The Use Case UML diagram presents the interactions between users (actors) and the main functionalities of the LearnRoyale system:

**Actors:**
- **Student** - the main user who accesses the platform for learning
- **System** - the automated component that manages internal processes

**Main Use Cases:**
1. **Authentication and Registration** - Users can create a new account or authenticate to access full functionalities
2. **Course Navigation** - Exploring available courses by categories (Romanian, Mathematics, Informatics)
3. **Quiz Participation** - Solving tests with questions from various subjects
4. **Progress Visualization** - Monitoring XP, streak, and completed lessons
5. **AI Tutor Interaction** - Conversations with the AI assistant for clarifications
6. **Newsletter Subscription** - Receiving updates via email
7. **Leaderboard Visualization** - Comparing performance with other users

![Diagram](./UML.png)

---

# Context Diagram

## Context Diagram Explanation (C4 Level 1)

The context diagram presents the LearnRoyale system from a high-level perspective (C4 Model - Level 1), showing how it interacts with users and external systems:

**Main Entities:**
- **Student** - The end user who accesses the platform through a browser (HTTPS)
- **LearnRoyale Account** - The user identity management system
- **LearnRoyale** - The main application (monolith)

**External Systems:**
- **PostgreSQL Database** - Relational database for storing users, progress, courses, and reviews (SQL communication)
- **Google Gemini API** - External AI service for the intelligent tutor (HTTPS/JSON communication)
- **SMTP E-mail Server** - Server for sending welcome emails and notifications (SMTP/TLS protocol)
- **YouTube** - Video platform for educational content (iframe embedding)

**Data Flows:**
- Student creates an account and authenticates in the platform
- Platform saves/reads data from PostgreSQL
- AI tutor communicates with Gemini for intelligent responses
- Newsletter sends emails through the SMTP server

```mermaid
flowchart TB
   student["Student"]
   account["LearnRoyale Account\n(User Account)"]
   learnroyale["LearnRoyale"]
   db[("PostgreSQL Database")]
   gemini["Google Gemini API"]
   smtp["SMTP E-mail Server"]
   yt["YouTube"]

   student -->|"Creates/uses account"| account
   account -->|"Authenticates in"| learnroyale
   student -->|"Uses platform (HTTPS)"| learnroyale
   learnroyale -->|"Reads/Writes (SQL)"| db
   learnroyale -->|"Prompts & responses (HTTPS/JSON)"| gemini
   learnroyale -->|"Welcome email (SMTP/TLS)"| smtp
   smtp -->|"Delivers emails to"| student
   learnroyale -->|"Embeds video (iframe)"| yt
```

---

# Container Diagram

## Container Diagram Explanation (C4 Level 2)

The container diagram details the major components of the LearnRoyale system and how they communicate with each other:

**Main Containers:**
- **Single-Page Application (SPA)** - The React frontend that runs in the user's browser, providing the interactive interface
- **API Application** - The Go backend (net/http) that handles business logic and communication with external services
- **PostgreSQL Database** - The persistent storage system for all application data

**Container Interactions:**
1. The SPA communicates with the API Application via HTTP/JSON (ports localhost:5173 → localhost:8000)
2. API Application accesses PostgreSQL through SQL queries using the pgx driver
3. API Application integrates Google Gemini for AI Tutor functionality
4. The SMTP server handles email sending (Newsletter)

![Container Diagram](./container.jpeg)

---

# Component Diagram (API Application)

## Component Diagram Explanation (C4 Level 3)

The component diagram details the internal structure of the API Application (Go Backend):

**Infrastructure Components:**
- **CORS Middleware** - Handles Cross-Origin policies to allow communication with the frontend
- **HTTP Router/Mux** - Distributes HTTP requests to the corresponding handlers

**Controllers (Handlers):**
- **Auth Controller** (`/api/register`, `/api/login`) - Handles user registration and authentication with JWT
- **Profile Controller** (`/api/profile`) - Returns the authenticated user's profile data
- **Progress Controller** (`/api/progress`) - Saves quiz progress and updates XP
- **Stats Controller** (`/api/stats`) - Provides detailed statistics (XP, streak, weekly activity)
- **Courses Controller** (`/api/courses`, `/api/course`) - Manages the list of courses by categories
- **Lessons Controller** (`/api/lessons`, `/api/lesson`, `/api/userlessons`) - Manages lessons and per-lesson progress
- **Reviews Controller** (`/api/reviews`) - Allows creating and viewing reviews
- **Newsletter Controller** (`/api/subscribe`) - Handles newsletter subscriptions with confirmation email
- **AI Chat Controller** (`/api/ai-chat`) - Integrates Google Gemini for AI tutoring

**Access Components:**
- **DB Access** - PostgreSQL access layer using pgx and parameterized SQL
- **Email Component** - Sends HTML emails via SMTP/TLS
- **AI Client** - HTTP client for communication with Google Gemini API

```mermaid
flowchart LR
   %% External containers/systems
   spa["Single-Page Application<br/>(React)<br/>Container"]
   db[("PostgreSQL Database<br/>(External System)")]
   gemini["Google Gemini API<br/>(External System)"]
   smtp["SMTP E-mail Server<br/>(External System)"]

   %% API container boundary
   subgraph API_APP["API Application (Go net/http)<br/>Container"]
      direction LR
      cors["CORS Middleware<br/>(Component)"]
      router["HTTP Router / Mux<br/>(Component)<br/>/api/* endpoints"]

      subgraph CTL["Controllers (Components)"]
         direction LR
         auth["Auth<br/>/api/register<br/>/api/login"]
         profile["Profile<br/>/api/profile"]
         progress["Progress<br/>/api/progress"]
         stats["Stats<br/>/api/stats"]
         courses["Courses<br/>/api/courses<br/>/api/course"]
         lessons["Lessons<br/>/api/lessons<br/>/api/lesson<br/>/api/userlessons"]
         reviews["Reviews<br/>/api/reviews"]
         newsletter["Newsletter<br/>/api/subscribe"]
         aichat["AI Chat<br/>/api/ai-chat"]
      end

      dbaccess["DB Access<br/>(Component)<br/>pgx + SQL"]
      email["Email<br/>(Component)"]
      aiclient["AI Client<br/>(Component)<br/>Gemini HTTP"]
   end

   %% Incoming calls (simplified to reduce clutter)
   spa -->|"HTTP /api/*"| cors
   cors --> router
   router --> auth
   router --> profile
   router --> progress
   router --> stats
   router --> courses
   router --> lessons
   router --> reviews
   router --> newsletter
   router --> aichat

   %% Core dependencies
   auth --> dbaccess
   profile --> dbaccess
   progress --> dbaccess
   stats --> dbaccess
   courses --> dbaccess
   lessons --> dbaccess
   reviews --> dbaccess
   newsletter --> dbaccess
   newsletter -->|"welcome email"| email
   aichat --> aiclient

   %% External calls
   dbaccess -->|"SQL"| db
   aiclient -->|"HTTPS/JSON"| gemini
   email -->|"SMTP/TLS"| smtp

   %% Styling
   classDef container fill:#2E7D32,stroke:#1B5E20,color:#fff;
   classDef external fill:#9E9E9E,stroke:#757575,color:#fff;

   class spa container;
   class db,gemini,smtp external;
   class cors,router,auth,profile,progress,stats,courses,lessons,reviews,newsletter,aichat,dbaccess,email,aiclient container;
   style API_APP fill:#ffffff,stroke:#2E7D32,stroke-width:2px,stroke-dasharray: 6 6;
   style CTL fill:#ffffff,stroke:#2E7D32,stroke-width:1px,stroke-dasharray: 4 4;
```

---

# Architectural Description

## Product Summary

LearnRoyale is a complete educational platform for Baccalaureate exam preparation, built using a modern client-server architecture. The application combines an interactive React frontend with a performant Go backend, offering advanced features such as AI tutoring, gamification, and progress monitoring.

### Comparison with Intermediate Deliverable

**Newly implemented features:**
- Complete authentication system with JWT
- Google Gemini AI integration for intelligent tutoring
- Newsletter with welcome emails (SMTP/TLS)
- Leaderboard with national ranking
- Reviews and feedback system
- XP tracking and daily streaks
- Quizzes for multiple subjects (Romanian, Mathematics M1-M4, Informatics)

### Architectural Decisions

1. **Monolith Architecture**: We chose a monolithic architecture for simplicity and development speed, with the possibility of migrating to microservices in the future.

2. **SPA + REST API**: The React frontend communicates with the Go backend via REST API, allowing independent development and scalability.

3. **PostgreSQL**: Relational database for data consistency and support for complex queries.

4. **JWT Authentication**: Stateless authentication for horizontal scalability.

5. **External AI Service**: Gemini API integration instead of a local model for superior response quality.

## C4 Diagrams Description

### System Diagram (Context)
Visualized above - shows the LearnRoyale system and interactions with external actors (Student, PostgreSQL, Gemini, SMTP, YouTube).

### Container Diagram
Visualized above - details the 3 main containers: React SPA, Go API Application, PostgreSQL Database.

### Component Diagram
Visualized above - shows the internal structure of the API Application with all controllers and access components.

## Non-Functional Requirements and Architectural Solutions

| Requirement | Implemented Solution |
|-------------|---------------------|
| **Performance** | Go for backend (compiled, concurrent), React with Vite for optimized builds |
| **Scalability** | Stateless architecture with JWT, PostgreSQL with optimized indexing |
| **Security** | bcrypt for passwords, JWT with 24h expiration, configured CORS, HTTPS for external APIs |
| **Availability** | Persistent PostgreSQL, graceful shutdown in Go |
| **Maintainability** | Modular code, clear frontend/backend separation, inline documentation |

---

# QA (Quality Assurance)

## Building a Test Plan

### Testing Objectives

| Level | Tested Artifacts | Objective |
|-------|------------------|----------|
| **Unit** | Individual Go functions (hashPassword, generateToken, validateToken) | Verify correctness of atomic functions |
| **Integration** | API Handlers (register, login, progress) | Verify complete request-response flows |
| **E2E** | Complete user flows | Validate real scenarios from user perspective |

### Testing Process in SDLC

Testing is integrated into the development cycle as follows:
1. **Development**: Unit tests written in parallel with code
2. **Code Review**: Automatic test execution before merge
3. **Staging**: Integration tests on staging environment
4. **Production**: Smoke tests after deployment

### Testing Methods

#### 1. Unit Tests (Go Testing Framework)

```go
// Example from main_test.go
func TestHashPasswordAndCheck(t *testing.T) {
    pass := "mysecret"
    hash, err := hashPassword(pass)
    if err != nil {
        t.Fatalf("hash error: %v", err)
    }
    if !checkPasswordHash(pass, hash) {
        t.Fatalf("password should match hash")
    }
}

func TestTokenGenerateValidate(t *testing.T) {
    token, err := generateToken(1, "testuser")
    if err != nil {
        t.Fatalf("token error: %v", err)
    }
    claims, err := validateToken(token)
    if err != nil {
        t.Fatalf("validate error: %v", err)
    }
    if claims.UserID != 1 || claims.Username != "testuser" {
        t.Fatalf("claims mismatch")
    }
}
```

**Relevance**: Unit tests validate critical security functions (password hashing, JWT generation/validation) that are fundamental for user data protection.

#### 2. Integration Tests (httptest)

```go
// Example from main_test.go
func TestApiRegisterHandler_Success(t *testing.T) {
    reqBody := `{"username":"test-user","email":"test-user@example.com","password":"123456"}`
    req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(reqBody))
    w := httptest.NewRecorder()

    apiRegisterHandler(w, req)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d", w.Code)
    }
}

func TestApiLoginHandler_InvalidPassword(t *testing.T) {
    reqBody := `{"email":"test-login2@example.com","password":"wrong"}`
    req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(reqBody))
    w := httptest.NewRecorder()

    apiLoginHandler(w, req)

    if w.Code != http.StatusUnauthorized {
        t.Fatalf("expected 401, got %d", w.Code)
    }
}
```

**Relevance**: Integration tests validate complete authentication flows, ensuring the API responds correctly to various scenarios (success, validation errors, invalid credentials).

#### 3. E2E Tests (Selenium)

The `selenium_test.side` file contains automated tests for:
- Main page navigation
- New user registration
- Authentication
- Course navigation
- Quiz completion

### Testing Results

| Category | Total Tests | Passed | Failed | Coverage |
|----------|-------------|--------|--------|----------|
| Unit Tests | 5 | 5 | 0 | ~85% |
| Integration Tests | 4 | 4 | 0 | API Handlers |
| E2E Tests | 8 | 8 | 0 | Critical Flows |

**Implementation Observations:**
- Authentication tests exposed the need for cleanup in `TestMain`
- Middleware tests validated correct propagation of X-User-ID headers
- Rate limiting is not yet implemented (identified as improvement)

---

# Security Analysis

## Main Security Risks Analysis

### 1. Injection Attacks (SQL Injection)

**Risk**: Attackers can inject malicious SQL code through unvalidated inputs.

**Implemented Mitigation**:
```go
// We use prepared parameters ($1, $2) instead of string concatenation
err = db.QueryRow(
    "SELECT id, username, password_hash FROM users WHERE email = $1",
    req.Email,
).Scan(&userID, &username, &passwordHash)
```

### 2. Broken Authentication

**Risk**: Weak passwords, predictable tokens, unsecured sessions.

**Implemented Mitigation**:
```go
// bcrypt with default cost factor (10) for password hashing
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// JWT with 24h expiration
claims := &Claims{
    UserID:   userID,
    Username: username,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    },
}
```

### 3. Cross-Site Scripting (XSS)

**Risk**: Injection of malicious scripts into application pages.

**Implemented Mitigation**:
- React automatic escaping of JSX content
- Content-Type: application/json for all API responses
- Server-side input validation

### 4. Sensitive Data Exposure

**Risk**: Exposing passwords or sensitive data in logs or responses.

**Implemented Mitigation**:
```go
// Password is never returned in responses
type AuthResponse struct {
    Token    string `json:"token"`
    Username string `json:"username"`
    UserID   int    `json:"userId"`
    // We don't include password_hash
}
```

### 5. Cross-Origin Resource Sharing (CORS)

**Risk**: Unauthorized access from external domains.

**Implemented Mitigation**:
```go
func withCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") // In production: specific domain
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        // ...
    })
}
```

## Tactics for Addressing Risks

| Risk | Tactic | Status |
|------|--------|--------|
| SQL Injection | Prepared statements | ✅ Implemented |
| Weak Passwords | Minimum 6 character length validation | ✅ Implemented |
| Token Theft | JWT with 24h expiration | ✅ Implemented |
| Brute Force | Rate limiting | ⚠️ TODO |
| Data Breach | bcrypt hashing | ✅ Implemented |
| XSS | React escaping + JSON API | ✅ Implemented |
| CSRF | JWT in header (not cookie) | ✅ Implemented |

---

# CI/CD

## Environment Description

### 1. Development (Local)

**Characteristics:**
- Frontend: `localhost:5173` (Vite dev server with HMR)
- Backend: `localhost:8000` (Go server)
- Database: Local PostgreSQL (`localhost:5432/local_db`)
- AI: Google Gemini API (with test key)

**Configuration:**
```bash
# .env (development)
DATABASE_URL=postgres://dev:1234@localhost:5432/local_db?sslmode=disable
GEMINI_API_KEY=your-dev-api-key
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SENDER_EMAIL=your-email@gmail.com
SENDER_PASSWORD=your-app-password
```

### 2. Staging

**Characteristics:**
- Production replica for testing
- Database separate from production
- Test API keys

**Differences from Development:**
- Optimized frontend builds
- SSL/TLS enabled
- More detailed logging

### 3. Production

**Characteristics:**
- Frontend served by Go server (`SERVE_FRONTEND=1`)
- Managed PostgreSQL (cloud)
- HTTPS mandatory

**Configuration:**
```bash
# .env (production)
DATABASE_URL=postgres://user:pass@cloud-db:5432/prod_db?sslmode=require
GEMINI_API_KEY=production-api-key
SERVE_FRONTEND=1
```

## Environment Differences

| Aspect | Development | Staging | Production |
|--------|-------------|---------|------------|
| **Frontend** | Vite dev server | Static build | Static build served by Go |
| **Database** | Local PostgreSQL | Local PostgreSQL | Cloud (prod) |
| **Debug** | Verbose logging | Moderate | Minimal |
| **CORS** | `*` (all origins) | Staging domain | Production domain |
| **API Keys** | Test keys | Test keys | Production keys |

## Specific Configurations

### Build Frontend for Production:
```bash
cd frontend
npm run build
# Output in frontend/dist/
```

### Run Backend with Frontend served:
```bash
SERVE_FRONTEND=1 go run main.go
# Accessible at http://localhost:8000
```

### Run Tests:
```bash
go test -v ./...
# Runs all tests from main_test.go
```

### Database Migrations:
```bash
# Schema is automatically applied at server startup
# through the initTables() function in main.go
```

---

# Setup Local Development

### 1. Clone Repository
```bash
git clone https://github.com/unibuc-ro/proiect-inginerie-software-learnroyale.git
cd proiect-inginerie-software-learnroyale
```

### 2. Setup Environment Variables
Copy `.env.example` to `.env` and fill in your values:
```bash
cp .env.example .env
```

Then edit `.env` with your credentials:
- **DATABASE_URL**: PostgreSQL connection string
- **GEMINI_API_KEY**: Get from [Google AI Studio](https://aistudio.google.com/apikey)
- **SENDER_EMAIL**: Gmail account for sending emails
- **SENDER_PASSWORD**: Gmail App Password (not regular password!)

**For Gmail App Password:**
1. Enable 2FA on your Google account
2. Go to [myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords)
3. Generate app password for "Mail" + "Windows Computer"
4. Copy the 16-character password (without spaces) to `.env`

### 3. Start Backend
```bash
go run main.go
```

### 4. Start Frontend (separate terminal)
```bash
cd frontend
npm install
npm run dev
```

---

# Features
- ✅ Newsletter subscription with email notifications
- ✅ User authentication (JWT)
- ✅ Course management
- ✅ XP & Streak tracking
- ✅ AI-powered tutoring (Gemini API)
