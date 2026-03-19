<?php

require 'profiler_cross_file_function_call_helper.php';

$profiler = start_profiler(0.01, 0.05);
cross_file_work();
$profiler->stop();

echo $profiler->getLog();
