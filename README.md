# ns-scheduler

- [ns-scheduler](#ns-scheduler)
  - [Goal](#goal)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Internals](#internals)
      - [The watcher](#the-watcher)
      - [The suspender](#the-suspender)
    - [Flags](#flags)
    - [Resources](#resources)
    - [States](#states)
      - [On namespaces](#on-namespaces)
      - [On resources](#on-resources)
        - [Deployments and Stateful Sets](#deployments-and-stateful-sets)
        - [Cronjobs](#cronjobs)

## Goal

This controller is an attempt to scale down resources in a namespace in a "cron" way. 
You can for example scale down namespace resources from 8pm to 6am to save costs and compute resources.

Took a lot of inspiration from https://github.com/govirtuo/kube-ns-suspender

## Installation

A Dockerfile is available in the `manifests/deploy directory`. 
You can also use kubernetes manifests in the `manifests/deploy` to deploy the controller in your cluster. 
TODO: helm chart

## Usage

### Internals

This controller can be splitted into 2 parts:

* The watcher
* The suspender

#### The watcher

The watcher function is charged to check every 30 seconds (TODO: make it possible to change the value) all the namespaces. When it found namespace that have the `ns-scheduler: true` annotation, it sends it to the suspender. 

#### The suspender

The suspender function does all the work of reading namespaces/resources annotations, and (un)suspending them when required.

### Flags


| Environment variable   | Description  | Default  |
|---|---|---|
| `NS_SCHEDULER_LOG_LEVEL`  | Level of verbosity of the controller  | ""  |
| `NS_SCHEDULER_UPTIME_SCHEDULE`  | The default uptime schedule that the controller will use to scale up / scale down resources. During this interval, the controller will ensure that the resources are UP.  | `1-7 08:00-20:00`  |
### Resources

Currently supported resources are:

* [deployments](#deployments-and-stateful-sets)
* [stateful sets](#deployments-and-stateful-sets)
* [cronjobs](#cronjobs)

### States

Namespaces watched by `ns-scheduler` can be in 2 differents states:

* Running: the namespace is "up", and all the resources have the desired number of replicas.
* Suspended: the namespace is "paused", and all the supported resources are scaled down to 0 or suspended.

#### On namespaces

In order for a namespace to be watched by the controller, it needs to have the `ns-scheduler: true` annotation setted.

Then, the namespace will be attributed a state, which can be either `Running` or `Suspended` (depending if current date and time is in `UPTIME_SCHEDULE` or not).
You can override the global UPTIME_SCHEDULE for a specific namespace by adding the annotation `ns-scheduler/uptime` to a specific namespace. The [suspender](#the-suspender) will check if the namespace should be suspensed based on the namespace annotation instead of the global env variable `NS_SCHEDULER_UPTIME_SCHEDULE`.

#### On resources

Annotations are employed to save the original state of a resource. 

##### Deployments and Stateful Sets

As those resources have a `spec.replicas` value, they must have a `ns-scheduler/originalReplicas` annotation that must be the same as the `spec.replicas` value. This annotation will be used when a resource will be "unsuspended" to set the original number of replicas.

##### Cronjobs

Cronjobs have a `spec.Suspend` value that indicates if they must be runned or not. As this value is a boolean, **no other annotations are required**.