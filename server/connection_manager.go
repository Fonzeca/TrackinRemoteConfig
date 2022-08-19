package server

import (
	"fmt"
	"net"
	"time"

	"github.com/Fonzeca/TrackinRemoteConfig/server/decoder"
	"github.com/Fonzeca/TrackinRemoteConfig/server/encoder"
)

var ConnectionPool = make(map[string]([]chan string))

func NewConnection(connection net.Conn) {
	fmt.Println("Connection established: ", connection.RemoteAddr().String())

	//Leemos el primer mensaje
	buffer := make([]byte, 2048)
	mLen, err := connection.Read(buffer)
	//Si el mensaje es vacio, cerró la conexion
	if mLen == 0 {
		fmt.Println("Connection closed by client")
		connection.Close()
		return
	}
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	//verificamos que el mensaje sea de Login
	if buffer[2] != 0x01 {
		fmt.Println("Not Login message")
		connection.Close()
		return
	}

	//Obtenemos el imei
	imei := decoder.DecodeLogin(buffer)

	//Cerramos el pipe viejo si se quedo pegado en la lista
	if ConnectionPool[imei] != nil {
		if ConnectionPool[imei][0] != nil {
			close(ConnectionPool[imei][0])
		}

		if ConnectionPool[imei][1] != nil {
			close(ConnectionPool[imei][1])
		}

	}

	pipeIn := make(chan string)
	pipeOut := make(chan string)

	//Creamos la conexion en el pool
	ConnectionPool[imei][0] = pipeIn
	ConnectionPool[imei][1] = pipeOut

	//Funcion de cierre
	defer func() {
		close(pipeIn)
		close(pipeOut)
		ConnectionPool[imei] = nil
	}()

	//Loop para manter la conexion activa
	for {

		//Esperamos algun mensaje del prompt
		commandToSend := <-pipeIn

		//Encodeamos el mensaje
		dataToSend := encoder.Encode(imei, commandToSend)

		//Mandamos el mensaje al cliente
		_, err := connection.Write(dataToSend)
		if err != nil {
			fmt.Println("Intentelo nuevamente")
			continue
		}

		// Obtenemos el mensaje de respuesta
		buffer := make([]byte, 2048)
		mLen, _ := connection.Read(buffer)
		if mLen == 0 {
			fmt.Println("Cerró la conexion el cliente")
			connection.Close()
			break
		}
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		//Decodificamos el mensaje de respuesta
		imeiResp, content := decoder.Decode(buffer)

		//Verificamos el imei de respuesta
		if imeiResp != imei {
			fmt.Println("Error imei diferentes")
		}

		//Imprimimos el mensaje
		pipeOut <- content

		time.Sleep(time.Second)
	}
}