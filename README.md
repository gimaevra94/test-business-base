# Repair Service

## Запуск
```bash
docker-compose up --build

Пользователи (из seed)
Dispatcher: ID 1 (Header: X-User-ID: 1, X-User-Role: dispatcher)
Master 1: ID 2 (Header: X-User-ID: 2, X-User-Role: master)
Master 2: ID 3 (Header: X-User-ID: 3, X-User-Role: master)
API
POST /requests - Создать
GET /requests?status=new - Список
PATCH /requests/?id=1&action=assign&master_id=2 - Назначить
PATCH /requests/?id=1&action=start - Взять в работу (Master)
PATCH /requests/?id=1&action=finish - Завершить (Master)
Проверка гонки (Race Condition)
Запустите в двух терминалах одновременно:
# Terminal 1
curl -X PATCH "http://localhost:8080/requests/?id=2&action=start" -H "X-User-ID: 2" -H "X-User-Role: master"

# Terminal 2
curl -X PATCH "http://localhost:8080/requests/?id=2&action=start" -H "X-User-ID: 3" -H "X-User-Role: master"
Один вернет 200 OK, второй 409 Conflict