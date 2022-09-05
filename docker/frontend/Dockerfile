FROM node:14.20.0 AS build

WORKDIR /srv/

COPY frontend/package.json /srv/package.json
COPY frontend/package-lock.json /srv/package-lock.json

RUN npm ci

COPY frontend/tsconfig.json /srv/tsconfig.json
COPY frontend/src /srv/src
COPY frontend/public /srv/public

RUN npm run build

FROM nginx:1.23.1

COPY docker/frontend/default.conf /etc/nginx/conf.d/default.conf

COPY --from=build /srv/build/static /usr/share/nginx/html/static
COPY --from=build /srv/build/* /usr/share/nginx/html/
