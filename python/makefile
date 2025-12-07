# configuration for running on local
export DATABASE_URL=myuser:mypassword@tcp(localhost:3306)/course_management?parseTime=true
export APP_PORT=:8992

# run the Python app
run:
	python -m python.cmd.main

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
