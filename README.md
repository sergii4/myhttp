# Tool

## Build
Run
```
go build
```
there will be executable `myhttp`

## Run
Run just like in your example:
```
./myhttp -parallel 3 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
```
with parameter `parallel` or without parameters
```
./myhttp adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
```

## Notes
The tool itself is a pretty simple command app based on concurrency pattern *worker pool* where parameter _parallel_ is a number of workers doing job parallel(default is 10).
I allowed myself one enhancement http request timeout. 
`yandex.com` for example is blocked in Ukraine and to decrease waiting time I set timeout ot 10 seconds. If that happens or any other http error that result will be eliminated from the output
