# Go URL Shortener

## Overview
This is a simple URL shortener written in Go. It allows you to shorten long URLs into manageable links that never expire.

## How does a URL shortener work?
At a high level, the URL shortener executes the following operations:

- the server generates a unique short URL for each long URL
- the server encodes the short URL for readability
- the server persists the short URL in the data store
- the server redirects the client to the original long URL against the short URL
- 
## Features
- Shorten URLs
- Custom alias for URLs
- Basic analytics (click counts)

## REST API
```
POST /urls
GET /urls/:id
PUT /urls/:id
DELETE /urls/:id

```

## Structure of entities in the database
![image](https://github.com/yantay0/url-shortener/assets/93054482/0ccaa6ff-92d4-48ec-bb28-5d9cff66a734)



## High-level system design
![image](https://github.com/yantay0/url-shortener/assets/93054482/11b77a99-c41e-40ff-8710-24915dbfbc44)




### Prerequisites
- Go 1.16 or higher
- PostgreSQL 13.0 or higher

### Installation
1. Clone the repository:
```bash
git clone https://github.com/yantay0/url-shortener.git
