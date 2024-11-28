# Sport-Complex
This project is part of 1305308	Platform Development Course.
 ![alt text](assets/image.png)
## Overview

This project aims to redesign the existing sport complex system to improve its scalability, performance, and user experience. The primary goals are to handle a larger number of users, enable easy scaling using a Kubernetes cluster, and leverage modern technologies like Golang, Echo, gRPC, and Kafka for efficient communication and interaction.

## Acknowledgements

- [Golang] (https://golang.org/)
- [Echo] (https://echo.labstack.com/)
- [gRPC] (https://grpc.io/)
- [Kafka] (https://kafka.apache.org/)
- [OOP] (https://www.w3schools.com/java/java_oop.asp)
- [Mono-Microservice Architecture] (https://microservices.io/patterns/microservices.html)
- [Kubernetes] (https://kubernetes.io/)
- [Docker] (https://www.docker.com/)
- [MongoDB] (https://www.mongodb.com/)
- [Git] (https://git-scm.com/)
- [GitHub] (https://github.com/)
- [JWT] (https://jwt.io/)
- [Next.js] (https://nextjs.org/)
- [Next UI] (https://nextui.org/)
- [Tailwind CSS] (https://tailwindcss.com/)
- [React] (https://reactjs.org/)

## Features

- **User Management**: Register and manage users, including guests and different roles (e.g., campus users, outsiders, guest).
- **Activity Booking**: Book activities such as gym sessions, swimming, and badminton courts.
- **Dynamic Pricing**: Support for different pricing based on user roles (insiders, outsiders) and activities.
- **Scalability**: Designed to scale with Kubernetes for handling large numbers of concurrent users.
- **Real-time Communication**: Utilizes gRPC for fast and efficient communication between services.
- **Asynchronous Processing**: Employs Kafka for managing and processing asynchronous tasks.

## Architecture
![alt text](<assets/Screenshot 2567-09-01 at 21.57.59.png>)
The architecture should include the following components:

- **User Service**: Handles user registration, login, profile management, and other user-related functions.
- **Booking Service**: Manages booking of facilities, timeslots, and payments.
- **Facility Service**: Manages information about available facilities and their pricing.
- **Notification Servic**: Sends notifications to users about bookings, cancellations, and other updates.
- **Payment Service**: Handles payment processing and integration with payment gateways.
- **Middleware**: Acts as a gateway for incoming requests, routing them to the appropriate services.
- **Kafka**: Used for asynchronous communication between services, especially for events like booking confirmations and payment updates.

## Setup

### Prerequisites

- Golang 1.23.0
- Docker
- Kubernetes (Minikube or a K8s cluster)
- Apache Kafka
- MongoDB

### Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/6531503042/Sport-Complex.git
   cd sport complex

<h2>📃 Start App in Terminal</h2>
Start

## Contributors

### Team Members and Contributions

<table>
<tr>
    <th>Student ID</th>
    <th>Name</th>
    <th>Contributions</th>
    <th>Statistics</th>
</tr>

<tr>
    <td>
        <a href="https://github.com/6531503042">6531503042</a>
    </td>
    <td>Nimit Tanboontor</td>
    <td>
        <b>Role:</b> Lead Developer<br/>
        <b>Responsibilities:</b>
        <ul>
            <li>Project Foundation & Architecture</li>
            <li>User Authentication Service</li>
            <li>Facility Management System</li>
            <li>Booking Service Implementation</li>
        </ul>
    </td>
    <td>
        <img src="https://img.shields.io/badge/Commits-15-blue"/>
        <img src="https://img.shields.io/badge/PRs-3-green"/>
        <br/>
        <img src="https://github-readme-stats.vercel.app/api/pin/?username=6531503042&repo=Sport-Complex&show_owner=true" width="200"/>
    </td>
</tr>

<tr>
    <td>
        <a href="https://github.com/MABiuS1">6531503006</a>
    </td>
    <td>Klavivach Prajong</td>
    <td>
        <b>Role:</b> Backend Developer<br/>
        <b>Responsibilities:</b>
        <ul>
            <li>Payment Service Integration</li>
            <li>Booking-Facility Integration</li>
            <li>System Integration Testing</li>
        </ul>
    </td>
    <td>
        <img src="https://img.shields.io/badge/Commits-8-blue"/>
        <img src="https://img.shields.io/badge/PRs-2-green"/>
    </td>
</tr>

<tr>
    <td>
        <a href="https://github.com/2547phumiphat">6531503117</a>
    </td>
    <td>Phumiphat Wongsathit</td>
    <td>
        <b>Role:</b> Frontend Developer<br/>
        <b>Responsibilities:</b>
        <ul>
            <li>User Interface Development</li>
            <li>Frontend Testing</li>
        </ul>
    </td>
    <td>
        <img src="https://img.shields.io/badge/Commits-5-blue"/>
        <img src="https://img.shields.io/badge/PRs-1-green"/>
    </td>
</tr>

</table>

### Project Contribution Overview
<div align="center">
    <img src="https://github-profile-summary-cards.vercel.app/api/cards/profile-details?username=6531503042&theme=github" width="600"/>
</div>

### Repository Activity
<div align="center">
    <img src="https://github-readme-activity-graph.vercel.app/graph?username=6531503042&repo=Sport-Complex&theme=github-light" width="600"/>
</div>

### Commit Distribution
```mermaid
pie
    title Commit Distribution
    "Nimit (Lead)" : 15
    "Klavivach (Backend)" : 8
    "Phumiphat (Frontend)" : 5
```

### Repository Statistics
- **Total Commits:** 28
- **Active PRs:** 6
- **Completed Features:** 4
- **Active Contributors:** 3

View detailed contribution statistics: [Contributors Graph](https://github.com/6531503042/Sport-Complex/graphs/contributors)

Last updated: Manual update (Will be automated by GitHub Actions)