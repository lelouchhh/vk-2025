# DB/Dockerfile

# Используем официальный образ PostgreSQL
FROM postgres:15-alpine

# Копируем файл инициализации базы данных
COPY ./init.sql /docker-entrypoint-initdb.d/

# Открываем порт для работы PostgreSQL
EXPOSE 5432
