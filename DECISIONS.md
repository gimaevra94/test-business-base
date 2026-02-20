```markdown
1. **DB**: MySQL выбран как требование.
2. **Concurrency**: Optimistic Locking (column `version`). Позволяет вернуть 409 сразу без блокировок БД на долгий срок.
3. **Auth**: Header-based (X-User-ID). Упрощено для тестового задания.
4. **API**: JSON REST. Удобно для автоматических тестов и curl.
5. **Migrations**: Отдельный контейнер `migrate` гарантирует схему перед стартом приложения.
6. **Structure**: Monolith inside Docker. Достаточно для масштаба задачи.
7. **Tests**: Unit + Integration (через compose).