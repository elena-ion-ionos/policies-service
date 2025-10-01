[//]: # (TOTOUCH: Add your own meaningful and well written README.)

# IONOS Cloud Go Sample Service

Sample Go service for IONOS Cloud utilizing a Hexagonal Architecture

## Getting Started

The build image it will be automatized and it will publish this image in a harbor registry.
You only need to configure the harbor and set the environment variables in the github vars/secrets.

Steps to create a harbor registry:
Default Harbor:
1. Go to: [Harbor: "https://harbor.tp.infra.cluster.ionos.com/harbor/projects"]
2. Create a new Project in our case is: go-sample-project
3. Create a robot account with the permission you need
4. use credentials in github secrets.

If you want to use your own harbor, you can do it as well, just set the environment variables in the github secrets.
The main environment variables you need to set are:
- HARBOR_URL: The URL of the harbor registry(it needs to contain also the project name)
- HARBOR_USERNAME: The username of the robot account
- HARBOR_PASSWORD: The password of the robot account
