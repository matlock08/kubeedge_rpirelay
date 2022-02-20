FROM ubuntu:latest
RUN apt-get update 
RUN apt-get install wget -y
RUN wget --no-check-certificate https://project-downloads.drogon.net/wiringpi-latest.deb && dpkg -i wiringpi-latest.deb
RUN mkdir -p /rpirelay
COPY ./relay /rpirelay/
WORKDIR  rpirelay
