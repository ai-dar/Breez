# Breez Project

## Overview
Breez is a Twitter-like web application where users can register, log in, post tweets, like tweets, and admins can manage the platform, including sending bulk emails to users. This project demonstrates core web development features, including:

- **User Authentication**
- **Tweet Management**
- **Filtering, Sorting, and Pagination**
- **Error Handling**
- **Rate Limiting**
- **Email Notifications**
- **Administrative Panel**

## Features

### User Features
- **Registration and Login**: Secure user authentication with password hashing.
- **Post Tweets**: Users can post tweets.
- **Like Tweets**: Users can like or unlike tweets.
- **View Tweets**: Tweets are displayed with filtering, sorting, and pagination.

### Admin Features
- **Admin Registration**: Admins can register with elevated permissions.
- **Email Notifications**: Admins can send bulk emails to all registered users.
- **Admin Panel**: A dedicated admin interface for managing platform functionalities.

### Additional Features
- **Structured Logging**: Log all major actions and errors in a structured format for easier debugging.
- **Error Handling**: Handle errors gracefully with meaningful feedback to users.
- **Rate Limiting**: Protect the server from abuse with request rate limiting.
- **Graceful Shutdown**: Safely terminate the server while preserving ongoing operations.

## Setup

### Prerequisites
- Go (version 1.18 or higher)
- PostgreSQL
- [Gomail](https://pkg.go.dev/gopkg.in/gomail.v2) library for email

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/breez.git
   cd breez
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up the database:
   - Create a PostgreSQL database.
   - Update the `dsn` in `main.go` with your database credentials.
   ```go
   dsn := "user=postgres password=yourpassword dbname=breez sslmode=disable"
   ```
   - Run migrations automatically:
   ```bash
   go run main.go
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

5. Open the application in your browser:
   ```
   http://localhost:8080
   ```

## Project Structure

```plaintext
Breez/
├── handlers/          # Request handlers for routes
├── models/            # Database models
├── static/            # Static files (HTML, CSS, JS)
├── utils/             # Utility functions (email, etc.)
├── main.go            # Entry point of the application
├── README.md          # Project documentation
```

## API Endpoints

### User Endpoints
- `POST /register`: Register a new user.
- `POST /login`: Log in as a user.
- `GET /check-auth`: Check user authentication status.
- `POST /tweets`: Post a new tweet.
- `GET /tweets`: Fetch tweets with filtering, sorting, and pagination.
- `POST /like`: Like or unlike a tweet.

### Admin Endpoints
- `POST /admin/register`: Register a new admin.
- `POST /admin/send-emails`: Send bulk emails to users.

## Features in Detail

### Filtering, Sorting, and Pagination
- **Filtering**: Search tweets by content or user name.
- **Sorting**: Sort tweets by creation date or user.
- **Pagination**: View tweets in pages with a configurable limit.

### Rate Limiting
Limits users to 1 request per second with a burst capacity of 3 requests. Requests exceeding the limit return a `429 Too Many Requests` status.

### Email Notifications
Admins can send HTML emails to all registered users using the SMTP protocol. Configure the SMTP settings in `utils/email.go`.

### Structured Logging
Logs are output in JSON format for better readability and debugging:
```json
{
    "level": "info",
    "action": "fetch_tweets",
    "count": 5,
    "filter": "golang",
    "sort": "created_at",
    "page": 1
}
```

## Future Enhancements
- Add more roles and permissions.
- Enhance admin panel with additional management features.
- Implement a dashboard for analytics.

## Authors
- [Aidar]
- [Miras]

## License
This project is licensed under the MIT License.
