
all:
	go build
	docker build -t godnsimg:v0.0.1 .
	@echo "编译dockerfile"