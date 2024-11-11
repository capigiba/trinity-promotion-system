# Trinity App Report

## Table of Contents

1. [Technical Decisions](#technical-decisions)
2. [Assumptions Made](#assumptions-made)
3. [Future Improvements](#future-improvements)
4. [Local Setup Guide](#local-setup-guide)

---

## Technical Decisions

### 1. Framework Selection: Gin

- **Reasoning:** As a developer with limited experience, I chose the [Gin](https://github.com/gin-gonic/gin) framework for building RESTful APIs due to its simplicity and clear syntax. Gin provides a robust set of features out-of-the-box, such as routing, middleware support, and JSON validation, which streamline the development process.
  
- **Flexibility:** Gin's modular architecture allows for easy integration with other tools and services (e.g., AWS). This flexibility ensures that we can adapt or switch frameworks in the future if project requirements evolve.

### 2. Localization Support

- **Implementation:** To facilitate multi-language support, I integrated a localization module using YAML files stored in the `locales` directory. This approach allows for easy addition of new languages by simply creating new YAML files.
  
- **Flexibility:** Developers can choose their preferred programming language for localization, making the system adaptable to various team preferences and requirements.

### 3. Database Choice: MongoDB

- **Reasoning:** MongoDB, a NoSQL database, was selected for its ability to handle high-speed read and write operations, which is essential for our voucher campaign functionalities that require quick data access and manipulation.
  
- **Considerations:** While MongoDB excels in performance and scalability, it prioritizes flexibility and speed over strict consistency. This trade-off is acceptable for our project, where rapid voucher generation and redemption are critical.

- **Note:** If strict transactional consistency becomes a requirement in the future, we might need to consider alternative databases or implement additional consistency mechanisms.

### 4. Dependency Injection

- **Current Approach:** For simplicity, dependencies are injected manually within the application.
  
- **Future Consideration:** To enhance modularity and testability, we can adopt dependency injection frameworks like [Wire](https://github.com/google/wire). Wire automates the generation of dependency injection code, reducing boilerplate and potential errors.

---

## Assumptions Made

1. **Unique Identifiers:**
   - Each entity (Campaign, Voucher, Purchase, Subscription, User) is uniquely identified by a string-based `_id`.

2. **Voucher Redemption:**
   - A voucher can only be redeemed once and is linked to a single purchase upon redemption.

3. **User Authentication:**
   - User authentication mechanisms (e.g., password hashing, session management) are handled securely, though not detailed in the current scope.

4. **Localization:**
   - The application initially supports English (`en`), with the infrastructure in place to add more languages as needed.

5. **Database Consistency:**
   - MongoDB's eventual consistency model is sufficient for the application's current requirements, especially for voucher campaigns where speed is prioritized.

6. **Environment Configuration:**
   - Configuration settings (e.g., database URI, server port) are managed through environment variables, although currently injected manually without a configuration management tool.

---

## Future Improvements

### 1. Dynamic Discount Values

- **Current State:** The default discount value is hardcoded at 30%.
  
- **Enhancement:** Introduce a dynamic parameter for discounts, allowing administrators to set different discount percentages based on campaign requirements or user tiers.

### 2. Voucher Code Generation Output

- **Current State:** Voucher codes are printed to the terminal for simplicity.
  
- **Enhancement:** Implement functionality to export voucher codes in various formats such as `.txt` or `.xlsx`. This will facilitate bulk distribution and management of vouchers.

### 3. Expanded Localization

- **Current State:** Localization is implemented in the handler layer with support for English only. Logging and other layers lack localization support.
  
- **Enhancement:** Extend localization to all layers of the application, including services and repositories. Add support for multiple languages by creating additional YAML files in the `locales` directory.

### 4. Database Flexibility

- **Current State:** The application exclusively uses MongoDB for data persistence.
  
- **Enhancement:** Introduce support for alternative databases, such as Redis for in-memory data storage, to handle scenarios requiring faster access times or caching mechanisms.

### 5. Configuration Management with Viper

- **Current State:** Configuration settings are manually managed within the `config/config.go` file.
  
- **Enhancement:** Integrate [Viper](https://github.com/spf13/viper) for robust configuration management. Viper allows for reading configuration from environment variables, JSON, TOML, YAML, and more, enhancing the application's flexibility and security.

  - **Reasoning:** Using Viper will streamline the configuration process, support multiple environments (development, staging, production), and improve the application's scalability by making it easier to manage complex configuration settings.

### 6. Dependency Injection with Wire

- **Current State:** Dependencies are injected manually, leading to increased boilerplate code and potential errors.
  
- **Enhancement:** Adopt [Wire](https://github.com/google/wire) for automated dependency injection, reducing boilerplate and enhancing code maintainability.

### 7. Implement Persistent Caching

- **Enhancement:** Explore the use of persistent caching solutions like Redis to improve performance for frequently accessed data, reducing load on MongoDB and enhancing user experience.

### 8. Enhanced Error Handling and Logging

- **Enhancement:** Implement more granular error responses and extend logging across all application layers (services, repositories) to facilitate better monitoring and debugging.

### 9. API Rate Limiting and Security

- **Enhancement:** Introduce rate limiting to protect against abuse and implement comprehensive security measures (e.g., JWT authentication, input validation) to safeguard the application.

---

## Local Setup Guide

### Prerequisites

- **Go:** Ensure that Go is installed on your machine. You can download it from the [official website](https://golang.org/dl/).
- **Docker:** Required for running MongoDB via Docker. Install Docker from the [official website](https://www.docker.com/get-started).

### Setup

- move to this file: [setup.md](setup.md)