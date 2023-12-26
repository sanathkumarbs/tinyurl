# TinyURL Service

## Context and Scope
The project aims to provide an implementation for TinyURL service. It will focus on primarily being deployable locally, with possibly extending it to be deployed on AWS.

## Goals
A TinyURL service has two key use-cases (listed below) which are table stakes for this project.
* Given a website link, provide a shorterned version of it
* Given a shortnered version, redirect the link to the original link

Additionally, this service can be extended for the following use-cases which will be explored in the future.
* Generate a TinyURL based on user-defined shorterned name
* Provide user-management capabilities through user accounts
* Provide a user to look-up the website destination for a given shorterned name
* Provide a user to set time-to-live (expiration) for a TinyURL
* Provide analytics on TinyURLs such as visits

## Requirements

### Functional Requirements
* Generating the TinyURL
    - Should not be longer than 10 characters
    - Should have a default expiration of 5 years and cannot be modified
* Redirection from TinyURL to original website
    - The user should be automatically redirected to the original link when they visit the TinyURL
    - Handle link expirations in a graceful manner

### Non-Functional Requirements
* User Scale
    - Supports a user-base of creation of 10,000 TinyURLs per day
    - Supports an access-pattern of 20% URLs being heavily used at 10,000 per day, where as, the remaining 80% is used at a rate of 1,000 per day
    - Has a design that can evolve without large refactors as the service gets traction
* Rate Limit
	- Apply sane rate-limiting to ensure the services will not be overwhelmed and provides safeguaring against malicious users
* Latency
    - The service needs to have low latency, especially when a user is visiting a TinyURL to reach the original destination
* High Reliability and Availability
    - The service needs to be be highly realiable and available as we will a high number of users trying to use TinyURL to reach their original destination

## Design

### Constraints
* The service will be designed for a 2x scale than anticipated
    - 20,000 TinyURL Creations
    - 10,000/1,000 access per day for 20/80 percent of TinyURLS available in the system
* For the scope of the project (hobby), we will ensure the choices made for storage and compute are reasonable
    - A more detailed analysis for understanding trade-offs for various options available will be made in-future.

### System Diagram
TBD: Actual design diagram
User -> API Gateway (NGINX) -> TinyURL Services -> Datastore (Postgres)



## Design Principles

### Project Design
* The project will be built using Go, and, uses standard Go Project layout, idomatic Go and general best practices around it
* The project will be designed such that it can be deployed in a Docker Runtime on Kubernetes
* The project will use existing open-source projects and toolchain where possible and will not reinvent the wheel

### Service Design
* The services will be designed such that they can be scaled horizontally to match the traffic patterns
* All services will follow the standard RESTful Design Practices
* All the services are stateless in nature to help with scaling and being fault-tolerant

### Authentication and Authorization
The services will not use any authentication or authorization as we do not have any user accounts, all the APIs will be accessible publicly

### Auditing
The service logs, along with the API Gateway logs will provide us the necessary Auditing logs. 

### Observability
All the services and its dependencies will natively integrate and provide observatbility data through OpenTelemetry.


## APIs

### Services
The TinyURL service will expose an API called `/generate` whose responsibility is to create a TinyURL for a given original website link.

TBD: Use OpenAPI Schema definition if possible?

**Generate a TinyURL**
* Path: `/generate`
* Methods: POST
* Encoding: JSON

Request Data:

	request: 
		object: generateURLRequest

	generateURLRequest:
		required:
			url: string
		optional:
			expiry: date 

	response: 
		object: generateURLResponse

		generateURLResponse:
			required:
				url: string
				short: string
				expiry: date 

### Persistence
The TinyURL data will be persisted in the `postgres` datastore for long-term storage. A table schema design will be added here for reference. 

We will use various caching strategies to provide low latency environment, and, details around the same will be provided here as well. 

## Estimations

### Resource
This section will evaulate the resources required to meet the desired scale as defined the requirements section.

#### Storage
TBD 

#### Compute
TBD 

### Costs
As the project will be eventually deployed in an AWS environment, we need to identify the costs associated with deploying and running this project for the scale defined above. The information here will guide us towards making cost optimzations where needed.

#### Storage
TBD 

#### Compute
TBD 
