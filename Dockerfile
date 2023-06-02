FROM  golang:1.16.7-alpine3.14 as builder

#ENV GOPATH /usr/local/go
#ENV PATH $GOPATH/bin:$PATH

# Set up build directories
RUN mkdir -p /app && \
    mkdir -p /BUILD && \
    mkdir -p /BUILD/db

# Build the 411 binary
COPY cmd /BUILD/cmd
COPY go.sum  /BUILD/go.sum
COPY go.mod /BUILD/go.mod
COPY internal /BUILD/internal
RUN cd /BUILD && go mod vendor && go mod download
RUN cd /BUILD && go build -o /BUILD/warmachine cmd/warmachine/main.go 



# Production container
FROM alpine

# Add user and set up temporary account
RUN mkdir /app && \
    mkdir app/temp && \
    addgroup web && \
    adduser --home /app --system --no-create-home web web && \
    chown -R web:web /app && \
    chmod 1777 app/temp 

#Copy Stuff
COPY --from=builder /BUILD/warmachine /app/warmachine
COPY sites /app/sites

WORKDIR /app

USER 411
EXPOSE 5000
    
CMD ["./warmachine"]
