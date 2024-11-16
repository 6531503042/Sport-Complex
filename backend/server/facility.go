package server

import (
	"encoding/json"
	"log"
	"main/modules/booking"
	"main/modules/facility/handler"
	facilityPb "main/modules/facility/proto"
	"main/modules/facility/repository"
	"main/modules/facility/usecase"
	"main/pkg/grpc"
	"main/pkg/queue"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func (s *server) facilityService() {
	// Check if port is already in use
	ln, err := net.Listen("tcp", s.cfg.App.Url)
	if err != nil {
		log.Fatalf("Port %s is already in use", s.cfg.App.Url)
		return
	}
	ln.Close()

	repo := repository.NewFacilityRepository(s.db)
	facilityUsecase := usecase.NewFacilityUsecase(repo)
	httpHandler := handler.NewFacilityHttpHandler(s.cfg, facilityUsecase)
	grpcHandler := handler.NewFacilityGrpcHandler(facilityUsecase)

	// Initialize Kafka consumer
	consumer, err := queue.ConnectConsumer(
		[]string{s.cfg.Kafka.Url},
		s.cfg.Kafka.ApiKey,
		s.cfg.Kafka.Secret,
	)
	if err != nil {
		log.Printf("Warning: Failed to connect to Kafka: %v", err)
	} else {
		// Start consuming booking events in a goroutine
		go func() {
			defer consumer.Close()
			
			partitionConsumer, err := consumer.ConsumePartition("booking.created", 0, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Error creating partition consumer: %v", err)
				return
			}
			defer partitionConsumer.Close()

			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

			for {
				select {
				case msg := <-partitionConsumer.Messages():
					var bookingEvent booking.Booking
					if err := json.Unmarshal(msg.Value, &bookingEvent); err != nil {
						log.Printf("Error unmarshalling booking event: %v", err)
						continue
					}
					log.Printf("Processing booking event: %+v", bookingEvent)
				case err := <-partitionConsumer.Errors():
					log.Printf("Error from consumer: %v", err)
				case <-sigchan:
					log.Println("Shutting down facility Kafka consumer...")
					return
				}
			}
		}()
	}

	// HTTP routes
	facility := s.app.Group("/facility_v1")
	facility.GET("/facilities", httpHandler.FindManyFacility)
	facility.GET("/facility/:facility_id", httpHandler.FindOneFacility)
	facility.POST("/facility", httpHandler.CreateFacility)

	// Slot Routes
	facilitySlot := facility.Group("/:facilityName/slot_v1")
	facilitySlot.POST("/slots", httpHandler.InsertSlot)
	facilitySlot.GET("/slots/:slot_id", httpHandler.FindOneSlot)
	facilitySlot.GET("/slots", httpHandler.FindAllSlots)

	// Badminton Routes
	badminton := facility.Group("/badminton_v1")
	badminton.POST("/court", httpHandler.InsertBadCourt)
	badminton.POST("/slot", httpHandler.InsertBadmintonSlot)
	badminton.GET("/slots", httpHandler.FindBadmintonSlot)
	badminton.GET("/courts", httpHandler.FindCourt)

	// Start gRPC server
	go func() {
		grpcServer, lis := grpc.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.FacilityUrl)
		facilityPb.RegisterFacilityServiceServer(grpcServer, grpcHandler)
		log.Printf("Facility gRPC server listening on %s", s.cfg.Grpc.FacilityUrl)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Println("Facility service initialized")
}

func checkPortAvailable(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	ln.Close()
	return nil
}
