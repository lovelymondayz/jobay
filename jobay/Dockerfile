FROM node:20-slim

# Install git for build
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy everything
COPY backend/package*.json ./backend/
RUN cd /app/backend && npm install --omit=dev

COPY frontend/package*.json ./frontend/
RUN cd /app/frontend && npm install --omit=dev

COPY frontend/ ./frontend/
RUN cd /app/frontend && npm run build

COPY backend/ ./backend/

WORKDIR /app/backend

ENV PORT=3001
EXPOSE 3001

CMD ["node", "server.js"]
