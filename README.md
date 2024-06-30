# Rate limiter

Sistema CLI em Go para realizar testes de carga em um serviço web.

### Instalação
```sh
download 
ou
git clone git@github.com/nagahshi/pos_go_stress_test
cd pos_go_stress_test
```

### Parâmetros
```sh
-h --help:        Ajuda
-u --url:         URL do serviço a ser testado.
-r --requests:    Número total de requests.
-c --concurrency: Número de chamadas simultâneas.
```

### Como usar
local
```sh
go run main.go
# ex. go run main.go --url=http://example.com --requests=1000 --concurrency=10
```

via docker
```sh
docker build -t stress .
docker run stress --url=http://example.com --requests=1000 --concurrency=10
```