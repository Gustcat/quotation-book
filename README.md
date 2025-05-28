# quotation-book
Мини-сервис для хранения и управления цитатами, реализованный на Go.
Сервис предоставляет HTTP API для добавления, просмотра, фильтрации и удаления цитат.

В проекте используются только стандартные библиотеки Go, за исключением gorilla/mux для маршрутизации.
Все данные хранятся в оперативной памяти (in-memory storage), без использования внешней базы данных.

## Развертывание
Клонировать репозиторий
```
git clone git@github.com:Gustcat/people-info-service.git
```
Скомпилировать и заупстить бинарный файл
```
cd quotation-book/cmd
go run main.go
```
Сервис будет доступен по адресу _http://localhost:8080_

## Примеры запросов
- Создание записи:
```bash
curl -X POST http://localhost:8080/quotes \ -H "Content-Type: application/json" \ -d
'{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}
# response:
{
    "id": 1
}
```
- Получение всех цитат:
```bash
curl http://localhost:8080/quotes
# response
[
    {
        "id": 1,
        "author": "Confucius",
        "quote": "Life is simple, but we insist on making it complicated."
    },
    {
        "id": 2,
        "author": "Second",
        "quote": "Life is simple, but we insist on making it complicated."
    },
    {
        "id": 3,
        "author": "Second",
        "quote": "Life is simple."
    }
]
```
- Фильтрация по автору:
```bash
curl http://localhost:8080/quotes?author=Confucius
# response
[
    {
        "id": 1,
        "author": "Confucius",
        "quote": "Life is simple, but we insist on making it complicated."
    }
] 
```
- Получение рандомной цитаты:
```bash
curl http://localhost:8080/quotes/random
# response
{
    "id": 3,
    "author": "Second",
    "quote": "Life is simple."
} 
```
- Удаление цитаты:
```bash
curl -X DELETE http://localhost:8080/quotes/1
```

## Статус проекта
Проект находится в стадии разработки.
## Технологии
- Go 1.23

## Автор
https://github.com/Gustcat