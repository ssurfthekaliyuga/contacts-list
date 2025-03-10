FROM node:22.8-alpine AS builder

WORKDIR /app

COPY frontend/src/package*.json ./
RUN npm ci

COPY frontend/src/ .
RUN npm run build

FROM node:22.8-alpine

USER node

WORKDIR /home/node/app

COPY --from=builder /app/.output ./

CMD ["node", "./server/index.mjs"]