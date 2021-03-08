# rollover
[![Go Reference](https://pkg.go.dev/badge/github.com/tomasen/rollover.svg)](https://pkg.go.dev/github.com/tomasen/rollover)

Safely restart current process in golang.

By safely we mean it waits parent process 
to exit completely before the child process starts to allocate resources.
The child process will inherit parent process' arguments and environment variables.

To put everything simple, we did NOT implement anything "graceful" such as pass network 
resources to child process to avoid service interruption because it's way too complicate 
for programs that is listening to multiple network ports or sockets.

This package is more suitable for an agent program that updates its own binary and requires 
to restart itself to complete the upgrade.

## Unit test
It's hard to write unit tests for a process keep restarting(exit/fork). 

Here is how it can be tested:

1. Build and start a parent process
```
go build -o ./bin/example ./example
./bin/example -any=arguments -we=want -inherited=to -be=tested
```

2. Following message is to be expected
```
Current PID: 29539
29539 : this is a parent
29539 : have arguments [./bin/example -any=args]
kill -HUP 29539  -> to rollover(restart).
kill 29539  -> to quit
```

3. Run `kill -HUP 29539` to initiate rollover(restart).
   
4. Following message is to be expected
```
initiating rollover
child started runing, pid: 29544
```

5. And after 3 second we should see parent cleaned up and exit
```
29539 : safely exit after clean up
```

6. Also, following message is to be expected:
```
Current PID: 29544
29544 : this is a child
29544 : have arguments [./bin/example -any=args]
kill -HUP 29544  -> to rollover(restart).
kill 29544  -> to quit
```

7. Loop between step 3 to 6 to see if it works for multiple rollovers(restart).

8. Now we can kill the child by run `kill 29544`

## Few considerations

Q. Why kill the parent from child instead of quitting the parent when rollover.Restart()?

A. In case the child is failed to spawn.
