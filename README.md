# rollover
safe restart current process or updated binary of the current executable file in golang

## Unit test
It's hard to write unit test for a process keep restarting(exit/fork). Here is how it can be tested

```
go build -o ./bin/example ./example
./bin/example -any=arguments -we=want -inherited=to -be=tested
```

Following message is to be expected
```
Current PID: 29539
29539 : this is a parent
29539 : have arguments [./bin/example -any=args]
kill -HUP 29539  -> to rollover(restart).
kill 29539  -> to quit
```

Run `kill -HUP 29539` to initiate rollover(restart).
Following message is to be expected
```
initiating rollover
child started runing, pid: 29544
```

And after 3 second we should see parent cleaned up and exit
```
29539 : safely exit after clean up
```

Also, following message is to be expected:
```
Current PID: 29544
29544 : this is a child
29544 : have arguments [./bin/example -any=args]
kill -HUP 29544  -> to rollover(restart).
kill 29544  -> to quit
```

Now we can kill the child by run `kill 29544`

