<?php

function inner() {
    // Wait in a loop so the profiler samples us inside this function about 10 times
    for ($i = 0; $i < 15; $i++) {
        usleep(1000);
    }
}

function outer() {
    inner();
}

$profiler = start_profiler(0.01, 0.5);
outer();
$profiler->stop();
$log = $profiler->getLog();
echo $log;
echo "---END---";
