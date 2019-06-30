# go-github-wip

[![Docker pulls](https://img.shields.io/docker/pulls/twinproduction/go-github-wip.svg)](https://cloud.docker.com/repository/docker/twinproduction/go-github-wip)

A very small Docker image that...


## Table of Contents

- [go-github-wip](#go-github-wip)
  * [Table of Contents](#table-of-contents)
  * [Usage](#usage)
  * [Environment variables](#environment-variables)


## Usage

Pull the image from Docker Hub:

```
docker pull twinproduction/go-github-wip:latest
```

Run it:

```
docker run --name go-github-wip -p 0.0.0.0:80:80 twinproduction/go-github-wip
```


## Environment variables

You must set `GO_GITHUB_WIP_APP_PRIVATE_KEY` to the name of the file containing your Github App private key.

e.g.

```
GO_GITHUB_WIP_APP_PRIVATE_KEY="github-app-private-key.pem"
```


### Building locally

```
docker build . -t go-github-wip
```