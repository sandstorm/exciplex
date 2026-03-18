<?php

class TestClass {
    function method()
    {
        middle_call();
    }

    static function staticMethod()
    {
        deep_call();
    }
}

function deep_call() {
    for($i = 0; $i < 15; $i++) {
        usleep(10000);
    }
}

function middle_call() {
    deep_call();
}

function top_call() {
    middle_call();
}

$instance = new TestClass();

$profiler = start_profiler(0.01, 0.05);
top_call();
$instance->method();
TestClass::staticMethod();
$profiler->stop();

echo $profiler->getLog();