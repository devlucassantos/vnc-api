# vnc-read-api

🌍 *[English](README.md) ∙ [Português](README_pt.md)*

`vnc-read-api` é o repositório responsável por realizar a leitura dos dados nos bancos de dados da plataforma
[Você na Câmara (VNC)](#você-na-câmara-vnc). Neste repositório você encontrará o código-fonte da API de leitura do VNC e
também o container responsável por executar este código, deste modo você poderá facilmente rodar o projeto.

## Como Executar

> Observe que para executar corretamente o `vnc-read-api` você precisará ter os [containers do `vnc-database`](https://github.com/devlucassantos/vnc-database)
em execução de modo que o container desta aplicação tenha acesso aos bancos de dados necessários para a consulta dos dados.

Para executar a API você precisará ter o [Docker](https://www.docker.com) instalado na sua máquina e executar o seguinte
comando no diretório raiz deste projeto:

````shell
docker compose up
````

## Você Na Câmara (VNC)

Você na Câmara (VNC) é uma plataforma de notícias que busca simplificar as proposições que tramitam pela Câmara dos
Deputados do Brasil visando sintetizar as ideias destas proposições através do uso da Inteligência Artificial (IA)
de modo que estes documentos possam ter suas ideias expressas de maneira simples e objetiva para a população em geral.
