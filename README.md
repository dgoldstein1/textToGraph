# crawler

Script to crawl html and add href links, crawling and indexing 5k sites / second into [big-data graph DB](https://github.com/dgoldstein1/graphApi).

![Maintainability]
![Test Coverage]
![CircleCI]

## Build it

#### Binary

```sh
go get -u github.com/dgoldstein1/textToGraph
```

#### Docker
```sh
docker build . -t dgoldstein1/wikipedia-path
```

## Run it


```sh
docker-compose up -d
```

or with dependencies running locally

```sh
export GRAPH_DB_ENDPOINT="http://localhost:5000" # endpoint of graph database
export TWO_WAY_KV_ENDPOINT="http://localhost:5001" # endpoint of k:v <-> v:k 
export MAX_APPROX_NODES=1000 # approximate number of nodes to visit (+/- one order of magnitude), set to '-1' for unlimited pars
textToGraph parse ./documents/moby_dick.txt
```

## Development

#### Local Development

- Install [inotifywait](https://linux.die.net/man/1/inotifywait)
```sh
./watch_dev_changes.sh
```

#### Testing

```sh
go test $(go list ./... | grep -v /vendor/)
```

## Authors

* **David Goldstein** - [DavidCharlesGoldstein.com](http://www.davidcharlesgoldstein.com/?github-textToGraph) - [Decipher Technology Studios](http://deciphernow.com/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
