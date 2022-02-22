# syntax=docker/dockerfile:1
FROM ubuntu:latest
RUN apt-get update 
RUN apt-get install wget -y
RUN mkdir -p /client
COPY ./client /client/
WORKDIR  client
EXPOSE 1323
