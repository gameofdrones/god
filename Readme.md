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
POST /action/stop
```

```
POST /action/up
```

```
POST /action/down
```

```
POST /action/left
```

```
POST /action/right
```

```
PUT /action/fire
```




