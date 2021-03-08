# rollover
[![Go Reference](https://pkg.go.dev/badge/github.com/tomasen/rollover.svg)](https://pkg.go.dev/github.com/tomasen/rollover)

Safely restart current process in golang.

By safely we mean wait parent process 
to exit completely before the child start to allocate resources.
The child process will inherit parent process' arguments and environment variables.

To put every thing simple, we did NOT implement anything "graceful" such as pass network 
resources to child process to avoid service interruption because it's way to complicate 
for programs that listen to multiple network port or socket.

This package is more suitable for agent program that updated its own binary and require 
to restart itself to complete the upgrade.

## Unit test
It's hard to write unit test for a process keep restarting(exit/fork). 

Here is how it can be tested:

1. Build and start the parent process
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

7. Now we can kill the child by run `kill 29544`

