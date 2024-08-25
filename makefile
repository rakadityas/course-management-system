# do go build and run the binary
run:
	go build -o bin/course-management-system ./cmd && ./bin/course-management-system -configuration_path etc/development.json

# compose up
# Host: mysql
# Username: myuser
# Password: mypassword
# Database: course_management
compose-up:
	docker-compose up -d