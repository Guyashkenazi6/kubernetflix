# 🎬 Kubernetflix

Kubernetflix is a movie management service, with CRUD operations for movies. The service is designed to run on Kubernetes and uses MongoDB for data persistence.

## 🌟 Features

- **Web Interface**: Welcome page with a Kubernetflix logo.
- **Movie CRUD Operations**: Manage movies with endpoints for:
  - Retrieving all movies
  - Fetching a movie by ID
  - Creating a new movie
  - Updating a movie
  - Deleting a movie

## 🛠️ Tech Stack

- **Web Framework**: Gin-Gonic
- **Database**: MongoDB
- **Containerization**: Docker
- **Orchestration**: Kubernetes
- **Continuous Deployment**: ArgoCD & GitHub Actions


## 🔗 Endpoints

- **Welcome Page**: `/`
- **All Movies**: `/gmovies`
- **Single Movie (by ID)**: `/gmovies/:movie_id`
- **Create Movie**: `/cmovies`
- **Update Movie**: `/umovies/:movie_id`
- **Delete Movie**: `/dmovies/:movie_id`

## 📦 Deployment

The project includes ArgoCD deployment configurations and a GitHub Workflow for CI/CD.

- **Kubernetes Deployment**: Refer to `deployment.yaml` for details. This deployment sets up both the main application and a MongoDB instance.
- **Kubernetes Service**: The service is exposed via LoadBalancer on port `8010`. Check `service.yaml` for more details.
- **GitHub Workflow**: There's an automated pipeline to build and push Docker images upon pushes to the `main` branch.
