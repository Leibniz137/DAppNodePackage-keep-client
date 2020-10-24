FROM golang

ENV APPDIR /usr/src/app/
WORKDIR $APPDIR
COPY ./prepare.go $APPDIR

# disable cgo to fix this: `./entrypoint.sh: line 3: ./prepare: not found`
# see: https://stackoverflow.com/a/36308464/1527293
ENV CGO_ENABLED=0
RUN go get github.com/ethereum/go-ethereum/ethclient \
  && go get github.com/ethereum/go-ethereum/accounts \
  && go get github.com/ethereum/go-ethereum/accounts/keystore \
  && go build prepare.go

COPY ./config.toml.tmpl $APPDIR
COPY ./entrypoint.sh $APPDIR

FROM keepnetwork/keep-client:v1.3.1

ENV APPDIR /usr/src/app/
COPY --from=0 $APPDIR $APPDIR

WORKDIR $APPDIR

ENTRYPOINT ["./entrypoint.sh"]
