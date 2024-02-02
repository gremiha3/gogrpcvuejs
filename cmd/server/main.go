package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "gogrpcvuejs/pkg/api" //алиас на папку со сгенерированным protobuf файлом

	"google.golang.org/grpc"
)

var ( //опрашиваем параметры командной строки
	// -port 8080
	// если ничего не указывать, то по умолчанию будет 50051
	// если запустить с параметром -h/--help , то выведет справку
	port = flag.Int("port", 50051, "The server port")
)

// структура server используется для реализации интерфейса pb.AdderServer.
type server struct {
	pb.UnimplementedAdderServer
}

// Добавляем реализацию интерфейса pb.AdderServer через метод Add структуры server
//
//	type AdderServer interface {
//		Add(context.Context, *AddRequest) (*AddResponse, error)
//		mustEmbedUnimplementedAdderServer()
//	}
//
//	структура         контекст			  входное сообщение	  выходное сообщение
func (s *server) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	log.Printf("Received: %v, %v", in.GetX(), in.GetY())
	return &pb.AddResponse{Result: in.GetX() + in.GetY()}, nil
}

func main() {
	flag.Parse()                                             //обрабатываем входные аргументы программы
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port)) //запускаем прослушивание сети
	if err != nil {                                          //если запуск сети произошел с ошибкой, то
		log.Fatalf("failed to listen: %v", err) //выводим лог в терминал и закрываем программу
	}
	s := grpc.NewServer()                //создаем gRPC-сервер, который еще не имеет зарегистрированных служб и не начал принимать запросы.
	pb.RegisterAdderServer(s, &server{}) //регистрируем службу AdderServer на сервере "s" с указанием
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
