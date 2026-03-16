# kts-cli

Personal command-line toolkit used to simplify recurring daily tasks.

This repository contains a collection of small CLI commands I built over time to reduce friction in common workflows involving Kubernetes, cloud infrastructure, containers, Git, and general system operations. The commands are opinionated and optimized for my own environment and habits.

The goal of this project is not to be a general-purpose tool, but rather a personal automation layer on top of tools I already use.

## Purpose

Instead of remembering long commands, scripts, or pipelines, I centralize frequently used workflows into a single CLI entry point.

Typical use cases include:

* inspecting Kubernetes resources
* interacting with cloud providers
* managing Docker images
* automating Git operations
* manipulating YAML and files
* inspecting infrastructure state
* running helper utilities for development

Most commands provide small interactive flows that make it faster to select resources or execute tasks.

## Structure

The project follows a Go CLI layout.

* `main.go`
  CLI entrypoint.
* `cmd/`
  Command definitions grouped by domain.
* `pkg/app/`
  Application logic that wraps external tools and APIs.
* `internal/`
  Internal helpers such as logging, providers, system utilities, and instrumentation.
* `bin/`
  External helper binaries used by some commands.
* `out/`
  Compiled binaries.

## Main Command Groups

Commands are organized by domain.

Examples include:

* `aws` — helpers for ECR, EKS, Route53, and profiles
* `kubernetes` — utilities for pods, deployments, logs, labels, services, etc.
* `docker` — image and network helpers
* `dockerhub` — repository and tag inspection
* `git` — shortcuts for commits, branches, staging, and tags
* `gcp` — container registry helpers
* `terraform` / `terraformer` — infrastructure inspection utilities
* `helm` — release and revision helpers
* `istio` — sidecar and mesh inspection tools
* `files` — file operations such as compression, backup, and removal
* `yaml` — YAML manipulation utilities

## Usage

Commands follow a hierarchical pattern.

```
kts <domain> <resource> <action>
```

Examples:

```
kts kubernetes pods list
kts kubernetes deployments logs
kts docker images list
kts git commits list
kts aws ecr
```

Tab completion support is available.

## Installation

Build locally:

```
make build
```

The binary will be generated at:

```
out/kts
```

Add it to your PATH if desired.

## Notes

This repository evolves together with my workflows. Commands appear, change, or disappear as needed.
