<?php

exciplex_set_timeout(
    function () {
        die ("The allowed time has been exceeded");
    },
    0.1
);

$num = 0;
exciplex_set_interval(
    function () use (&$num) {
        echo ("interrupt $num -- ");
        $num = $num+1;
    },
    0.03,
    0.005
);


echo "Hello from FrankenPHP with continuous profiling!<br /><br />\n";

for($i = 0; $i < 20; $i++) {
    echo ($i < 10 ? " " : "") . "$i, ";
    usleep(10000);
}