# loadbalancer

## How to run this project

- clone this repository using **git clone https://github.com/saurabhsisodia/loadbalancer.git**
- make sure Go(1.18) is installed in your system
- **mkdir loadbalancer**
- Run **go run main.go**, make sure that port 8080 is free in your system

## How to use this Load Balancer

- make **POST** request using postman or any client to **localhost:8080/urls/register** to register your URLs in LB
  - Request body should be array of urls
  ```
  [
  	{"url":"https://www.google.com"},
  	{"urls":"https://www.abc.com"},
  	{"urls":"efefeferf"},
  	------
  	------
  ]
  ```
- make HTTP request to **localhost:8080/proxy** to get response from one of the registered URL in Round Robin fashion