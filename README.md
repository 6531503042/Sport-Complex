# Sport-Complex
This project is part of 1305308	Platform Development Course.
 ![alt text](assets/admin_facility.png)
## Overview

This project aims to redesign the existing sport complex system to improve its scalability, performance, and user experience. The primary goals are to handle a larger number of users, enable easy scaling using a Kubernetes cluster, and leverage modern technologies like Golang, Echo, gRPC, and Kafka for efficient communication and interaction.

## 📂 Project Structure
The frontend is structured into key pages, each serving a specific role in the system:

![alt text](assets/admin_homepage.png)
**Admin Homepage**: The main entry point for user.

![alt text](assets/admin_booking.png)
**Booking**: A page for client booking.

![alt text](assets/admin_payment.png)
**User Payment History**: A page for user to view payment history.

![alt text](assets/user_payment.png)
**Payment**: A page for client payment.



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

### Frontend
- **User Interface**
  - Responsive design for all devices
  - Dark/Light mode support
  - Accessible UI components
  - Real-time updates
  - Interactive booking calendar

### Backend Services

#### Authentication Service
- JWT-based authentication
- Role-based access control (RBAC)
- Session management
- Password encryption
- OAuth integration

#### Booking Service
- Real-time slot availability
- Concurrent booking handling
- Booking validation
- Cancellation management
- Slot locking mechanism

#### Facility Service
- Facility status tracking
- Maintenance scheduling
- Dynamic pricing
- Capacity management
- Resource allocation

#### Payment Service
- Secure payment processing
- Transaction history
- Refund handling
- Receipt generation
- Payment verification

#### Admin Service
- User management dashboard
- Booking oversight
- Facility management
- Analytics and reporting
- System configuration


