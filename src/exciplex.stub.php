<?php

/** @generate-class-entries */

final class ExciplexTimer {
    public function stop(): void {}
}

function exciplex_set_timeout(callable $callable, float $interval): ExciplexTimer {}

function exciplex_set_interval(callable $callable, float $initialDelay, float $interval): ExciplexTimer {}

final class ExciplexProfiler {
    public function stop(): void {}
    public function getLog(): string {}
}

function start_profiler(float $initialDelay, float $interval): ExciplexProfiler {}
