# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Andre Passos <andrebritopassos@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

ARG PROJECT

ENV PROJECT_ID=$PROJECT

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main"]