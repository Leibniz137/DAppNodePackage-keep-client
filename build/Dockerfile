FROM golang

ENV APPDIR /usr/src/app/
WORKDIR $APPDIR
COPY ./prepare.go $APPDIR
RUN go get github.com/ethereum/go-ethereum/ethclient && go build prepare.go

COPY ./config.toml.tmpl $APPDIR
COPY ./entrypoint.sh $APPDIR

ENTRYPOINT ["./entrypoint.sh"]
