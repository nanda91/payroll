# PAYROLL

_**Empowering Payroll**_

[![Last Commit](https://img.shields.io/github/last-commit/yourusername/payroll?style=flat-square)](https://github.com/yourusername/payroll)
![Go](https://img.shields.io/badge/go-100.0%25-brightgreen.svg?style=flat-square)
![Languages](https://img.shields.io/github/languages/count/yourusername/payroll?style=flat-square)

_Built with the tools and technologies:_

![Go](https://img.shields.io/badge/-Go-00ADD8?logo=go&logoColor=white&style=flat-square)
![Gin](https://img.shields.io/badge/-Gin-black?logo=go&style=flat-square)
![GORM](https://img.shields.io/badge/-GORM-gray?style=flat-square)
![Docker](https://img.shields.io/badge/-Docker-2496ED?logo=docker&logoColor=white&style=flat-square)
![YAML](https://img.shields.io/badge/-YAML-red?style=flat-square)

---

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Testing](#testing)

---

### Overview

A payroll management system built with Go (Golang) and Gin framework.
- ✅ **Modular Repository Layer**: Facilitates clean separation of concerns for managing users, attendance, overtime, reimbursements, and payroll data.

- 🔐 **Secure Authentication**: Implements JWT-based role management and middleware for role-specific access control.

- 🧾 **Audit Trail Support**: Tracks user actions and data modifications for accountability and compliance.

- 🌐 **RESTful API**: Organized routing and handlers enable seamless integration with client applications.

- 🛠️ **Database Migration & Seeding**: Simplifies setup with automated schema management and initial data population.

- 🔁 **Extensible & Maintainable**: Designed for scalability, supporting future feature additions with ease.


### Project Structure

```text
/payroll
├── configs/         # Configuration files
│   ├── app.go       # Application settings
│   └── database.go  # DB configuration
├── database/
│   ├── migrations/  # Database migration files
│   └── seeds/       # Test data seeders
├── delivery/
│   ├── http/        # HTTP handlers
├── domain/          # Core business models
├── repositories/    # Data access layer
├── routes/          # API endpoint definitions
├── usecase/         # Business logic
├── utils/           # Helper functions
├── main.go          # Application entry point
└── tests/           # Test suites
    └── integration/ # Integration tests
```


---

## Getting Started

### Prerequisites

This project requires the following dependencies:

- **Programming Language**: Go
- **Package Manager**: Go modules

---

### Installation

Build payroll from the source and install dependencies:

1. **Clone the repository**:

   ```bash
   git clone https://github.com/nanda91/payroll.git

2. **Navigate to the project directory**:

   ```bash
   cd payroll

3. **Install the dependencies**:

   ```bash
   go mod tidy

4. **Run Migration**:

   ```bash
   go run main.go migrate

5. **Run Seeder**:

   ```bash
   go run main.go seed
   

### Usage

run the project:

    go run {entrypoint}

### Testing

payroll uses test framework. Run suite with

    go test ./tests/integration/... -v