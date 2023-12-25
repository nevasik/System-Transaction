# Тестовое задание: Транзакционная система на Golang

Это тестовое задание представляет собой реализацию транзакционной системы на языке Golang с использованием технологий NATS (как брокера сообщений) и PostgreSQL (в качестве базы данных).

## Установка и запуск NATS из Docker

1. Установите Docker, если у вас его нет: [Docker Install](https://docs.docker.com/get-docker/)

2. Запустите NATS из Docker:

   ```bash
   docker run -p 4222:4222 -p 6222:6222 -p 8222:8222 --name nats-main -d nats

Эта команда запустит контейнер NATS с именем "nats-main" и откроет порты для внешнего доступа.

## Установка и настройка PostgreSQL из Docker

1. Установите Docker, если у вас его нет: [Docker Install](https://docs.docker.com/get-docker/)

2. Запустите PostgreSQL из Docker:

   ```bash
   docker run --name test-postgres -p 5432:5432 -e POSTGRES_PASSWORD=admin -d postgres

Эта команда запустит контейнер PostgreSQL с именем "test-postgres" и откроет порты для внешнего доступа.

## Запуск приложения

1. Склонируйте репозиторий :
```bash 
  git clone [https://github.com/nevasik/System-Transaction]
  cd System-Transaction
 ```
2. Установите зависимости :
```bash 
    go mod download
```
3. Настройте переменные окружения в файле .env в корне проекта
  ```bash 
DB_USER=
DB_PASSWORD=admin
DB_NAME=postgres
  ```  
4. Запустите приложение
```bash 
go run cmd/transaction-system/main.go
  ```

### Использование API
Ваше приложение теперь готово к использованию. Вы можете взаимодействовать с API, используя ручки, такие как /invoice и /withdraw для создания транзакций, а также /balances для получения актуального и замороженного баланса клиентов.

