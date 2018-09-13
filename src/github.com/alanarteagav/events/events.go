// ChatEvent constants, to be used in the protocol between the
// client and the server.
package events

type ChatEvent string

// ChatEvent constants, to be used in the protocol between the
// client and the server.
const (
    USERS  ChatEvent = "USERS"
    STATUS  ChatEvent = "STATUS"
    INVITE ChatEvent = "INVITE"
    MESSAGE ChatEvent = "MESSAGE"
    IDENTIFY  ChatEvent = "IDENTIFY"
    JOINROOM ChatEvent = "JOINROOM"
    CREATEROOM ChatEvent = "CREATEROOM"
    DISCONNECT ChatEvent = "DISCONNECT"
    ROOMESSAGE ChatEvent = "ROOMESSAGE"
    PUBLICMESSAGE  ChatEvent = "PUBLICMESSAGE"

    ERROR   ChatEvent = "ERROR"
    UNDEFINED   ChatEvent = "UNDEFINED"
)

// Function that converts strings into ChatEvent constants.
func ToChatEvent(str string) ChatEvent {
    switch str {
    case "USERS":
        return USERS
    case "STATUS":
        return STATUS
    case "INVITE":
        return INVITE
    case "MESSAGE":
        return MESSAGE
    case "IDENTIFY":
        return IDENTIFY
    case "JOINROOM":
        return JOINROOM
    case "CREATEROOM":
        return CREATEROOM
    case "DISCONNECT":
        return DISCONNECT
    case "ROOMESSAGE":
        return ROOMESSAGE
    case "PUBLICMESSAGE":
        return PUBLICMESSAGE

    case "ERROR":
        return ERROR
    default:
        return UNDEFINED
    }
}
