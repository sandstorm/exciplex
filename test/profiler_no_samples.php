<?php

// Start profiler with long delay so no samples are collected
$profiler = start_profiler(10.0, 10.0);
// Stop immediately
$profiler->stop();
echo $profiler->getLog();
echo "empty";
