# go-github-wip

[![Docker pulls](https://img.shields.io/docker/pulls/twinproduction/go-github-wip.svg)](https://cloud.docker.com/repository/docker/twinproduction/go-github-wip)

**go-github-wip** is an application that creates a GitHub check run on pull requests that 
have a title starting with `[WIP]`. This is used in order to prevent collaborators from
accidentally merging a PR that isn't completed yet.


## Table of Contents

- [go-github-wip](#go-github-wip)
  * [Table of Contents](#table-of-contents)
  * [Usage](#usage)
  * [Github App Requirements](#github-app-requirements)
    - [Permissions](#permissions)
    - [Events](#events)
  * [Environment variables](#environment-variables)


## Usage

Pull the image from Docker Hub:

```
docker pull twinproduction/go-github-wip:latest
```

Run it:

```
docker run --name go-github-wip -p 0.0.0.0:8080:8080 twinproduction/go-github-wip
```


## Github App Requirements

### Permissions

| Permission    | Access       | Use                                                         | 
|---------------|--------------|-------------------------------------------------------------|
| Checks        | Read & Write | Create a check suite to prevent users from merging          |
| Issues        | Read         | Check whether an issue already has a label associated to it |
| Pull requests | Read & Write | Create and delete a label from a PR                         |


### Events

- Pull request


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

If you wish, you can also configure the prefix that will trigger setting a PR as `work in progress` by setting `GO_GITHUB_WIP_PREFIXES`. The values are comma separated, meaning that `WIP!,[WIP]` would set both `WIP!` and `[WIP]` as prefixes. If no prefixes are defined, it will default to `WIP` and `[WIP]`.

```
GO_GITHUB_WIP_PREFIXES="WIP,[WIP]"
```

Optionally, you can also enable debugging by setting the `GO_GITHUB_WIP_DEBUG` to `true`:

```
GO_GITHUB_WIP_DEBUG="true"
```


### Building locally

```
docker build . -t go-github-wip
```
