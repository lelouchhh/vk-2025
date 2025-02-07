# Стадия сборки
FROM node:18-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и package-lock.json (или yarn.lock)
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install --legacy-peer-deps

# Копируем все остальные файлы
COPY . .

# Собираем приложение для продакшн
RUN npm run build

# Стадия выполнения
FROM node:18-alpine AS runner

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем только production-зависимости
COPY package*.json ./
RUN npm install --production --legacy-peer-deps

# Копируем собранные файлы из стадии сборки
COPY --from=builder /app/dist ./dist

# Указываем порт
EXPOSE 3001

# Запускаем приложение
CMD ["npm", "run", "start"]
