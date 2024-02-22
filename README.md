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


## Structure of entities in the database
![image](https://github.com/yantay0/url-shortener/assets/93054482/59f8bfac-fa69-4d1a-a1d6-36f5cde5ded5)



### Prerequisites
- Go 1.16 or higher
- PostgreSQL 13.0 or higher

### Installation
1. Clone the repository:
```bash
git clone https://github.com/yantay0/url-shortener.git
