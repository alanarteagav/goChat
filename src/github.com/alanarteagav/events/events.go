package events

type ChatEvent string

const (
    LOG_IN  ChatEvent = "LOG_IN"
    LOG_OUT ChatEvent = "LOG_OUT"
    MESSAGE ChatEvent = "MESSAGE"
    ERROR   ChatEvent = "ERROR"
    UNDEFINED   ChatEvent = "UNDEFINED"
)

func ToChatEvent(str string) ChatEvent {
    switch str {
    case "MESSAGE":
        return MESSAGE
    case "LOG_IN":
        return LOG_IN
    case "LOG_OUT":
        return LOG_OUT
    case "ERROR":
        return ERROR
    default:
        return UNDEFINED
    }
}
