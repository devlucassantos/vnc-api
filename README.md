# vnc-api

🌍 *[English](README.md) ∙ [Português](README_pt.md)*

`vnc-api` is the service responsible for providing data to the frontend of the [Você na Câmara (VNC)](#você-na-câmara)
platform. In this repository, you will find the source code for the VNC API, which uses technologies such as Go, Echo,
PostgreSQL, and Redis. Additionally, the Docker container responsible for running this code is available,
allowing you to execute the project quickly and easily.

## How to run

### Prerequisites

To properly run `vnc-api`, you will need to have the [`vnc-databases`](https://github.com/devlucassantos/vnc-databases)
service containers running, so that this application's container has access to the necessary databases for reading and
writing data.

Additionally, you will also need to fill in some variables in the `.env` file, located in the _config_ directory
(`./src/config/.env`). In this file, you’ll notice that some variables are already filled in — this is because they
refer to default configurations, which can be used if you choose not to modify any of the pre-configured containers
used to run the repositories that make up VNC. However, feel free to change any of these variables if you wish to adapt
the project to your environment. Also note that some of these variables are not filled in — this is because their use is
tied to specific user accounts on platforms external to VNC, and therefore their values must be provided individually by
whoever intends to use these features. These variables are:

* `SMTP_HOST` → Address of the email server that will be used to send account activation codes, such as
  [Gmail](https://support.google.com/a/answer/176600?hl=en) or 
  [Outlook](https://support.microsoft.com/en-us/office/pop-imap-and-smtp-settings-for-outlook-com-d088b986-291d-42b8-9564-9c414e2aa040)
* `SMTP_PORT` → Port of the server used to send emails
* `SMTP_USER_EMAIL` → Email address of the account responsible for sending activation codes
* `SMTP_USER_PASSWORD` → Account password defined in `SMTP_USER_EMAIL`, according to the requirements of the email
  service used

### Running via Docker

To run the API, you will need to have [Docker](https://www.docker.com) installed on your machine and run the following
command in the root directory of this project:

````shell
docker compose up --build
````

### Documentation

After running the project, all the available routes for accessing the API can be found through the link:

> [http://localhost:8083/api/documentation](http://localhost:8083/api/documentation)

## Você na Câmara

Você na Câmara (VNC) is a news platform developed to simplify and make accessible the legislative propositions being
processed in the Chamber of Deputies of Brazil. Through the use of Artificial Intelligence, the platform synthesizes the
content of these legislative documents, transforming technical and complex information into clear and objective
summaries for the general public.

This project is part of the Final Paper of the platform's developers and was conceived based on architectures such as
hexagonal and microservices. The solution was organized into several repositories, each with specific responsibilities
within the system:

* [`vnc-databases`](https://github.com/devlucassantos/vnc-databases): Responsible for managing the platform's data
  infrastructure. Main technologies used: PostgreSQL, Redis, Liquibase, and Docker.
* [`vnc-pdf-content-extractor-api`](https://github.com/devlucassantos/vnc-pdf-content-extractor-api): Responsible for
  extracting content from the PDFs used by the platform. Main technologies used: Python, FastAPI, and Docker.
* [`vnc-domains`](https://github.com/devlucassantos/vnc-domains): Responsible for centralizing the platform's domains
  and business logic. Main technology used: Go.
* [`vnc-summarizer`](https://github.com/devlucassantos/vnc-summarizer): Responsible for the software that extracts data
  and summarizes the propositions available on the platform. Main technologies used: Go, PostgreSQL,
  Amazon Web Services (AWS), and Docker.
* [`vnc-api`](https://github.com/devlucassantos/vnc-api): Responsible for providing data to the platform's frontend.
  Main technologies used: Go, Echo, PostgreSQL, Redis, and Docker.
* [`vnc-web-ui`](https://github.com/devlucassantos/vnc-web-ui): Responsible for providing the platform's web interface.
  Main technologies used: TypeScript, SCSS, React, Vite, and Docker.
