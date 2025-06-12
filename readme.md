# Blyncq POC

This is a Go-based proof-of-concept project, Dockerized for easy deployment.

---

## 📦 Requirements

* Go (1.21+)
* Docker
* make

---

## 🏃‍♂️ Run the App Locally (Without Docker)

You can run the app directly using Go:

```bash
make run
```

---

## 🐳 Build and Push Docker Image

This project includes `Makefile` targets to easily build a Docker image for the application and push it to a Docker registry.

### Build the Docker Image

To build the Docker image, use the `build` target:

```bash
make build
```

This command will:
1.  Format the Go source code (`go fmt ./...`).
2.  Run Go vet to check for common errors (`go vet ./...`).
3.  Build the Docker image using the `build/dockerfile` and tag it as `hitesharma/blyncq-poc:latest` (or with a specified version if `VERSION` is set in the `Makefile`).

### Push to Docker Registry

After building the image, you can push it to your Docker registry (e.g., Docker Hub) using the `push` target:

```bash
make push
```

This command will push the previously built Docker image (`hitesharma/blyncq-poc:latest` or your specified version) to the registry. Make sure you are logged in to your Docker registry before attempting to push.

### Build and Push (Combined)

For convenience, you can combine the build and push operations into a single command:

```bash
make build-push
```

This will first build the Docker image and then push it to the registry.
