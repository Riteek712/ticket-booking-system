# Event Ticketing System

## Overview
This project is a **ticket-event management system** built using **Go (Golang)** with **GORM** for database interactions, **PostgreSQL** as the database, **Redis** for caching / rate limiting, and **RabbitMQ** for message queuing.

### Features
- **User Registration & Authentication**
- **Event Creation & Management**
- **Ticket Booking System**
- **Database Schema Management (Migrations)**
- **Caching with Redis**
- **Message Queueing with RabbitMQ**
- **Dockerized Setup with `docker-compose`**

---

## Queue-Based Ticket Processing with RabbitMQ
This system leverages **RabbitMQ** as a **message queue** to handle ticket booking operations asynchronously.  

### **Why Use a Queue-Based System?**
- Prevents blocking the main API request/response cycle.
- Ensures **high availability** by handling **a large number of concurrent ticket bookings**.
- Provides **fault tolerance**—failed jobs can be **retried** without data loss.

### **How it Works**
1. When a user **books a ticket**, the request is **pushed into a queue** in RabbitMQ.
2. A **worker service** consumes the message **asynchronously** and processes the booking.
3. The worker **validates ticket availability** and stores booking details in the database.
4. A **confirmation email** is sent to the user upon successful booking.

---


---

## Tech Stack
| Technology  | Description |
|-------------|------------|
| Go (Golang) | Backend API Development |
| GORM        | ORM for Go (Database Handling) |
| PostgreSQL  | Relational Database Management System |
| Redis       | Caching Layer |
| RabbitMQ    | Message Broker |
| Docker      | Containerization |
| Swagger     | API Documentation |

---

## Setup Instructions

### Prerequisites
Ensure you have the following installed:
- **Go** (1.19 or later)
- **Docker** & **Docker Compose**
- **PostgreSQL** (if running locally)

###  Clone the Repository
```sh
git clone https://github.com/Riteek712/ticket-booking-system.git
cd event-ticketing-system
```
