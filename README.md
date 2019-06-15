# go-github-wip

[![Docker pulls](https://img.shields.io/docker/pulls/twinproduction/go-github-wip.svg)](https://cloud.docker.com/repository/docker/twinproduction/go-github-wip)

A very small Docker image that test


## Table of Contents

- [go-github-wip](#go-github-wip)
  * [Table of Contents](#table-of-contents)
  * [Usage](#usage)


## Usage

Pull the image from Docker Hub:

```
docker pull twinproduction/go-github-wip:latest
```

Run it:

```
docker run --name go-github-wip -p 0.0.0.0:80:80 twinproduction/go-github-wip
```

### Building locally

```
docker build . -t go-github-wip
```
