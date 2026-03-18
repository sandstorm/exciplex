<?php

function test() {
    // Wait in a loop so the profiler could sample us multiple times (but should only do so once)
    for ($i = 0; $i < 15; $i++) {
        usleep(1000);
    }
}

// Start profiler but stop immediately so no samples are collected
$profiler = start_profiler(0.05, 0.01);
$profiler->stop();

// if not stopped this will result in multiple stack traces
test();

echo $profiler->getLog();
echo "empty";
