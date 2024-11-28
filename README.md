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

### 6531503042
![GitHub Stats](https://github-readme-stats.vercel.app/api?username=6531503042&show_icons=true&count_private=true&theme=radical)