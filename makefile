# configuration for to running on local
export DATABASE_URL=myuser:mypassword@tcp(localhost:3306)/course_management?parseTime=true
export APP_PORT=:8991

# do go build and run the binary
run:
	go build -o bin/course-management-system ./cmd && ./bin/course-management-system

# building the dockerfile
compose-build:
	docker-compose build

# compose up
compose-up:
	docker-compose up -d

# stop the Docker containers
compose-down:
	docker-compose down

# rebuild and restart the Docker containers
compose-restart: compose-down compose-up
