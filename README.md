# exciplex

A sampling profiler for PHP running inside [FrankenPHP](https://frankenphp.dev).

## Inspiration

exciplex is directly inspired by [Excimer](https://www.mediawiki.org/wiki/Excimer), Wikimedia's sampling profiler for PHP. exciplex targets FrankenPHP specifically — see [How it works](#how-it-works) for why a separate implementation is needed.

## Usage

The typical pattern is continuous production profiling: start a profiler with a random initial delay so that samples are spread across many requests, then write the result on shutdown.

```php
$initialDelay = rand(min: 0, max: 60000) / 1000.0;
$profiler = start_profiler(initialDelay: $initialDelay, interval: 60);

register_shutdown_function(function () use ($profiler) {
    $profiler->stop();
    $data = $profiler->getLog();

    // return early if no samples were collected
    if ($data === '') {
        return;
    }

    // write log to file or send it to an API endpoint
    file_put_contents('/tmp/exciplex-traces.log', $data, FILE_APPEND);
});
```

With a random initial delay of up to 60 seconds and a 60-second interval, most requests finish before the first sample is taken and contribute nothing. Across many requests the random delay distributes samples uniformly over time. The collected logs can be aggregated and rendered as a flame graph with [flamegraph.pl](https://github.com/brendangregg/FlameGraph) or any tool that accepts the collapsed-stack format.

`getLog()` returns data in the collapsed-stack format:

```
/app/index.php;processRequest;executeQuery 42
/app/index.php;processRequest 8
```

Each line is a unique call stack observed during profiling, with frames separated by semicolons (outermost first) and the sample count at the end.

You can also use `exciplex_set_timeout` and `exciplex_set_interval` directly for general-purpose timer callbacks:

```php
// throw after 500 ms if the function hasn't returned
$timer = exciplex_set_timeout(function (): void {
    throw new Exception('timed out after 500ms');
}, 0.5);
potentiallyLongRunningFunction();
$timer->stop();
```

## Installation

> exciplex is not yet published as a standalone package. Installation instructions will be added once a release mechanism is in place.

## How it works

### Interrupt mechanism

PHP's virtual machine checks an interrupt flag (`EG(vm_interrupt)`) between opcode dispatches. When the flag is set, the VM calls a registered interrupt handler before continuing — allowing code to run safely on the PHP thread at a predictable point. exciplex uses this mechanism to deliver timer callbacks and capture stack traces. The idea comes directly from Excimer.

### Architecture

exciplex is built in three layers: Go handles timers and scheduling via goroutines, C provides the Zend extension glue and stack trace capture, and PHP exposes the public API.

Excimer triggers sampling using POSIX signals, which doesn't work in FrankenPHP. exciplex is compiled directly into FrankenPHP as a Go module, which means goroutines can set the Zend interrupt flag at the correct time. This works, but it ties exciplex to FrankenPHP: the same approach cannot be used in a standard PHP process.
