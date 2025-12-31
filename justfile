set dotenv-load := true
DEPLOY_TARGET := env_var('DEPLOY_TARGET')

# Build for Linux x86_64 (OCI Intel/AMD Instance)
build:
    GOOS=linux GOARCH=amd64 go build -o oci-checker .

# Deploy binary and .env to server
deploy: build
    rsync -avz ./oci-checker ./.env {{DEPLOY_TARGET}}

# Run locally (for testing)
run:
    go run main.go

# Check logs on server (requires DEPLOY_TARGET to be user@host:path)
logs:
    ssh oracle "tail -f /home/ubuntu/projects/oci-checker/oci-checker.log"
