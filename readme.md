## Jobsity Dev Test

This repository contains implmenentation of tech test in [go-challenge-financial-chat_5cd0c06df1e48.pdf](./go-challenge-financial-chat_5cd0c06df1e48.pdf).

### Dependencies

[Docker](https://www.docker.com/) with support for [compose](https://docs.docker.com/compose/).

### Execution

1. Clone this repository.
2. Run the services with command `docker compose up -d`.
3. Open `http://0.0.0.0:9090` in your browser.

### Services

This project contains main 3 services, chat server, web interface and bot, and supported by mongodb as database persistence and redis as pubsub.

### Chat server

Chat server is a backend that provide chat resources by API. The endpoint is available in url `http://0.0.0.0:9001`, with documentation below.

<details>
 <summary><code>GET</code> <code><b>/api/status</b></code> <code> - </code> <code>Application status</code></summary>

##### Response

> | Status code | Response |
> |-----------|-----------|
> | `200` | Status with datetime |

##### Example

> ```bash
> curl --request GET --url 'http://127.0.0.1:9007/api/status'
> ```

</details>

<details>
 <summary><code>POST</code> <code><b>/api/register</b></code> <code> - </code> <code>Register new user</code></summary>

##### Parameters

> | Name | Type | Description |
> |-----------|-----------|-----------|
> | name | string | Name of user |
> | email | string | E-mail of user |
> | password | string | Password with minimum lenght of 8 characters |

##### Response

> | Status code | Response |
> |-----------|-----------|
> | `201` | User created with success |
> | `400` | One of parameters above is incorret |
> | `422` | User already registered with e-mail filled |
> | `500` | Internal server error |

##### Example

> ```bash
> curl --request POST \
>   --url http://127.0.0.1:9001/api/register \
>   --header 'Content-Type: application/json' \
>   --data '{
>   "name": "John Doe",
>   "email": "john@doe.tld",
>   "password": "pa$sw0rd"
> }'
> ```

> ```json
> {
>   "id": "65207247d040de340e853cc9",
>   "name": "John Doe",
>   "email": "john@doe.tld",
>   "password": "it's a secret :)"
> }
> ```

</details>

<details>
 <summary><code>POST</code> <code><b>/api/login</b></code> <code> - </code> <code>Authentication of user</code></summary>

##### Parameters

> | Name | Type | Description |
> |-----------|-----------|-----------|
> | email | string | E-mail of user |
> | password | string | Password with minimum lenght of 8 characters |

##### Response

> | Status code | Response |
> |-----------|-----------|
> | `200` | Login with success and access token generated |
> | `400` | One of parameters above is incorret |
> | `422` | E-mail or password is incorrect |
> | `500` | Internal server error |

##### Example

> ```bash
> curl --request POST \
>   --url http://127.0.0.1:9001/api/login \
>   --header 'Content-Type: application/json' \
>   --data '{
>   "email": "john@doe.tld",
>   "password": "pa$sw0rd"
> }'
> ```

> ```json
> {
>   "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6IkpvaG4gRG9lIiwidXNlckVtYWlsIjoiam9obkBkb2UudGxkIiwiZXhwIjoxNjk3MDM0MzU4fQ.Cyrr_QJQdKOgtNcW2jF9_UFxFSe8StgkOc_OLTqkzKg"
> }
> ```

</details>

<details>
 <summary><code>POST</code> <code><b>/api/messages</b></code> <code> - </code> <code>Create new message</code></summary>

##### Parameters

> | Name | Type | Format | Description |
> |-----------|-----------|-----------|-----------|
> | Authorization | string | Header | Token of authenticated user |
> | messageText | string | Json body | Text message |

##### Response

> | Status code | Resposta |
> |-----------|-----------|
> | `201` | Message created with success |
> | `202` | Message accepted with success if a message contains a command |
> | `400` | One of parameters above is incorret |
> | `403` | Authenticated user required |
> | `500` | Internal server error |

##### Example

> ```bash
> curl --request POST \
>   --url http://127.0.0.1:9001/api/messages \
>   --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6IkpvaG4gRG9lIiwidXNlckVtYWlsIjoiam9obkBkb2UudGxkIiwiZXhwIjoxNjk3MDM0MzU4fQ.Cyrr_QJQdKOgtNcW2jF9_UFxFSe8StgkOc_OLTqkzKg' \
>   --header 'Content-Type: application/json' \
>   --data '{
>   "messageText": "Did you see that?"
> }'
> ```

> ```json
> {
>   "id": "01HC89B67PBJT8JXMMARZGAP1F",
>   "userName": "John Doe",
>   "userEmail": "john@doe.tld",
>   "messageText": "Did you see that?",
>   "datetime": "2023-10-08 15:38:41"
> }
> ```

</details>

<details>
 <summary><code>GET</code> <code><b>/api/messages</b></code> <code> - </code> <code>Get list of messages</code></summary>

##### Parameters

> | Name | Type | Format | Description |
> |-----------|-----------|-----------|-----------|
> | Authorization | string | Header | Token of authenticated user |
> | page | integer | Query string | Page of messages (optional) |
> | itemsPerPage | integer | Query string | Items per page of messages (optional) |

##### Response

> | Status code | Resposta |
> |-----------|-----------|
> | `200` | List of messages |
> | `400` | One of parameters above is incorret |
> | `403` | Authenticated user required |
> | `500` | Internal server error |

##### Example

> ```bash
> curl --request GET \
>   --url 'http://127.0.0.1:9001/api/messages?page=1&itemsPerPage=50' \
>   --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6IkpvaG4gRG9lIiwidXNlckVtYWlsIjoiam9obkBkb2UudGxkIiwiZXhwIjoxNjk3MDM0MzU4fQ.Cyrr_QJQdKOgtNcW2jF9_UFxFSe8StgkOc_OLTqkzKg'
> ```

> ```json
> [
>   {
>     "id": "01HC89B67PBJT8JXMMARZGAP1F",
>     "userName": "John Doe",
>     "userEmail": "john@doe.tld",
>     "messageText": "Did you see that?",
>     "datetime": "2023-10-08 15:38:41"
>   }
> ]
> ```

</details>

### Bot

Bot is a project that provides a worker responsible to receive stock commands and retrieve quote for stock filled in command. This bot receives and dispatch data by events.

### Web

Web is a simple web interface made for usage purposes. This project is available at `http://0.0.0.0:9090`.

### Project explained

This project use concepts and strategies of layered architecture to provide a separation between business code from other code that support application, like application and infra. I used Domain Driven Design for this implementation, with Repository pattern, Usecases, Data Transfer Objects and Presenter to create code that it's easy to maintain, add new features, fixes, tests and scale. I used Table driven tests approach in various tests for testing all possible results by mocking the behavior of dependencies. This project use concepts of Event driven for dispatch events that can be consumed for other projects.

For communication between services from Event driven approach, I used Redis as pubsub, then bot project can receive and dispatch events for the commands received. The bot project is completely decoupled, then if necessary, code can be moved to a new repository and continue as standalone project, without any problems, just configure the pubsub topic for consume and dispatch events.

Tech test mentioned that front-end isn't the focus, then I used polling approach for retrieving new messages from chat server.
