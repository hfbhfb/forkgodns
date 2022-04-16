#FROM alpine:3.14
FROM ubuntu:18.04
COPY . /program/
WORKDIR /program
CMD ["sh", "-c", "nohup /program/godns"]
