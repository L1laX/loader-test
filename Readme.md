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
: make start
This command starts the application in detached mode (-d) with one instance of the backend service.

Start the Application with 5 Backend Instances
: make start-5
This command scales the backend service to 5 instances using --scale=5.

Rebuild and Restart the Application
If you make changes to the code or configuration, rebuild and restart the application:
: make rebuild

Stop the Application
To stop all running containers:
make down


2. Test the Application
Manual Testing
You can manually test the API using curl:



Run the k6 script to simulate high loads:
: make test
The k6.js script is configured to simulate 100,000 requests per minute (RPM) . You can modify the script in ./test/k6.js to adjust the load.