### Figma Design
- [UI/UX Design](https://www.figma.com/file/xxxxx)
- [Component Library](https://www.figma.com/file/xxxxx)
- [Design System](https://www.figma.com/file/xxxxx)


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
            <li>Project Architecture & Foundation</li>
            <li>User Authentication Service</li>
            <li>Admin Dashboard Implementation
                <ul>
                    <li>User Management System</li>
                    <li>Facility Management Interface</li>
                    <li>Booking Enhancement</li>
                </ul>
            </li>
            <li>UI/UX Enhancement with Next.UI</li>
            <li>System Integration</li>
        </ul>
    </td>
    <td align="center">
        <img src="https://img.shields.io/badge/Commits-450-blue?style=for-the-badge"/>
        <img src="https://img.shields.io/badge/PRs-15-green?style=for-the-badge"/>
        <br/>
        <b>45% of total contributions</b>
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
            <li>Booking System Implementation
                <ul>
                    <li>Badminton Court Booking</li>
                    <li>Football Field Reservation</li>
                    <li>Fitness Room Scheduling</li>
                    <li>Swimming Pool Booking</li>
                </ul>
            </li>
            <li>Slot Management System</li>
            <li>Real-time Slot Refresh</li>
        </ul>
    </td>
    <td align="center">
        <img src="https://img.shields.io/badge/Commits-280-blue?style=for-the-badge"/>
        <img src="https://img.shields.io/badge/PRs-8-green?style=for-the-badge"/>
        <br/>
        <b>28% of total contributions</b>
    </td>
</tr>

<tr>
    <td>
        <a href="https://github.com/Kritsasoft">6531503005</a>
    </td>
    <td>Kritsakorn Sukkasem</td>
    <td>
        <b>Role:</b> Frontend Developer<br/>
        <b>Responsibilities:</b>
        <ul>
            <li>Authentication System
                <ul>
                    <li>Login Interface</li>
                    <li>Registration System</li>
                    <li>Role-Based Access Control</li>
                </ul>
            </li>
            <li>User Role Management</li>
            <li>Authentication Endpoints</li>
        </ul>
    </td>
    <td align="center">
        <img src="https://img.shields.io/badge/Commits-150-blue?style=for-the-badge"/>
        <img src="https://img.shields.io/badge/PRs-5-green?style=for-the-badge"/>
        <br/>
        <b>15% of total contributions</b>
    </td>
</tr>

<tr>
    <td>
        <a href="https://github.com/kongphop1209">6531503008</a>
    </td>
    <td>Kongphop Saenphai</td>
    <td>
        <b>Role:</b> Frontend Developer<br/>
        <b>Responsibilities:</b>
        <ul>
            <li>Homepage Development</li>
            <li>Route Management</li>
            <li>Admin Sidebar Implementation</li>
            <li>Page Transitions & Effects</li>
        </ul>
    </td>
    <td align="center">
        <img src="https://img.shields.io/badge/Commits-80-blue?style=for-the-badge"/>
        <img src="https://img.shields.io/badge/PRs-4-green?style=for-the-badge"/>
        <br/>
        <b>8% of total contributions</b>
    </td>
</tr>

<tr>
    <td>
        <a href="https://github.com/MABiuS1">6531503006</a>
    </td>
    <td>Klavivach Prajong</td>
    <td>
        <b>Role:</b> Full Stack Developer<br/>
        <b>Responsibilities:</b>
        <ul>
            <li>Payment System Integration</li>
            <li>Payment Page Development</li>
            <li>Payment Route Management</li>
        </ul>
    </td>
    <td align="center">
        <img src="https://img.shields.io/badge/Commits-50-blue?style=for-the-badge"/>
        <img src="https://img.shields.io/badge/PRs-3-green?style=for-the-badge"/>
        <br/>
        <b>4% of total contributions</b>
    </td>
</tr>

</table>

### Contribution Distribution
```mermaid
%%{init: {'theme': 'base', 'themeVariables': { 'pie1': '#4CAF50', 'pie2': '#2196F3', 'pie3': '#FF9800', 'pie4': '#9C27B0', 'pie5': '#F44336'}}}%%
pie
    title Total Commits: 1162
    "Nimit (Lead Dev & Full Stack)" : 450
    "Phumiphat (Booking System Frontend)" : 280
    "Kritsakorn (Auth System Frontend)" : 150
    "Kongphop (UI/Routes Frontend)" : 80
    "Klavivach (Payment & Full Stack)" : 50
```

### Weekly Contribution Timeline
```mermaid
gantt
    title Project Development Timeline
    dateFormat  YYYY-MM-DD
    section Project Setup
    Initial Planning        :2023-09-01, 7d
    Architecture Design     :2023-09-07, 14d
    Docker & K8s Setup     :2023-09-14, 10d
    
    section Frontend Development
    UI/UX Design           :2023-09-10, 14d
    Login/Register         :2023-09-15, 21d
    Homepage & Navigation  :2023-09-20, 14d
    Booking System        :2023-09-25, 30d
    Admin Dashboard       :2023-10-05, 25d
    Payment Interface     :2023-10-15, 14d
    
    section Backend Development
    Auth Service          :2023-09-15, 21d
    User Service          :2023-09-20, 21d
    Facility Service      :2023-09-25, 21d
    Booking Service       :2023-10-01, 25d
    Admin Service         :2023-10-10, 20d
    Payment Service       :2023-10-15, 21d
    
    section Integration & Testing
    API Integration       :2023-10-20, 14d
    gRPC Implementation   :2023-10-25, 10d
    Kafka Setup          :2023-10-30, 7d
    System Testing       :2023-11-01, 21d
    Bug Fixes           :2023-11-15, 14d
```

### Repository Statistics
- **Total Commits:** 1,010
- **Active Pull Requests:** 35
- **Completed Features:** 12
- **Active Contributors:** 5
- **Lines of Code:** 15,000+

View detailed contribution statistics: [Contributors Graph](https://github.com/6531503042/Sport-Complex/graphs/contributors)

Last updated: 2025-03-27 00:37:12 UTC
