# Gym Membership System - Sequence Diagrams

---

## 1. User Registration Flow (`POST /register`)

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant Client App
    participant API Gateway/Server
    participant Database

    User->>Client App: Fill registration form (email, password)
    Client App->>API Gateway/Server: POST /register { email, password }
    API Gateway/Server->>Database: Query user by email
    Database-->>API Gateway/Server: User not found

    API Gateway/Server->>API Gateway/Server: Hash password (Bcrypt)
    API Gateway/Server->>Database: INSERT INTO users (email, password_hash, deposit_amount = 0)
    Database-->>API Gateway/Server: User created successfully

    API Gateway/Server-->>Client App: 201 Created { message: "Registration successful" }
    Client App-->>User: Show success message & redirect to Login
```

---

## 2. User Login Flow (`POST /login`)

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant Client App
    participant API Gateway/Server
    participant Database

    User->>Client App: Enter credentials (email, password)
    Client App->>API Gateway/Server: POST /login { email, password }
    API Gateway/Server->>Database: SELECT * FROM users WHERE email = ?
    Database-->>API Gateway/Server: User record

    alt Invalid Credentials / User Not Found
        API Gateway/Server-->>Client App: 401 Unauthorized { error: "Invalid email or password" }
        Client App-->>User: Show error notification
    else Valid Credentials
        API Gateway/Server->>API Gateway/Server: Verify password hash
        API Gateway/Server->>API Gateway/Server: Generate JWT Token (with user_id payload)
        API Gateway/Server-->>Client App: 200 OK { token: "<JWT_TOKEN>" }
        Client App-->>User: Save token & redirect to Dashboard
    end
```

---

## 3. Balance Top-Up Flow (`POST /top-up`)

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant Client App
    participant API Gateway/Server
    participant Auth Middleware
    participant Database

    User->>Client App: Input top-up amount
    Client App->>API Gateway/Server: POST /top-up Header: Bearer <JWT> { amount }
    API Gateway/Server->>Auth Middleware: Validate Token
    
    alt Token Invalid / Expired
        Auth Middleware-->>Client App: 401 Unauthorized { error: "Invalid Token" }
        Client App-->>User: Prompt to re-login
    else Token Valid
        Auth Middleware->>API Gateway/Server: Pass user_id
        API Gateway/Server->>Database: UPDATE users SET deposit_amount = deposit_amount + ? WHERE id = ?
        Database-->>API Gateway/Server: Updated record
        API Gateway/Server-->>Client App: 200 OK { current_balance: new_amount }
        Client App-->>User: Show updated balance
    end
```

---

## 4. Product Rental / Booking Flow (`POST /rent-products`)

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant Client App
    participant API Gateway/Server
    participant Auth Middleware
    participant Database

    User->>Client App: Select Gym Class / Facility & click Rent
    Client App->>API Gateway/Server: POST /rent-products Header: Bearer <JWT> { product_id, quantity }
    API Gateway/Server->>Auth Middleware: Validate Token
    Auth Middleware->>API Gateway/Server: Extract user_id

    API Gateway/Server->>Database: SELECT price, stock FROM products WHERE id = ?
    Database-->>API Gateway/Server: Product details
    API Gateway/Server->>Database: SELECT deposit_amount FROM users WHERE id = ?
    Database-->>API Gateway/Server: User balance

    alt Insufficient Balance or Stock
        API Gateway/Server-->>Client App: 400 Bad Request { error: "Insufficient balance or stock" }
        Client App-->>User: Display failure message
    else Balance & Stock Sufficient
        API Gateway/Server->>Database: BEGIN TRANSACTION
        API Gateway/Server->>Database: UPDATE users SET deposit_amount = deposit_amount - total_price WHERE id = ?
        API Gateway/Server->>Database: UPDATE products SET stock = stock - quantity WHERE id = ?
        API Gateway/Server->>Database: INSERT INTO bookings (user_id, product_id, amount, status)
        API Gateway/Server->>Database: COMMIT TRANSACTION
        Database-->>API Gateway/Server: Transaction completed
        
        API Gateway/Server-->>Client App: 201 Created { booking_id, status: "SUCCESS" }
        Client App-->>User: Display booking confirmation
    end
```

---

## 5. Third-Party Payment Gateway Integration (e.g., Midtrans / Xendit)

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant Client App
    participant API Gateway/Server
    participant Payment Gateway API
    participant Database

    User->>Client App: Choose Payment via Payment Gateway
    Client App->>API Gateway/Server: POST /checkout Header: Bearer <JWT> { product_id }
    API Gateway/Server->>Database: CREATE pending order record
    
    API Gateway/Server->>Payment Gateway API: Request Snap Transaction / Payment URL (API Key)
    Payment Gateway API-->>API Gateway/Server: Return Payment URL & Order Token
    API Gateway/Server-->>Client App: 200 OK { redirect_url: "https://payment..." }

    Client App->>User: Redirect to Payment Gateway Checkout Page
    User->>Payment Gateway API: Complete payment (Bank Transfer / E-Wallet)
    
    Payment Gateway API->>API Gateway/Server: POST /webhook/payment-notification (Callback)
    API Gateway/Server->>API Gateway/Server: Verify Webhook Signature Key
    
    alt Payment Successful
        API Gateway/Server->>Database: UPDATE order_status = 'PAID', create booking
        API Gateway/Server-->>Payment Gateway API: 200 OK
    else Payment Failed / Expired
        API Gateway/Server->>Database: UPDATE order_status = 'EXPIRED'
        API Gateway/Server-->>Payment Gateway API: 200 OK
    end

    User->>Client App: Return to Application
    Client App->>API Gateway/Server: GET /booking-report Header: Bearer <JWT>
    API Gateway/Server->>Database: Fetch updated booking status
    Database-->>API Gateway/Server: Booking details
    API Gateway/Server-->>Client App: Display current booking status
```

---

## 6. Booking Report Flow (`GET /booking-report`)

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant Client App
    participant API Gateway/Server
    participant Auth Middleware
    participant Database

    User->>Client App: Navigate to Booking History
    Client App->>API Gateway/Server: GET /booking-report Header: Bearer <JWT>
    API Gateway/Server->>Auth Middleware: Validate Token
    Auth Middleware->>API Gateway/Server: Pass user_id

    API Gateway/Server->>Database: SELECT b.*, p.name FROM bookings b JOIN products p ON b.product_id = p.id WHERE user_id = ?
    Database-->>API Gateway/Server: List of booking records

    API Gateway/Server-->>Client App: 200 OK { data: [ booking_list ] }
    Client App-->>User: Render history list / report
```