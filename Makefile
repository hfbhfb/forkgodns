
all:
	export CGO_ENABLED=0 && go build
	docker build -t godnsimg:v0.0.1 .
	-docker stop godnslocal
	-docker rm godnslocal
	docker run --name godnslocal --restart always  -p 23153:53 -p 23153:53/udp -v /etc/hosts:/etc/hosts -d godnsimg:v0.0.1
	@echo "编译dockerfile"

