rm -f /tmp/sigmaos-perf/*
./stop.sh
export SIGMADEBUG=""
# export SIGMADEBUG="WATCH_V2;WATCH_PERF;WATCH_TEST" 
# export SIGMAPERF="WATCH_TEST_WORKER_PPROF;WATCH_TEST_WORKER_PPROF_MUTEX;WATCH_TEST_WORKER_PPROF_BLOCK;WATCH_PERF_WORKER_PPROF"
export USE_OLD_WATCH=""
export WATCHPERF_MEASURE_MODE="0"
go test sigmaos/watch -v --start --run "TestWatchPerfSingle"

