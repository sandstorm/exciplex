<?php

$profiler = start_profiler(0.01, 0.05);
require 'profiler_cross_file_no_functions_helper.php';
$profiler->stop();

echo $profiler->getLog();