<?php

$timer = exciplex_set_interval(
    function () use (&$timer) {
        echo "callback";
        $timer->stop();
    },
    0.01,
    0.01
);

usleep(150000);
echo "done";
