# WarMachine - A Go Project

Built to be fast and scalable.  

# Project Layout

```
warmachine
  cmd/warmachine/main.go -> Main go program, entry point into application
  internal --> libraries internal to this code base
    config
      config.go --> Loads configuration from environment such as SITE=411.com
    handlers
      defaulthandler.go --> Single handler file, but can be split up into multiple for easier code readability
      handlercontext.go --> Context is use to pass statefull constructs such as database pools, or redis pools, stuff you want to keep around
    metrics 
      metrics.go --> loads some of the common metrics, but more metrics can be injected in other places
    middleware
      middlewhare.go --> Contains various bits of logic that can be chained to a request to add functionality, such as logging or authentication
    routes
      routes.go --> Contains all the routes of the application, this could be static or dynamic routes.  In addition this could be modified to load from a configuration file.
    util
      util.go --> Various utility functions
  sites
    411.com --> 411.com site files, more sites can be created and referenced from the "SITE" environmental variable
  Dockerfile
  go.mod --> contains versioned dependencies.  These are pinned until you want to adjust them
  go.sum --> contains the md5 sums for binaries
  vendor --> may exist, may not
  README.md --> This file you are reading now!  


```

# Setting Up Development Environment

- Install GoLang [Go](https://golang.org/doc/install)
- Optional Visual Studio Code, has nice intellisense, linter and formatting.  

# Running Project

In your terminal, in the root directory of the project:

```bash
SITE=411.com ENV=dev go run cmd/warmachine/main.go 
```

Open your browser and visit:

http://localhost:5001


