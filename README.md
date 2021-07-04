# Ecology

__This project was superceded by [fns](https://github.com/gbdubs/fns), which is WIP and has not yet reached this feature parity.__

This repository holds utilities for creating and manipulating Cloud Resources within a layer of orchestration. This tool will eventually support a variety of high-level primitives that allow effictive and safe management of Cloud Resources on both GCP and AWS. Primitives I plan on building:

## Implemented Commands

No commands have been implemented yet.

## Planned Commands

### `initialize`

```
$ ecology initialize
```

`initialize` will walk the user through a variety of prompts that set up the environment that they will be using, including Github, AWS, Google and Travis credentials, a common parent directory for all ecology projects, and preferences surrounding geography and preferred runtime zones.

### Project Management

#### `create_project`

```
$ ecology create_project --project_name=MyFirstProject --lambda_name=HelloWorld --provider=AWS
```

`new_project` will construct a new project in the user's Ecology directory, with a lambda named HelloWorld (with a test), upload the Lambda to Amazon, verify that it runs correctly, setup a github repo for the project, set up a Travis continuous delivery for the created Github project, and verify that modifying the lambda results in an observed change in the serving lambda.

After `create_project` is invoked, a variety of other commands will be accessible on the project. 

#### `rename_project`

```
$ ecology rename_project --project_name=MyFirstProject --new_project_name=RenamedProject
```

`rename_project` will systematically change the name of the project across all surface areas (both on hosted service providers and on local machines). This operation will take a long time to complete.

#### `delete_project`

```
$ ecology delete_project --project_name=MyFirstProject
```

`delete_project` will carefully unwind all of the elements of the project that were constructed, both on local service providers, and on local machines

#### `pull_project`

```
$ ecology pull_project --project_name=MyFirstProject
```

`pull_project` will update the local configuration for the project by using the service providers as a source of truth (so if a cloud resource was deleted, pull project will move the state of the resource to be undeployed).

#### `push_project`

```
$ ecology push_project --project_name=MyFirstProject
```

`push_project` will push the local configuration of the project to the cloud service providers. This is by default what will be called by Travis on each update.


### Compute

#### `new_lambda`

```
$ ecology create_lambda --project_name=MyFirstProject --lambda_name=MySecondLambda --provider=AWS
```

`new_lambda` will construct a new Lambda within the existing project on AWS, and similar to `new_project`, set up continuous delivery. 

#### `rename_lambda`

```
$ ecology rename_lambda --project 

#### `delete_lambda`

```
$ ecology delete_lambda --project_name=MyFirstProject --lambda_name=MySecondLambda
```

`delete_lambda` will delete the lambda and all of its roles.

### Routing

#### `initialize_routing`

```
$ ecology initialize_routing --project_name=MyFirstProject --domain=mysubdomain.mydomain.com --registrar=google
```

`initialie_routing` will construct the local configuration of routing of URLs to resources, and configure the records of the domain with the provider to route to the appropriate service provider.

Further changes to routing will be reflected in the project's ecology manifest.

#### `update_domain_records`

```
$ ecology update_domain_records --project_name=MyFirstProject
```

`update_domain_records` will get the latest routing configurations from the service providers and update them at the domain registrar.

## Lol, why?

At google one of the most successful patterns I've seen is opinionated orchestration layers on top of powerful and flexible backing infrastructure. I've had trouble getting myself excited about cloud computing in the concrete, because of the breadth and depth of possibilities. I'm hoping this project gives me both concrete experience working with cloud service providers, and provids me a foundation on which I can quicly build resilient and scalable side projects. 

A key piece of what makes Google successful is that word "opinionated". There are thousands of ways to do any given task, but usually only one or two ways of doing it aligned with all mandates and recommendations and policies. I'm hoping to establish a similar degree of narrow rigor in my side projects. The only real opinionated choice I've made thus far has been to build compute on top of lambda. Serverless seems to offer exceptional benefits in terms of flexibility and scalability. The constraints it requires (no assumptions about serving location, volume or state) seem to force an engineer to make consistent choices that scale well. We'll see!

## Why Ecology?

The only problem in software is managing complexity - creating and managing systems to be stable, performant, constantly improving, and self-healing _without oversight_. That also seems to be the role of an ecologist - figuring out patterns, dependencies, and stability in a natural environment. I think when we look at software as an ecosystem it shifts our priorities toward stability, maintainance, and automation, and shifts emphasis away from features and rapid change. My hope in this project is to let myself explore these relationships, their interplay, and develop the skills to be a contientious steward of a technological ecosystem. 

