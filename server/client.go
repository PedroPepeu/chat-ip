package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    // Endereço do servidor (substitua pela IP e porta desejados)
    serverAddress := "127.0.0.1:1234"

    // Estabelecer uma conexão TCP com o servidor
    conn, err := net.Dial("tcp", serverAddress)
    if err != nil {
        fmt.Println("Erro ao conectar:", err)
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println("Conectado ao servidor:", serverAddress)

    // Enviar uma mensagem para o servidor
    message := "Olá, servidor!"
    _, err = conn.Write([]byte(message))
    if err != nil {
        fmt.Println("Erro ao enviar mensagem:", err)
        os.Exit(1)
    }
    fmt.Println("Mensagem enviada:", message)

    // Ler a resposta do servidor (opcional)
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Erro ao ler resposta:", err)
        os.Exit(1)
    }
    fmt.Println("Resposta do servidor:", string(buffer[:n]))
}
