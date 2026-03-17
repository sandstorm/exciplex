<?php

exciplex_set_timeout(
    function () {
        die("The allowed time has been exceeded");
    },
    0.1
);


echo "Hello from FrankenPHP with continuous profiling!<br /><br />\n";

for($i = 0; $i < 20; $i++) {
    echo ($i < 10 ? " " : "") . "$i, ";
    usleep(10000);
}