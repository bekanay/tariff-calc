# --- build stage ---
FROM node:22-alpine AS build
WORKDIR /app
COPY package*.json ./
# Use a stable mirror to dodge npmjs 500s and keep moderate retries.
RUN npm config set registry https://registry.npmmirror.com \
    && npm config set fetch-retries 4 \
    && npm config set fetch-retry-factor 2 \
    && npm config set fetch-retry-mintimeout 2000 \
    && npm config set fetch-retry-maxtimeout 30000 \
    && npm ci
COPY . .
ARG VITE_URL
ENV VITE_URL=$VITE_URL
RUN npm run build

# --- runtime stage ---
FROM nginx:1.27-alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
