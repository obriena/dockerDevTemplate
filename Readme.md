# Docker Development Template

## Goals
This is a project to show how to setup a working development environment using docker containers

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
### Creating migrate scripts
Here is a helpful command
```
migrate create -ext sql -dir migrations -seq -digits 3 <function_detail_table>
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
docker build -f dev.build.DockerFile  -t flyingspheres/develpment:0.0.1 .
```

This only builds an image.  I took this image `flyingspheres/develpment:0.0.1` and replaced the `from golang` to 
`from flyingspheres/development:0.0.1`.  If I didn't do that, then everytime I launched the dev enviornment it would download migrate and 
reflex again.  Not a huge deal but this will be faster, and having more control over the tools necesary will be helpful things like ping and curl are now available to me as well.  

## Application Flow
Retrieve all posts
```mermaid
sequenceDiagram
    cmd/main.go-->>httpController: handle
    httpController-->>postInteractor: ReadAll(context)
    postInteractor-->>repo: ReadAll(Context)
    repo->>db: Find(allPosts)
    db-->>repo: data
    repo->>repo: toDomainList
    repo->>postInteractor: response
    postInteractor->>postInteractor: convertToPostOutput
    postInteractor->>httpController: RespondJson
```
```mermaid
---
title: Class Diagram
---
classDiagram
    domain_PostRepo <|-- repository_mysqlrepo_rep
    domain_PostInteractor <|-- post_interactor
    main_service 
    domain_PostRepo: Create(context, Post)
    domain_PostRepo: ReadAll(context)
    domain_PostRepo: ReadById(context, int)
    domain_PostRepo: Update(context, Post)

    domain_PostInteractor: Add(context, PostInteractorAddIniput) (output, error)
    domain_PostInteractor: ReadAll(context) (output, error)
    domain_PostInteractor: ReadById(context, ReadyByIdInput) (output, error)
    domain_PostInteractor: ReadByCreatedById(context, ReadByCreatedIdInput) (output, error)
    
    main_service : Service
    main_service: main()
    main_service: newService()
    main_service: (s *Service) Init() 
    main_service: pushProcessingDuration()
    main_service: pushProcessingCount()
    main_service: RespondJSON()
    main_service: connectDB()

    repository_mysqlrepo_rep: NewRepository(db *gorm.DB) 
    repository_mysqlrepo_rep: toDomainList([]postModel)
    repository_mysqlrepo_rep: (p *postModel)toDomain(postModel)
    post_interactor: NewInteractor(repo domain.PostRepo)
```

## Testing
There are two types of testing we are concerned about.  Unit testing and Integration testing.  Unit testing shoudl be code that tests the basic functionality of the code that is written.  Integration testing is testing who this system works in conjunction with the servers that support it - things like MySql, or other services.  
Another issue I ran into while exectuing tests is tha Go caches tests.  So if you are testing application code and it is changing your test may not be exectued with the regular command of 
`go test`
Passing in the commands will solve that problem
` go test -count=1`
To run a specific test you can do something like this
`go gest -count ./test/delete_test.go`

### Unit Testing
Go Lang has top tier support for testing, there are a few caveats that need to be understood before writing tests.  The test compiler/runner does it out of convention not configuration.  What this means is that if you have a test case that is not named correctly i.e <name>_test.go then the tests will not be run.  Also for test to run the test case methods must support this signagure:
```
func TestUpdate(t *testing.T) {
}
```
Typically in GoLang there is not a tests folder.  Tests live in the same directory as the file they are written to test.  I have created a test folder to specifically execute integration tests.  

### Inegration Testing
These tests typically invoke the service via an http call.  Here is an example of the delete test case:
```
package tests

import (
	"log"
	"net/http"
	"testing"
)

func TestDelete(t *testing.T) {
	// URL of the resource to delete
	url := "http://localhost/posts/1"

	// Create a new DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("Error creating DELETE request: %v", err)
	}

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing DELETE request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	log.Println("Resource deleted successfully!")
}
```
The whole point ofintegration tests is to create an easy mechanism to execute tests against a running system.  Where unit tests are tests that run against just the code.  There should be no server running when unit tests are run.