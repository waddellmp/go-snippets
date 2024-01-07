IMAGE_NAME = mysql:latest
CONTAINER_NAME = snippetbox

build:
	docker build -t $(IMAGE_NAME) .
run:
	docker run -d --name $(CONTAINER_NAME) -p 3306:3306 $(IMAGE_NAME) 