
## Yield Chaser

Infra for monitoring Defi assets


### run prometheus
docker run \
    -p 9090:9090 \
    -v /Users/nityananda/projects/ethodyssey/yield-chaser/prometheus/config:/etc/prometheus \
    prom/prometheus
    
https://tomgregory.com/how-and-when-to-use-a-prometheus-gauge/
https://www.robustperception.io/how-does-a-prometheus-gauge-work