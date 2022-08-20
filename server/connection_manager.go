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

	defer connection.Close()

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
	imei, err := decoder.DecodeLogin(buffer)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

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
	ConnectionPool[imei] = []chan string{pipeIn, pipeOut}

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

		connection.SetReadDeadline(time.Now().Add(time.Second * 10))

		// Obtenemos el mensaje de respuesta
		buffer := make([]byte, 2048)
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		if mLen == 0 {
			fmt.Println("Cerró la conexion el cliente")
			connection.Close()
			break
		}

		//Decodificamos el mensaje de respuesta
		imeiResp, content, err := decoder.Decode(buffer)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		//Verificamos el imei de respuesta
		if imeiResp != imei {
			fmt.Println("Error imei diferentes")
			return
		}

		//Imprimimos el mensaje
		pipeOut <- content

		time.Sleep(time.Second)
	}
}
