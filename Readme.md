# Game of Drones

Raspberry powered rocket launcher

## API

```
GET /position
Accept: application/json

{
  "x": 60,
  "y": 3
}
```

```
PUT /rocket
Content-Type: application/json

{
  "x": 100,
  "y": 6
}
```

```
PUT /actions/stop
```

```
POST /actions/up
```

```
POST /actions/down
```

```
POST /actions/left
```

```
POST /actions/right
```

```
PUT /actions/fire
```


## Dev setup

```
brew install go
brew install libusb
go get github.com/baloo/gousb/usb
```

Mac OS:
```
vi /usr/local/include/libusb-1.0/libusb.h
# Follow  https://gist.github.com/pjvds/4578277
```

