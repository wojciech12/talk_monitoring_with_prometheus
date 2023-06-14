# Monitoring Golang App with Prometheus

## How-to

1. Start the application, *Prometheus*, *Alertmanager*, and *Grafana* with docker-compose:

   ```bash
   make start
   ```

2. Access the services:

   - 8080 - our application
   - 9090 - Prometheus webgui
   - 9093 - AlertManager
   - 3000 - Grafana.

3. Single calls, see targets with prefix *srv_* in [Makefile](Makefile).

4. Generate traffic (open Grafana dashboard to see the metrics):

   ```bash
   make srv_wrk_random
   ```

   if you have issues in making *wrk* work, you will find the [locust](https://locust.io/) configuration [here](../demo_py).

5. Questions:

   - What can we learn from the graphs?
   - Can we say sth about out random calls?
   - Naming? Is it good?

## Changing the default configuration

Docker-compose mounts all configuration from the git repo. You can change it locally on your laptop.

1. To reload prometheus configuration after changes:

   ```bash
   make prometheus_reload_config
   ```

2. To reload grafana configuration, restart the grafana docker:

   ```bash
   docker restart talkgolangpromred_grafana_1
   ```

## Example of Prometheus Queries

- simple:

  ```promql
  order_mgmt_duration_seconds_sum{status_code='200'}
  ```

  ```promql
  order_mgmt_duration_seconds_sum{job=~".*"}
  or
  order_mgmt_database_duration_seconds_sum{job=~".*"}
  or
  order_mgmt_audit_duration_seconds_sum{job=~".*"}
  ```

- based on weave blog (https://www.weave.works/blog/of-metrics-and-middleware/):

  - QPS:

    ```promql
    sum(irate(order_mgmt_duration_seconds_count{job=~".*"}[1m])) by (status_code)
    ```

  - will give you the rate of requests returning 500s:

    ```promql
    sum(irate(order_mgmt_duration_seconds_count{job=~".*", status_code=~"5.."}[1m]))
    ```

  - by status_code:

    ```promql
    sum(irate(order_mgmt_duration_seconds_count{job=~".*"}[1m])) by (status_code)
    ```

  - 500s:

    ```promql
    sum(irate(order_mgmt_duration_seconds_count{job=~".*", status_code=~"5.."}[1m]))
    ```

  - will give you the 5-min moving 99th percentile request latency:

    ```promql
    histogram_quantile(0.99, sum(rate(order_mgmt_duration_seconds_count{job=~".*",ws="false"}[5m])) by (le))
    ```

### Related Work

- https://prometheus.io/docs/prometheus/latest/querying/functions/
- https://www.robustperception.io/combining-alert-conditions/
