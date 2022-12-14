package server

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Fonzeca/TrackinRemoteConfig/server/decoder"
	"github.com/Fonzeca/TrackinRemoteConfig/server/encoder"
)

var ConnectionPool = make(map[string]([]chan string))

func NewConnection(connection net.Conn) {
	fmt.Println("Connection established: ", connection.RemoteAddr().String())

	defer connection.Close()

	imei, err := readLoginMessage(connection)
	if err != nil {
		fmt.Println(err)
		return
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
		delete(ConnectionPool, imei)
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
			pipeOut <- "Intentelo nuevamente"
			continue
		}

		connection.SetReadDeadline(time.Now().Add(time.Second * 20))

		// Obtenemos el mensaje de respuesta
		buffer := make([]byte, 2048)
		mLen, err := connection.Read(buffer)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				pipeOut <- err.Error()
				continue
			}
			pipeOut <- "Error al leer el mensaje"
			continue
		}
		if mLen == 0 {
			pipeOut <- "Cerró la conexion el cliente"
			connection.Close()
			return
		}

		//Decodificamos el mensaje de respuesta
		imeiResp, content, err := decoder.Decode(buffer)
		if err != nil {
			pipeOut <- err.Error()
			return
		}

		//Verificamos el imei de respuesta
		if imeiResp != imei {
			pipeOut <- "Error imei diferentes"
			return
		}

		//Imprimimos el mensaje
		pipeOut <- content

		time.Sleep(time.Second)
	}
}

func readLoginMessage(connection net.Conn) (string, error) {

	//Leemos el primer mensaje
	buffer := make([]byte, 2048)
	mLen, err := connection.Read(buffer)
	//Si el mensaje es vacio, cerró la conexion
	if mLen == 0 {
		return "", errors.New("Connection closed by client")
	}
	if err != nil {
		return "", err
	}

	//verificamos que el mensaje sea de Login
	if buffer[2] != 0x01 {
		return "", errors.New("Not Login message")
	}

	//Obtenemos el imei
	imei, err := decoder.DecodeLogin(buffer)
	if err != nil {
		return "", err
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

	return imei, nil
}
