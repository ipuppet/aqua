FROM node:22-alpine AS build
WORKDIR /app

COPY frontend/package.json frontend/package-lock.json ./
COPY frontend/apps ./apps
COPY frontend/packages ./packages
RUN npm ci

ARG VITE_API
ARG VITE_AUTH_COOKIE_DOMAIN
ENV NODE_ENV=production
ENV VITE_API=${VITE_API}
ENV VITE_AUTH_COOKIE_DOMAIN=${VITE_AUTH_COOKIE_DOMAIN}

RUN npm run build:admin && npm run build:blog && npm run build:passport

FROM nginx:1.29-alpine

COPY deploy/nginx/proxy.inc /etc/nginx/proxy.inc
COPY deploy/nginx/templates /etc/nginx/templates
COPY --from=build /app/apps/admin/dist /usr/share/nginx/html/admin
COPY --from=build /app/apps/blog/dist /usr/share/nginx/html/blog
COPY --from=build /app/apps/passport/dist /usr/share/nginx/html/passport

EXPOSE 80
