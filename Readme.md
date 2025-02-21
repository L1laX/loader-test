File Structure

project-root/
│
├── api/                     # Backend Service (Golang + Gin)
│   ├── main.go              # Main application file
│   ├── go.mod               # Go module file
│   └── go.sum               # Go dependencies lock file
│
├── nginx/                   # NGINX Configuration
│   └── nginx.conf           # NGINX configuration for load balancing
│
├── test/                    # Load Testing Scripts
│   └── k6.js                # k6 script for simulating high loads
│
├── docker-compose.yml       # Docker Compose file to orchestrate services
├── Dockerfile               # Dockerfile for building the Golang app
├── Makefile                 # Makefile for common commands
└── README.md                # Documentation
1. Start the Application

Use the provided Makefile to manage the application lifecycle.

Start the Application with Default Configuration
bash
Copy
1
make start
This command starts the application in detached mode (-d) with one instance of the backend service.

Start the Application with 5 Backend Instances
bash
Copy
1
make start-5
This command scales the backend service to 5 instances using --scale=5.

Rebuild and Restart the Application
If you make changes to the code or configuration, rebuild and restart the application:

bash
Copy
1
make rebuild
Stop the Application
To stop all running containers:
make down

profile
Qwen2.5-Max
11:08 pm
Here’s a well-structured README.md file that explains the purpose of your project, how to use the Makefile, and provides instructions for running and testing the application. This will make it easier for others (or your future self) to understand and work with the project.

High-Performance API with Golang, Gin, Docker, and Load Balancing
This project demonstrates how to build a high-performance API using Golang , Gin Framework , Docker , and NGINX as a load balancer. The system is designed to handle up to 100,000 requests per second (RPS) by leveraging horizontal scaling and caching with Redis.

Features
High Performance : Optimized for handling high traffic using Golang and Gin.
Load Balancing : NGINX distributes traffic across multiple backend instances.
Caching : Redis is used to cache user validation and balance checks, reducing external API calls.
Scalability : Easily scale the backend service using Docker Compose.
Hot Reload : Integrated with air for live-reloading during development.
Load Testing : Includes a k6 script for simulating high loads (e.g., 100,000 RPM).
Prerequisites
Before running the project, ensure you have the following installed:

Docker : Install Docker
Docker Compose : Included with Docker Desktop or install separately.
k6 : For load testing. Install from k6.io
File Structure
Copy
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
project-root/
│
├── api/                     # Backend Service (Golang + Gin)
│   ├── main.go              # Main application file
│   ├── go.mod               # Go module file
│   └── go.sum               # Go dependencies lock file
│
├── nginx/                   # NGINX Configuration
│   └── nginx.conf           # NGINX configuration for load balancing
│
├── test/                    # Load Testing Scripts
│   └── k6.js                # k6 script for simulating high loads
│
├── docker-compose.yml       # Docker Compose file to orchestrate services
├── Dockerfile               # Dockerfile for building the Golang app
├── Makefile                 # Makefile for common commands
└── README.md                # Documentation
Usage
1. Start the Application
Use the provided Makefile to manage the application lifecycle.

Start the Application with Default Configuration
bash
Copy
1
make start
This command starts the application in detached mode (-d) with one instance of the backend service.

Start the Application with 5 Backend Instances
bash
Copy
1
make start-5
This command scales the backend service to 5 instances using --scale=5.

Rebuild and Restart the Application
If you make changes to the code or configuration, rebuild and restart the application:

bash
Copy
1
make rebuild
Stop the Application
To stop all running containers:

bash
Copy
1
make down
2. Test the Application
Manual Testing
You can manually test the API using curl:

bash
Copy
1
curl http://localhost/check/12345
Load Testing with k6
Run the k6 script to simulate high loads:

bash
Copy
1
make test
The k6.js script is configured to simulate 100,000 requests per minute (RPM) . You can modify the script in ./test/k6.js to adjust the load.

3. Monitor Logs
To monitor logs from the backend service:

bash
Copy
1
docker-compose logs app
To monitor NGINX logs:

bash
Copy
1
cat ./nginx/logs/access.log