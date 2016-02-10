# alertd

Alertd is a simple lightweight alerting daemon for InfluxDB. It provides a simple threshold based alerting service and can send alerts via Email and PagerDuty. Metrics are sent to alertd via a REST API and are checked against alerting rules which have been previously defined via alertD, a simple DSL. The motivation for alertd was the lack of a decent alerting tool for Influx. This is now no longer the case as now Kapacitor is available as well as Bosun providing support for InfluxDB.

### alertql
Alert rule are created via alertql, a simple SQL inspired domain specific language. 

#### Creating Alert Rules
Alert rules have a name which identifies them, a metric which they are associated, a condition which is the central part of the rule and text, which describes the rule. The query to create a new alert rule takes the following form:
```
ALERT <alert name> IF <metric name> <operator> <threshold value> TEXT <description of alert>
```

The alert name is simply an identifier for the alert. The metric name is the metric which the alert corresponds to. The operator and threshold specify when the alert is triggered. Here are several more concrete examples:

```
ALERT cpuOnFireAlert IF superImportantServer.cpuUsage > 100 TEXT "Critical production server is heavily loaded"
ALERT noplayers IF tq.currentPlayers == 0 TEXT "something has gone badly wrong"
```

#### Creating Scheduled Database Queries
Once a rule has been created, they can be evaluated against data points. Data points can be passed to alertd via the API or alternatively they be can be pulled from InfluxDB
scheduled queries. A scheduled query is an InfluxDB query which alertd executes at regular intervals and evaluates against stored alerting rules. A scheduled InfluxDB query must
return a single value and single point, multi value or multi point queries are not accepted. E.g. the query: 

````
select cpufree, cpuused, cpuidle from host1.cpu
select cpuidle from host1.cpu
select * from host1.cpu
````
are not a valid scheduled query, where as:
````
select max(cpuidle) from host1.cpu
select last(cpufree) from host1.cpu
````
are acceptble scheduled queres. InfluxDB's aggregate functions are useful to return single points. So long as the query returns a single point and single value, the full range of InfluxQL functionality can be used. It is advisable to test your query via the InfluxDB web interfce, Chronograf or Grafana before scheduling it in alertd.

A scheduled query is pased to alertd by encapsulating it in a SCHEDULE statement. A schedule statement takes the form:

````
SCHEDULE <metric name> INFLUXDB <influx query> ON <influx database>
````
To provide specific examples:
````
SCHEDULE cpuOnFire INFLUXDB "select max(value) from host1.cpu where time > now() - 1h" ON public
SCHEDULE noplayers INFLUXDB "select last(value) from tq.currentPlayers" ON production
````
Note that the metric name in the schedule query should correspond to the metric name in the alert rule. Once a query is scheduled, the resulting metric will be checked against all
corresponding alert rules.
