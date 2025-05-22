# Finala ![Lint](https://github.com/similarweb/finala/workflows/Lint/badge.svg) ![Fmt](https://github.com/similarweb/finala/workflows/Fmt/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/similarweb/finala)](https://goreportcard.com/report/github.com/similarweb/finala) [![Gitter](https://badges.gitter.im/similarweb-finala/community.svg)](https://gitter.im/similarweb-finala/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

**Note**: The `master` branch represents the latest developed version and it may be in an *unstable or even broken*.

In order to get the latest stable version please use the [releases pages](https://github.com/similarweb/finala/releases).

![alt Logo](https://raw.githubusercontent.com/similarweb/finala/master/docs/images/main-logo.png)
![Finala Processing](https://raw.githubusercontent.com/similarweb/finala/master/docs/images/finala.png)

----

## Overview

Finala is an open-source resource cloud scanner that analyzes, discloses, presents and notifies about wasteful and unused resources.

With Finala you can achieve 2 main objectives: **Cost saving & Unused resources detection**.

## Features

* **YAML Definitions**: Resources definitions are described using a high-level YAML configuration syntax. This allows Finala consumers easily tweak the configuration to help it understand their infrastructure, spending habits and normal usage.
* **1 Click Deployment**: Finala can be deployed via Docker compose or a [Helm chart](https://github.com/similarweb/finala-helm).
* **Graphical user interface**: Users can easily explore and investigate unused or unutilized resources in your cloud provider.
* **Resource Filtering by Cloud Provider Tags**: Users can filter unused resources by just providing the tags you are using in your cloud provider.
* **Schedule Pro Active Notifications**: Finala has the ability to configure scheduled based notifications to a user or a group.

## Supported Services

### Finala's Definitions

* **Potential Cost Optimization** - is the price you can save for untilized resources in your infrastructure
* **Unused Resource** - are resources which don't necessarily cost money and can be removed.

### AWS

Resource            | Potential Cost Optimization| Unused Resource         |
--------------------| ---------------------------|-------------------------|
API Gateway         | :heavy_minus_sign:         | :ballot_box_with_check:
DocumentDB          | :ballot_box_with_check:    | :heavy_minus_sign:
DynamoDB            | :ballot_box_with_check:    | :heavy_minus_sign:
EC2 ALB,NLB         | :ballot_box_with_check:    | :heavy_minus_sign:
EC2 Elastic IPs     | :ballot_box_with_check:    | :heavy_minus_sign:
EC2 ELB             | :ballot_box_with_check:    | :heavy_minus_sign:
EC2 NAT Gateways    | :ballot_box_with_check:    | :heavy_minus_sign:
EC2 Instances       | :ballot_box_with_check:    | :heavy_minus_sign:
EC2 Volumes         | :ballot_box_with_check:    | :heavy_minus_sign:
ElasticCache        | :ballot_box_with_check:    | :heavy_minus_sign:
ElasticSearch       | :ballot_box_with_check:    | :heavy_minus_sign:
IAM User            | :heavy_minus_sign:         | :ballot_box_with_check:
Kinesis             | :ballot_box_with_check:    | :heavy_minus_sign:
Lambda              | :heavy_minus_sign:         | :ballot_box_with_check:
Neptune             | :ballot_box_with_check:    | :heavy_minus_sign:
RDS                 | :ballot_box_with_check:    | :heavy_minus_sign:
RedShift            | :ballot_box_with_check:    | :heavy_minus_sign:

## Recent Project Upgrades

This project has recently undergone significant upgrades to modernize its stack and improve its core functionalities. Key changes include:

*   **Search Backend Modernization**: Migrated from Elasticsearch to Meilisearch for an improved search experience.
*   **Go Version Update**: Upgraded to the latest Go version along with its dependencies, ensuring better performance and security.
*   **Frontend Overhaul**:
    *   Upgraded React to v18.
    *   Updated Material-UI (MUI) to v5.
    *   Migrated to React Router v6.
    *   Upgraded Webpack to v5 and updated related loaders and plugins for a more efficient build process.
*   **Authentication & Security**:
    *   Added secure login functionality with JWT-based authentication.
    *   Implemented protected routes for authenticated users.
    *   Enhanced API security with authentication middleware.
*   **Containerization Improvements**: Updated Dockerfiles for both development and production environments, optimizing build layers and improving security by using non-root users and newer base images.
*   **General Dependency Updates**: Various other packages and dependencies across the project have been updated to their latest stable versions.

## Authentication

Finala now features a secure authentication system to protect access to your cloud resource information.

### Login System

* **User Interface**: A clean, modern login page that matches Finala's visual identity.
* **API Authentication**: Backend API routes are now protected with JWT authentication.
* **Protected Routes**: Access to resource data requires successful authentication.

### Configuration

Authentication credentials are configured in `/etc/finala/config.yaml` under the `auth` section:

```yaml
auth:
  username: admin
  password: your-secure-password
```

### Auto-Generated Credentials

If the configuration file is missing, incomplete, or authentication credentials are not set:

1. The system will use `admin` as the default username.
2. A secure random password will be automatically generated at startup.
3. The generated password will be displayed in the startup logs for initial access.
4. The configuration will be written to `/etc/finala/config.yaml` for future reference.

Example startup log with generated credentials:
```
INFO: Generated random password for admin user: XyzT7q2PwC8rLzV5
INFO: To use custom credentials, set auth.username and auth.password in /etc/finala/config.yaml
```

For security in production environments, it is recommended to set your own credentials in the configuration file.

## QuickStart

Follow the [quick start](https://finala.io/docs/getting-started/quick-start) in our documentation to get familiar with Finala.


## Web User Interface

You can access Finala's user interface via http://localhost:8080/  (After you have finished with the quick start guide)
![dashboard](https://raw.githubusercontent.com/similarweb/finala/master/docs/images/main-dashboard.png)

## Installation

Please refer to [Installation instructions](https://finala.io/docs/installation/getting-started).

## Documentation & Guides

Documentation is available on the Finala website [here](https://finala.io/).

## Community, discussion, contribution, and support

You can reach the Finala community and developers via the following channels:

* [Gitter Community](https://gitter.im/similarweb-finala/community):
  * [finala-users](https://gitter.im/similarweb-finala/users)
  * [finala-developers](https://gitter.im/similarweb-finala/developers)

## Contributing

Thank you for your interest in contributing! Please refer to [Contribution guidelines](https://finala.io/docs/contributing/submitting-pr) for guidance.
