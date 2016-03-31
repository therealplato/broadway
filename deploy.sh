COMMIT=$(git rev-parse HEAD | cut -c1-8)
echo "Version: $COMMIT"
docker-compose run -e CGO_ENABLED=0 test go build -a -installsuffix cgo -ldflags "-s"
docker build -t registry.namely.tech/namely/broadway:$COMMIT -f Dockerfile-build .
docker push registry.namely.tech/namely/broadway:$COMMIT
