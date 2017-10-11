FROM golang:alpine
RUN mkdir -p ./src/github.com/johnmcdnl/nfl-rank/nfl
ADD . ./src/github.com/johnmcdnl/nfl-rank
WORKDIR ./src/github.com/johnmcdnl/nfl-rank/cmd/nfl
RUN go build -o nflrank
CMD ./nflrank