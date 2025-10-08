[//]: # (TOTOUCH: Add your own meaningful and well written README.)

# IONOS Cloud Go Sample Service

Sample Go service for IONOS Cloud utilizing a Hexagonal Architecture

## Getting Started

### Harbor

The build image it will be automatized and it will publish this image in a harbor registry.
You only need to configure the harbor and set the environment variables in the github vars/secrets.

#### Steps to create a harbor registry

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

### CI Workflows

The CI workflows are configured to build and push the image to the harbor registry on every push to the main branch.
You can find the workflow file in the `.github/workflows` folder.

#### Test Workflow

The test workflow is configured to run the tests on every push to the main branch and on every pull request.
You can find the workflow file in the `.github/workflows` folder.

We assume that you will use some libraris from private repositories, so you need to set the following environment variables in the github secrets:

You will need to setup a ionos-cloud-repo-read App.
Please follow the steps in the documentation: [Creating a new repo-read App](https://confluence.united-internet.org/spaces/ICDEV/pages/285381664/Resolving+internal+dependencies+with+ionos-cloud-repo-read+App)

After you have created the App, you will get the following credentials:
- IONOS_CLOUD_REPO_READ_APP_ID: The App ID of the ionos-cloud-repo-read App
- IONOS_CLOUD_REPO_READ_SECRET: The App Secret of the ionos-cloud-repo-read App

You may need to set the sonar token as well if you want to use sonarcloud for code quality checks:
- SONAR_TOKEN: The token for sonarcloud

#### Running the Service on a kubernetes cluster

First you need to create a kubernetes cluster, you can use the Managed Kubernetes Service from IONOS Cloud.
You can find the documentation here: [Managed Kubernetes Service](https://docs.ionos.com/cloud/containers/managed-kubernetes)
You will need a cluster for each environment you want to deploy the service to.

Alternatively, you can use any other kubernetes cluster if you have one available.

#### Steps to create bootstrap cluster dev/test/prod environments

First you need to possess a cluster as we described above.
Then you need to contact [Platform Framework](https://chat.google.com/room/AAAAwmSJfyo?cls=7) Team to bootstrap the cluster for you.
You can start from here: [Documentation](https://confluence.united-internet.org/spaces/ICDEV/pages/267005963/Developer+Platform+Service+Catalog)

#### Steps to get the host for your service

In case your service needs to be exposed via an ingress controller, you need to get a host for your service.
You can start from here: [OneDNS Documentation](https://confluence.united-internet.org/spaces/ICDEV/pages/267006202/Howto+create+a+new+location+zone)

#### Steps to deploy the service to the kubernetes cluster

To deploy the service you will need a repository where you will store the helm charts and the kubernetes manifests.
An example of such a repository is: [platform-s3-deployment](https://github.com/ionos-cloud/platform-s3-deployment)

To automatize the deployment, you can create a workflow similar to the one in the
[Example Workflow](https://github.com/ionos-cloud/platform-s3-deployment/actions/workflows/bump-image-tag.yml)
This needs to be dispatched automatically [Here](https://github.com/ionos-cloud/go-sample-service/actions/runs/18314015515/workflow)
You need to uncomment the last step in the workflow file and adjust the parameters accordingly.
