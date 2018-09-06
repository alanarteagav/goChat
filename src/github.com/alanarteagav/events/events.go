// ChatEvent constants, to be used in the protocol between the
// client and the server.
package events

type ChatEvent string

// ChatEvent constants, to be used in the protocol between the
// client and the server.
const (
    LOG_IN  ChatEvent = "LOG_IN"
    LOG_OUT ChatEvent = "LOG_OUT"
    MESSAGE ChatEvent = "MESSAGE"
    CREATE_CHATROOM ChatEvent = "CREATE_CHATROOM"
    JOIN_CHATROOM ChatEvent = "JOIN_CHATROOM"
    LEAVE_CHATROOM ChatEvent = "LEAVE_CHATROOM"
    MESSAGE_CHATROOM ChatEvent = "MESSAGE_CHATROOM"
    ERROR   ChatEvent = "ERROR"
    UNDEFINED   ChatEvent = "UNDEFINED"
)

// Function that converts strings into ChatEvent constants.
func ToChatEvent(str string) ChatEvent {
    switch str {
    case "MESSAGE":
        return MESSAGE
    case "LOG_IN":
        return LOG_IN
    case "LOG_OUT":
        return LOG_OUT
    case "CREATE_CHATROOM":
        return CREATE_CHATROOM
    case "JOIN_CHATROOM":
        return JOIN_CHATROOM
    case "LEAVE_CHATROOM":
        return LEAVE_CHATROOM
    case "MESSAGE_CHATROOM":
        return MESSAGE_CHATROOM
    case "ERROR":
        return ERROR
    default:
        return UNDEFINED
    }
}
