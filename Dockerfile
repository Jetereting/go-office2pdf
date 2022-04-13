FROM alpine:3.15.4

WORKDIR /app

COPY fonts /usr/share/fonts
COPY office2pdf /

RUN apk add libreoffice && \
    fc-cache -fv

EXPOSE 3000

CMD [ "/office2pdf" ]