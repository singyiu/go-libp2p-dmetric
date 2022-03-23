# Example - dCounter

## Prerequisite
* go 1.18+
* prometheus
* grafana

## Instructions
```
git clone https://github.com/singyiu/go-libp2p-dmetric.git
cd go-libp2p-dmetric/examples/dcounter
```

start libp2p node that would publish a counter metric
```
go run main.go
```

start libp2p node that would collect metrics and publish to prometheus (localhost:2112)
```
go run main.go -role=collector
```

edit prometheus.yml and add the following to the scrap_configs:
```
   - job_name: 'go-libp2p-dmetric-examples-dcounter'
     static_configs:
       - targets: ['localhost:2112']
```

run Prometheus and visualize the metrics in Grafana
![grafana_screenshot_01](grafana_screenshot_01.png?raw=true "Grafana screenshot")