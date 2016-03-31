until kubectl -s http://localhost:8080 get pods &> /dev/null; do
   printf "."
done

kubectl create -s http://localhost:8080 -f broadway-namespace.yaml
