# REST API Chat (Clean Architecture)
## Entities:
* User
* Chat
* Message

## Docker Compose
------------------------------------------------
Start APi with PostgreSQL <br />
*BridgeNetworking: API, PostgreSQL* <br />
*HostNetworking: API on `0.0.0.0:5000`* <br />
```bash
docker compose up -d
```
Start API with MySQL <br />
*BridgeNetworking: API, MySQL* <br />
*HostNetworking: API on `0.0.0.0:5000`* <br />
```bash
docker compose -f docker-compose.mysql.yml up -d
```
Start API with frontend <br />
_(This frontend does not fully implement the capabilities of the backend and serves as an example)_ <br />
*BridgeNetworking: Frontend, API, PostgreSQL* <br />
*HostNetworking: Frontend on `0.0.0.0:8080`* <br />
```bash
docker compose up -f docker-compose.psql.front.yml up -d
```
------------------------------------------------
## Docker Development Build
Statrt PostgreSQL
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=api_chat_db -p 5433:5432 -d postgres:15
```
Start MySQL
```bash
docker run --name mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_PASSWORD=mysql -e MYSQL_USER=mysql -e MYSQL_DATABASE=api_chat_db -p 3305:3306 -d mysql:8
```
Build API
```bash
docker build -t api_chat:v1.0.0 .
```
Build Frontend
```bash
docker build -t front_chat:v1.0.0 ./front_example/.
```
## Generate certificate
```bash
git clone https://github.com/redrockstyle/api_chat_golang.git
cd api_chat_golang/certs
chmod +x gen_localhost.sh
./gen_localhost.sh
```