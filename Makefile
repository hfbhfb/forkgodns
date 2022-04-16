
all:
	export CGO_ENABLED=0 && go build
	docker build -t godnsimg:v0.0.1 .
	-docker stop godnslocal
	-docker rm godnslocal
	#docker run  --name godnslocal53  --net host --restart always -v /etc/hosts:/etc/hosts -d godnsimg:v0.0.1
	docker run  --name godnslocal53  --net host --env DNSHOST=1.1.1.1 --restart always -v /etc/hosts:/etc/hosts -d godnsimg:v0.0.1
	@echo "编译dockerfile"

