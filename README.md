# KubeVue

This is work in progress, check [TODO](TODO.md).

## Motivation

Provide UI for custom objects (argo workflows) with actions with authentication (OIDC) and authorization (group membership),
and expose services uniformly with authenticating reverse proxy.

## Use case

Provide per project (namespace) UI to run and manage argo workflows and expose pre-defined services with uniform access management.

## Architecture

KubeVue uses custom resource definitions to configure what objects and services to expose.

## Pre-requisites

OIDC server is required for the application to work. It could be either external OIDC provider (Okta, Auth0), or included
dex.

## Development

Make sure `dex.default` resolves to `127.0.0.1` on your development machine (e.g. put an entry to `/etc/hosts`), and then in skaffold folder:

```sh
skaffold dev --port-forward
```

After successful deployment point your browser to `http://localhost:8080/ui`.
