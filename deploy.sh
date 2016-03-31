rm broadway
docker-compose run -e CGO_ENABLED=0 test go build -a -installsuffix cgo -ldflags "-s"
docker build -t registry.namely.tech/namely/broadway:`echo $TRAVIS_COMMIT | cut -c1-8` -f Dockerfile-build .
docker push registry.namely.tech/namely/broadway:`echo $TRAVIS_COMMIT | cut -c1-8`
