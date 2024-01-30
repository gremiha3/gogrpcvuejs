package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"time"

	pb "gogrpcvuejs/pkg/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ( //входные аргументы программы
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()         //опрашиваем входные аргументы программы
	if flag.NArg() < 2 { //если аргументов недостаточно, то
		log.Fatal("not enough arguments") //пишем в консоль ошибку и останавливаем программу
	}

	x, err := strconv.Atoi(flag.Arg(0)) //считываем первый аргумент и приводим его к типу INT
	if err != nil {                     //если привести к типу INT нельзя, то
		log.Fatal(err) //пишем в консоль ошибку и останавливаем программу
	}
	y, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatal(err) //пишем в консоль ошибку и останавливаем программу
	}
	//создаем клиентское подключение к адресу с портом, используя при этом небезопасный метод
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { //если подключение не удалось, то
		log.Fatalf("did not connect: %v", err) //пишем в консоль ошибку и останавливаем программу
	}
	defer conn.Close()           //По достижению конца функции main, подключение закроется автоматически
	c := pb.NewAdderClient(conn) //создаем объект AdderClient, подключенный через созданное соединение

	//создаем контекст с таймаутом в 1 сек на основе пустого контекста.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() //если код ниже до конца программы не успеет завершиться за 1 секунду, то
	//контекст завершит программу принудительно, иначе
	//выполнится код ниже, и по окончании программы, сработает "defer cancel()" и отменит контекст

	//вызывая функцию Add объекта AdderClient, передаем :
	//
	//            созданный|      структура передаваемых
	//            контекст |      данных серверу
	res, err := c.Add(ctx, &pb.AddRequest{X: int32(x), Y: int32(y)})
	// и возвращаем результат обработки данных от сервера
	if err != nil { //если получили ошибку функции сложения, то
		log.Fatalf("could not greet: %v", err) //пишем в консоль ошибку и останавливаем программу
	}
	log.Println(res.GetResult()) //иначе, печатаем результат

}
