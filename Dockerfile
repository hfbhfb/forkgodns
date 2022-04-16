FROM alpine:3.14
COPY . /program/
WORKDIR /program
CMD ["sh", "-c", "nohup /program/godns"]
