# Используем многостадийную сборку
FROM node:22.8-alpine AS builder

WORKDIR /app

# Копируем зависимости
COPY frontend/src/package*.json ./
RUN npm ci

# Копируем исходники и собираем
COPY frontend/src/ .
RUN npm run build

# Финальный образ
FROM node:22.8-alpine

WORKDIR /app
# Копируем ТОЛЬКО собранное приложение
COPY --from=builder /app/.output ./

# Правильный путь к серверу
CMD ["node", "./server/index.mjs"]