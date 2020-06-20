FROM tiangolo/uwsgi-nginx-flask:python3.7

WORKDIR /app
RUN apt update

RUN groupadd -g 2000 verimapcom && useradd -u 2000 -g verimapcom verimapcom
RUN apt-get update -yq \
    && apt-get install curl gnupg -yq \
    && curl -sL https://deb.nodesource.com/setup_12.x | bash \
    && apt-get install nodejs -y


RUN npm install -g yarn
RUN yarn global add webpack webpack-cli

COPY requirements.txt .
# RUN pip3 install uwsgi
RUN pip3 install --no-cache-dir -r requirements.txt

COPY yarn.lock .
COPY package.json .
RUN yarn install
COPY . .
ENV STATIC_PATH /app/dist
ENV STATIC_URL /dist
RUN yarn run build
RUN ip addr
