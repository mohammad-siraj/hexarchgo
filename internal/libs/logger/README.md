## Summary
This code defines a function called `NewLogger` that creates a new logger based on the provided configuration. The logger uses the Zap logging library and supports log rotation using the Lumberjack library.

## Example Usage
```go
config := NewConfig()
logger := NewLogger(config)
logger.Info("This is an info log")
logger.Error("This is an error log")
```

## Code Analysis
### Inputs
- `config` (Iconfigs): The configuration object that specifies the logger settings.
___
### Flow
1. The function first creates an encoder configuration for the logger. If the configuration is for production, it uses the production encoder configuration, otherwise it uses the development encoder configuration.
2. The function sets the time key and time encoder for the encoder configuration.
3. The function calls the `genarateLogRotater` function to create a log rotater based on the provided configuration.
4. The function creates a new Zap core using the JSON encoder, the log rotater, and the atomic log level.
5. The function creates a new logger using the created core and returns it.
___
### Outputs
- `ILogger`: The created logger object that can be used to log messages.
___
