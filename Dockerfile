FROM alpine:latest

WORKDIR /app

COPY app .
RUN mkdir log
RUN mkdir data
RUN mkdir data/secret
COPY log .
COPY data/secret/tokensecret data/secret/tokensecret

RUN  chmod 777 app

EXPOSE 8081
EXPOSE 8090

CMD [ "./app" ]