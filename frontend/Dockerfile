# Etapa 1: Construcción (Build)
FROM node:20-alpine AS build-stage
WORKDIR /app
COPY package*.json ./
# npm ci es más estricto que npm install, respeta las versiones exactas del lockfile
RUN npm ci
COPY . .
RUN npm run build

# Etapa 2: Producción (Nginx)
FROM nginx:alpine AS production-stage
# Copiar el resultado de la compilación de la etapa 1
COPY --from=build-stage /app/dist /usr/share/nginx/html
# Sobrescribir la configuración por defecto de Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]