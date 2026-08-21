FROM mcr.microsoft.com/playwright:v1.62.1-jammy

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY src/ ./src/

RUN mkdir -p /app/data/screenshots

ENV WORKER_PORT=3011
ENV SCREENSHOT_DIR=/app/data/screenshots

EXPOSE 3011

CMD ["node", "src/index.js"]
