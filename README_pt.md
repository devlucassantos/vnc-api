# vnc-api

🌍 *[English](README.md) ∙ [Português](README_pt.md)*

`vnc-api` é o serviço responsável por disponibilizar os dados para o frontend da plataforma
[Você na Câmara (VNC)](#você-na-câmara). Neste repositório, você encontrará o código-fonte da API da VNC, que
utiliza tecnologias como Go, Echo, PostgreSQL e Redis. Além disso, está disponível o container Docker responsável
por executar este código, permitindo que você execute o projeto de forma simples e rápida.

## Como Executar

### Pré-requisitos

Para executar corretamente o `vnc-api` você precisará ter os containers do serviço
[`vnc-databases`](https://github.com/devlucassantos/vnc-databases) em execução, de modo que o container desta aplicação
tenha acesso aos bancos de dados necessários para a leitura e escrita dos dados.

Além disso, você precisará preencher também algumas variáveis do arquivo `.env`, localizado no diretório _config_
(`./src/config/.env`). Neste arquivo, você notará que algumas variáveis já estão preenchidas — isso ocorre porque se
referem a configurações padrão, que podem ser utilizadas caso você opte por não modificar nenhum dos containers
pré-configurados para rodar os repositórios que compõem a VNC. No entanto, sinta-se à vontade para alterar qualquer uma
dessas variáveis, caso deseje adaptar o projeto ao seu ambiente. Observe também que algumas destas variáveis não estão
preenchidas - isso ocorre porque seu uso está vinculado a contas específicas de cada usuário em plataformas externas ao
VNC e, portanto, seus valores devem ser fornecidos individualmente por quem deseja utilizar esses recursos. Essas
variáveis são:

* `SMTP_HOST` → Endereço do servidor do serviço de e-mail que será utilizado para o envio dos códigos de ativação de
  conta, como [Gmail](https://support.google.com/a/answer/176600?hl=pt-BR) ou
  [Outlook](https://support.microsoft.com/pt-br/office/configura%C3%A7%C3%B5es-pop-imap-e-smtp-para-outlook-com-d088b986-291d-42b8-9564-9c414e2aa040)
* `SMTP_PORT` → Porta do servidor utilizado para o envio dos e-mails
* `SMTP_USER_EMAIL` → Endereço de e-mail da conta responsável pelo envio dos códigos de ativação
* `SMTP_USER_PASSWORD` → Senha da conta definida em `SMTP_USER_EMAIL`, conforme os requisitos do serviço de e-mail 
  utilizado

### Executando via Docker

Para executar a API, você precisará ter o [Docker](https://www.docker.com) instalado na sua máquina e executar o
seguinte comando no diretório raiz deste projeto:

````shell
docker compose up --build
````

### Documentação

Após a execução do projeto, todas as rotas disponíveis para acesso à API podem ser encontradas através do link:

> [http://localhost:8083/api/documentation](http://localhost:8083/api/documentation)
<img width="2880" height="1800" alt="image" src="https://github.com/user-attachments/assets/6b623d88-0e84-4f99-9621-bb87a2d0a1db" />

## Você na Câmara

Você na Câmara (VNC) é uma plataforma de notícias desenvolvida para simplificar e tornar acessíveis às proposições
legislativas que tramitam na Câmara dos Deputados do Brasil. Por meio do uso de Inteligência Artificial, a plataforma
sintetiza o conteúdo desses documentos legislativos, transformando informações técnicas e complexas em resumos objetivos
e claros para a população em geral.

Este projeto integra o Trabalho de Conclusão de Curso dos desenvolvedores da plataforma e foi concebido com base
em arquiteturas como a hexagonal e a de microsserviços. A solução foi organizada em diversos repositórios, cada um com
responsabilidades específicas dentro do sistema:

* [`vnc-databases`](https://github.com/devlucassantos/vnc-databases): Responsável por gerenciar a infraestrutura de
  dados da plataforma. Principais tecnologias utilizadas: PostgreSQL, Redis, Liquibase e Docker.
* [`vnc-pdf-content-extractor-api`](https://github.com/devlucassantos/vnc-pdf-content-extractor-api): Responsável por
  realizar a extração de conteúdo dos PDFs utilizados pela plataforma. Principais tecnologias utilizadas: Python,
  FastAPI e Docker.
* [`vnc-domains`](https://github.com/devlucassantos/vnc-domains): Responsável por centralizar os domínios e regras de
  negócio da plataforma. Principal tecnologia utilizada: Go.
* [`vnc-summarizer`](https://github.com/devlucassantos/vnc-summarizer): Responsável pelo software que extrai os dados e
  sumariza as proposições disponibilizadas na plataforma. Principais tecnologias utilizadas: Go, PostgreSQL, Amazon Web
  Services (AWS) e Docker.
* [`vnc-api`](https://github.com/devlucassantos/vnc-api): Responsável por disponibilizar os dados para o frontend da
  plataforma. Principais tecnologias utilizadas: Go, Echo, PostgreSQL, Redis e Docker.
* [`vnc-web-ui`](https://github.com/devlucassantos/vnc-web-ui): Responsável por fornecer a interface web da plataforma.
  Principais tecnologias utilizadas: TypeScript, SCSS, React, Vite e Docker.
