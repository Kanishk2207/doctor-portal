
# Doctor Portal

A microservice-based doctor portal where receptionists can create, read, update, and delete (CRUD) patient records. Doctors can view and update patients' information. The application is written in raw Go, using primarily standard libraries. It is containerized with Docker and deployable via Kubernetes (K8s) scripts.

## Features
- CRUD operations for patient records (handled by receptionists).
- Doctors can view and update patient information.
- Written in Go using standard libraries.
- Containerized using Docker.
- Deployable via Kubernetes (K8s) scripts.
- Bash scripts for easy building and launching.
- Includes unit tests.

## Microservices
- **Auth Service**: Handles user authentication.
- **Patient Service**: Manages patient-related operations.

## Prerequisites
- Docker
- Linux terminal

## Setup Instructions

1. Clone the repository:

    ```bash
    git clone <your-repo-url>
    ```

2. **Auth Service**:

    - Navigate to the `auth_service` directory:
    
      ```bash
      cd auth_service
      ```

    - Build the service:
    
      ```bash
      ./build.sh Dockerfile
      ```

    - Launch the service container:
    
      ```bash
      ./launch.sh
      ```

    - To launch the service

        ```bash
            go run main.go 
        ```

        or
        
        ```bash
            ./launch_reload_server.sh
        ```
    - The Auth Service will be available on port `8081`.

3. **Patient Service**:

    - Open a new terminal and navigate to the `patient_service` directory:
    
      ```bash
      cd patient_service
      ```

    - Build the service:
    
      ```bash
      ./build.sh Dockerfile
      ```

    - Launch the service container:
    
      ```bash
      ./launch.sh
      ```

    - To launch the service

        ```bash
            go run main.go 
        ```

        or

        ```bash
            ./launch_reload_server.sh
        ```

    - The Patient Service will be available on port `8082`.

## Docs
For the swagger docs, You have to

  - the auth-service server

  - The docs are available at http://localhost:8081/docs/

## Deployment with Kubernetes
Kubernetes scripts are provided to deploy the application in a microservice architecture.

## Testing
Unit tests are included. To run the tests:

```bash
go test ./...
```

## Author
Kanishk Shrivastava
