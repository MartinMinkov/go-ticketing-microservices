# Microservice-based Event Ticketing System

## Description

This project demonstrates a reliable, full-stack event ticketing solution developed as part of a self-guided learning project. Inspired by the functionality of Ticketmaster, this scaled-down version effectively emulates key features like user registration, authentication, ticket purchasing, and order management, all through an intuitive user interface.

The application's backend is masterfully crafted in Go using the Gin framework, renowned for its performance and productivity, forming a set of microservices. Each service is responsible for a specific functionality, operating independently, yet smoothly communicating via a shared message queue implemented using NATS Streaming.

For data persistence, the application leverages MongoDB, a highly scalable and performance-oriented NoSQL database, ensuring efficient data management and retrieval.

On the front-end, the application uses TypeScript along with Next.JS, a popular and modern framework, to provide a responsive and engaging user experience. The JWT (JSON Web Tokens) protocol is used to manage user authentication in a secure and efficient manner.

Adopting containerization practices, the application is designed with Docker and Kubernetes, ensuring easy deployment, scalability, and system resilience. This comprehensive project serves as an exemplary demonstration of the intersection of various technologies and architectural patterns, reflecting strong technical competence and a dedication to continuous learning.

## Technologies Used

- Backend development with Go and Gin
- Frontend development with TypeScript and Next.JS
- Microservice architecture with NATS Streaming
- Containerization and orchestration with Docker and Kubernetes
- User Authentication with JWT (JSON Web Tokens)
- Database management with MongoDB
