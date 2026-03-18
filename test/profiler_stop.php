<?php

$profiler = start_profiler(0.01, 0.02);

// Let it collect some samples
function wait() {
    usleep(50000); // 100ms
}
wait();

$profiler->stop();

// After stop, wait and verify no new samples are added
$logBefore = $profiler->getLog();
usleep(100000); // 100ms
$logAfter = $profiler->getLog();

if ($logBefore === $logAfter) {
    echo "stop_works";
} else {
    echo "FAIL: log changed after stop";
}
