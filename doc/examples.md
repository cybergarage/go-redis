# Examples

## go-redis-server

The `go-redis` is disributed with the `go-redis-server` which is a example Redis compatible server implemetation using the `go-redis`.

- [go-redis-server](examples/go-redis-server)

This example is implemented to return the same responses and errors as the official Redis server, but its performance characteristics may vary due to the simple and naive implementation.

### Installation

To set up this example server, run the following command:

```
make install
```
The example server is compiled and installed to your GO_PATH/bin directory. The example server supports major Redis commands, so you can run YCSB (Yahoo! Cloud Serving Benchmark) to the example server as the following. 

```
YCSB Client 0.17.0

Loading workload...
Starting test.
DBWrapper: report latency for each error is false and specific error codes to track for latency are: []
2022-09-03 17:07:47:680 0 sec: 0 operations; est completion in 0 second 
2022-09-03 17:07:47:879 0 sec: 1000 operations; 4255.32 current ops/sec; [READ: Count=479, Max=13207, Min=83, Avg=146.54, 90=153, 99=298, 99.9=13207, 99.99=13207] [CLEANUP: Count=1, Max=688, Min=688, Avg=688, 90=688, 99=688, 99.9=688, 99.99=688] [UPDATE: Count=521, Max=934, Min=174, Avg=208.47, 90=236, 99=319, 99.9=512, 99.99=934] 
[OVERALL], RunTime(ms), 235
[OVERALL], Throughput(ops/sec), 4255.31914893617
[TOTAL_GCS_G1_Young_Generation], Count, 0
[TOTAL_GC_TIME_G1_Young_Generation], Time(ms), 0
[TOTAL_GC_TIME_%_G1_Young_Generation], Time(%), 0.0
[TOTAL_GCS_G1_Old_Generation], Count, 0
[TOTAL_GC_TIME_G1_Old_Generation], Time(ms), 0
[TOTAL_GC_TIME_%_G1_Old_Generation], Time(%), 0.0
[TOTAL_GCs], Count, 0
[TOTAL_GC_TIME], Time(ms), 0
[TOTAL_GC_TIME_%], Time(%), 0.0
[READ], Operations, 479
[READ], AverageLatency(us), 146.53653444676408
[READ], MinLatency(us), 83
[READ], MaxLatency(us), 13207
[READ], 95thPercentileLatency(us), 173
[READ], 99thPercentileLatency(us), 298
[READ], Return=OK, 479
[CLEANUP], Operations, 1
[CLEANUP], AverageLatency(us), 688.0
[CLEANUP], MinLatency(us), 688
[CLEANUP], MaxLatency(us), 688
[CLEANUP], 95thPercentileLatency(us), 688
[CLEANUP], 99thPercentileLatency(us), 688
[UPDATE], Operations, 521
[UPDATE], AverageLatency(us), 208.4721689059501
[UPDATE], MinLatency(us), 174
[UPDATE], MaxLatency(us), 934
[UPDATE], 95thPercentileLatency(us), 252
[UPDATE], 99thPercentileLatency(us), 319
[UPDATE], Return=OK, 521
````
