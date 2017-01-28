# go-sbus

### Pub/Sub
```golang
import "github.com/kulikov/go-sbus"

sbus.Sub("my-command-subject", func(msg sbus.Message, log *logrus.Entry) error {
    log.Infof("Receive: %v", msg)
    return nil
})

sbus.Pub("my-command-subject", map[string]string{"do": "any"})
```

### Request/Reply
```golang
import "github.com/kulikov/go-sbus"

sbus.Sub("get-user", func(msg sbus.Message, log *logrus.Entry) error {
    sbus.Reply(map[string]string{"name": "dima", "userId": msg.Data["userId"]})
    return nil
})

sbus.Request("get-user", map[string]string{"userId":"12345"}, func(msg sbus.Message, log *logrus.Entry) error {
    log.Infof("Result: %v", msg)
    return nil
})
```
