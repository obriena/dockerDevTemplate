FROM golang

RUN apt-get update && apt-get install -y wget
RUN apt install inetutils-ping -y
RUN apt install curl -y

RUN go get gorm.io/driver/mysql
RUN go get gorm.io/gorm

RUN wget https://github.com/cespare/reflex/releases/download/v0.3.1/reflex_linux_amd64.tar.gz
RUN gunzip reflex_linux_amd64.tar.gz
RUN tar -xvf reflex_linux_amd64.tar
RUN mv ./reflex_linux_amd64/reflex /usr/local/bin/ 
RUN rm -rf /reflex_linux_amd64
RUN rm reflex_linux_amd64.tar


RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz
RUN gunzip migrate.linux-amd64.tar.gz 
RUN tar -xvf migrate.linux-amd64.tar
RUN chmod a+x ./migrate
RUN mv ./migrate /usr/local/bin/ 
RUN rm ./migrate*tar
RUN rm README.md

