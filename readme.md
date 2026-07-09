# 🎟️ **Online Subscription Service**

### Сервис управления подписками.

* **Создание / просмотр / обновление / удаление подписок (CRUDL)**
* **Подсчет суммарной стоимости подписок за период**
* **Swagger/OpenAPI документация**
* **Логи через [Uber Zap](https://github.com/uber-go/zap)**
* **Автоматические миграции в PostgreSQL**

---

> ⚠️ Миграции применяются автоматически при старте приложения. Внутри контейнера с базой данных уже будет готова таблица
`subscriptions`.
> Файл `.env` выступает как здесь как пример для более удобного развертывания из GitHub.
---

## 📘 **API документация**

| Тип          | URL                                                                                  |
|--------------|--------------------------------------------------------------------------------------|
| Swagger UI   | [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) |
| OpenAPI JSON | [http://localhost:8080/swagger/doc.json](http://localhost:8080/swagger/doc.json)     |

---

## 🗃️ **Миграции**

* Папка: `migrations/`
* Автозапуск при старте приложения
* Таблица `subscriptions` создается автоматически

---

## 🌐 **Примеры HTTP запросов**

### Создание подписки

```http
POST http://localhost:8080/subscriptions
Content-Type: application/json

{
  "service_name": "Netflix",
  "monthly_price": 999,
  "start_date": "12-2025",
  "user_id": "54639c13-710c-48f1-80b0-d18e88a6e9f5"
}
```

### Получение всех подписок

```http
GET http://localhost:8080/subscriptions
```

### Получение подписки по ID

```http
GET http://localhost:8080/subscriptions/{id}
```

### Обновление подписки

```http
PATCH http://localhost:8080/subscriptions/{id}
Content-Type: application/json

{
  "service_name": "Spotify",
  "monthly_price": 1200,
  "start_date": "01-2026"
}
```

### Удаление подписки

```http
DELETE http://localhost:8080/subscriptions/{id}
```

### Суммарная стоимость подписок

```http
GET http://localhost:8080/subscriptions/summary?from=01-2025&user_id=54639c13-710c-48f1-80b0-d18e88a6e9f5&service_name=Netflix
```

> 💡 Больше готовых примеров запросов есть в `requests/request.http`. 

