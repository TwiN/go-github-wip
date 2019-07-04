# go-github-wip

[![Docker pulls](https://img.shields.io/docker/pulls/twinproduction/go-github-wip.svg)](https://cloud.docker.com/repository/docker/twinproduction/go-github-wip)

waw
## Table of Contents

- [go-github-wip](#go-github-wip)
  * [Table of Contents](#table-of-contents)

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

You must also specify the ID of your GitHub application using `GO_GITHUB_WIP_APP_ID`:

```
GO_GITHUB_WIP_APP_ID="12345"
```

Optionally, you can also enable debugging by setting the `GO_GITHUB_WIP_DEBUG` to `true`:

```
GO_GITHUB_WIP_DEBUG="true"
```


### Building locally

```
docker build . -t go-github-wip
```
