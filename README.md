# Go microservices

A repo designed to practise building microservices in Go, using an array of different protocols/technologies.
These microservices can be run individually in a Docker container, or be deployed as a Distributed App using either Docker Swarm or Kubernetes.

Available microservices:

- Broker
- Authentication
- Logger
- Mail
- Listener


A front-end is available to test each microservice through a simple web UI.

In a production environment, microservices should be able to directly communicate with each other (according to their needs) and the broker microservices won't expose public API endpoints to do so through actions.