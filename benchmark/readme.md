```shell
go build -o finonacci.exe .\benchmark\


.\finonacci.exe -engine=eval
# engine=eval, result=9227465, duration=10.4356113s

.\finonacci.exe -engine=vm
# engine=vm, result=9227465, duration=2.8180383s
```
