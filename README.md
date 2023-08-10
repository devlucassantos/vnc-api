# vnc-read-api

🌍 *[English](README.md) ∙ [Português](README_pt.md)*

`vnc-read-api` is the repository responsible for reading data from the [Você na Câmara (VNC)](#você-na-câmara-vnc)
platform databases. In this repository you will find the source code of the VNC Read API and also the container
responsible for executing this code, so you can easily run the project.

## How to run

> Note that to properly run `vnc-read-api` you will need to have the [`vnc-database` containers](https://github.com/devlucassantos/vnc-database)
running so that this application's container has access to the databases needed to query the data.

To build the databases you will need to have [Docker](https://www.docker.com) installed on your machine and run the
following commands in the root directory of this project:

````shell
docker compose up
````

## Você Na Câmara (VNC)

Você Na Câmara (VNC) is a news platform that seeks to simplify the proposals under debate in the Chamber of Deputies of
Brazil aiming to synthesize the ideas of these propositions through the use of Artificial Intelligence (AI) so that
these documents can have their ideas expressed in a simple and objective way for the general population.
