# Docker Development Template

## Goals
This is a project to show how to setup a workign development environment using docker containers

## Project Setup
There are two core configuration files from a development perspective.
1) `dev.Dockerfile`
1) `dev-docker-compose.yaml`

Through these files the development environment is configured as well as the other services required to build an application.

## Getting Started
Start up the supporting servers.  At the time of this writing those servers that are started up are:
1) MySQL
1) PHP MySql Admin
1) Prometheus
1) Prometheus Gateway

You should be able to start these components by executing the following command (this should be done in a terminal window):
`docker-compose -f ./dev-docker-compose.yaml`

Once those servers are started make sure that you have the [dev containers plugin](https://code.visualstudio.com/docs/devcontainers/containers) plugin installed you should see an icon like this: <img src="./images/devContainer.png" style="width:40px"/>  in the bottom left of the VSCode window.  If you click on that you will get an options menu that looks like this:  
<img src="./images/devContainerOptions.png" style="width:480px"/>  

You will want to select the option for "Reopen in container"
Follow the prompts that make the most sense for your use case and eventually select `Run`
If you are using this template there is a `.devcontainer/decontainer.json` file checked in,
so it may just close down your window of vscode and reopen with a new vscode window in the 
the docker container.  

Test the configuration by exectuing this command:  
```
go version
```
you should see something like this:
```
go version go1.22.2 linux/arm64
```

### Start the server
Assuming all is going well at this point, the development server is configured to run with a library called `reflex` this library monitors file changes so that when a change is made the code is rebuilt and the server is restarted automatically.  This usually happens in seconds and is typically unnoticeable.  You can use this by issuing the following command within the vscode terminal like this:  
```
reflex -c ./reflex.conf
```

## Architecture
This is strictly from a development point of view.
```mermaid
flowchart LR
    subgraph docker
    GoLangRestService
    end
    subgraph dockercompose
    MySql
    PHPMyAdmin
    Prometheus
    PrometheusPush
    end

    PHPMyAdmin --3306---> MySql
    Prometheus --9091 read---> PrometheusPush
    Prometheus --54719 read---> GoLangRestService

    GoLangRestService--6033---> MySql
    GoLangRestService--9091 write--->PrometheusPush

```

## Database Migrations
The Docker file we are developing has a library called [migrate](https://github.com/golang-migrate/migrate/) all structural changes to the database should be managed through this.  Here is a command to start:
```
migrate -path migrations -database "mysql://db_user:db_user_pass@tcp(host.docker.internal:6033)/app_db" up
```

## Database Integration
So once we have some data out there we are going to need to be able to write to it.  We will use the [GORM](https://gorm.io/index.html) library to manage this.  We are going to follow a domain driven technique of using data repositories, and interactors to achieve this. 

## Building the development image
Initially development went well enough with a docker image that looked like this:  
```
from golang
```
This worked well enough but it is better to be working with an image that has Migrage and Reflex pre-installed so I created another docker file to create an image that has the tools in it that I will need to develop in

to build the dev image here is the command:  
```
docker build -f dev.DockerFile  -t flyingspheres/develpment:0.0.1 .
```

This only builds an image.  I took this image `flyingspheres/develpment:0.0.1` and replaced the `from golang` to 
`from flyingspheres/development:0.0.1`.  If I didn't do that, then everytime I launched the dev enviornment it would download migrate and 
reflex again.  Not a huge deal but this will be faster, and having more control over the tools necesary will be helpful things like ping and curl are now available to me as well.  