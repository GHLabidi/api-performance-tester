# Build the Go application
FROM golang:1.20.5-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the container's workspace
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the local package files to the container's workspace
COPY . .


# Build the Go application inside the container
RUN go build -o app cmd/app/main.go



# Use the Python image
FROM python:3.9-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Python scripts to the container's workspace
COPY python/ python/

# Install Python dependencies from requirements.txt
RUN pip install --no-cache-dir -r python/requirements.txt

# Copy the Go binary from the build stage to the Python image
COPY --from=build /app/app ./

# Copy the templates to the container's workspace
COPY templates/ templates/

# Copy the .env file to the container's workspace
COPY .env ./

# Expose the port that your HTTP server will run on
EXPOSE 8081

# Command to run your application
CMD ["./app"]