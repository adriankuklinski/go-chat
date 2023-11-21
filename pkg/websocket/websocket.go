package websocket

import (
    "crypto/sha1"
    "encoding/base64"
    "log"
    "fmt"
    "net"
    "net/http"
)

const websocketGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

type WebSocketConnection struct {
    Conn net.Conn
}

func computeAcceptKey(clientKey string) string {
    h := sha1.New()
    h.Write([]byte(clientKey + websocketGUID))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
    clientKey := r.Header.Get("Sec-WebSocket-Key")
    if clientKey == "" {
        return nil, fmt.Errorf("no Sec-WebSocket-Key header present")
    }

    log.Println("WebSocket upgrade requested")
    acceptKey := computeAcceptKey(clientKey)

    w.Header().Set("Upgrade", "websocket")
    w.Header().Set("Connection", "Upgrade")
    w.Header().Set("Sec-WebSocket-Accept", acceptKey)
    w.WriteHeader(http.StatusSwitchingProtocols)

    hj, ok := w.(http.Hijacker)
    if !ok {
        return nil, fmt.Errorf("web server does not support hijacking")
    }

    conn, _, err := hj.Hijack()
    if err != nil {
        return nil, fmt.Errorf("hijacking failed: %v", err)
    }

    return conn, nil
}

func readFrame(conn net.Conn) ([]byte, error) {
    header := make([]byte, 2)
    if _, err := conn.Read(header); err != nil {
        return nil, err
    }

    length := int(header[1] & 0x7F)  // Length of payload
    isMasked := (header[1] & 0x80) > 0  // Check if payload is masked

    var maskingKey []byte
    if isMasked {
        maskingKey = make([]byte, 4)
        if _, err := conn.Read(maskingKey); err != nil {
            return nil, err
        }
    }

    payload := make([]byte, length)
    if _, err := conn.Read(payload); err != nil {
        return nil, err
    }

    if isMasked {
        for i := range payload {
            payload[i] ^= maskingKey[i % 4]
        }
    }

    return payload, nil
}

func writeFrame(conn net.Conn, data []byte) error {
    fmt.Println("before frame data:")
    fmt.Println(string(data))
    frame := []byte{0x81, byte(len(data))}
    frame = append(frame, data...)


    fmt.Printf("Frame header: [%x %x]\n", frame[0], frame[1])
    fmt.Printf("Frame payload: %s\n", string(frame[2:]))
    fmt.Println("after frame data: ");
    fmt.Println(string(frame));
    _, err := conn.Write(frame)
    return err
}

func (wsc *WebSocketConnection) ReadMessage() ([]byte, error) {
    return readFrame(wsc.Conn)
}

func (wsc *WebSocketConnection) WriteMessage(message []byte) error {
    return writeFrame(wsc.Conn, message)
}
