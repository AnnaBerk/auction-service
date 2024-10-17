package app

import (
	"auction/internal/domain"
	"auction/internal/infrastructure/notify"
	"auction/internal/infrastructure/payment"
	"auction/internal/infrastructure/repo"
	"auction/internal/interfaces/rpc"
	v1 "auction/internal/interfaces/rpc/pb"
	"context"
	"embed"
	"fmt"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	Cfg     Config
	Db      *pg.DB
	Log     *log.Logger
	Auction domain.AuctionService
	workers []Worker
}

type Worker interface {
	Start()
	Stop()
}

func NewApp(cfg Config, db *pg.DB) (*App, error) {
	log := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	lotRepo := repo.NewLotRepository(db)
	userRepo := repo.NewUserRepository(db)
	auctionRepo := repo.NewAuctionRepository(db)
	bidRepo := repo.NewBidRepository(db)

	notifyService := notify.NewNotifyService(userRepo)
	payment := payment.NewBalanceService()
	auctionService := NewAuctionService(lotRepo, userRepo, auctionRepo, bidRepo, notifyService, payment)

	auctionWorker := NewAuctionWorker(auctionService, log)

	return &App{
		Cfg:     cfg,
		Db:      db,
		Log:     log,
		Auction: auctionService,
		workers: []Worker{
			auctionWorker,
		},
	}, nil
}

func (a *App) awaitShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	for _, worker := range a.workers {
		worker.Stop()
	}

	if err := a.Db.Close(); err != nil {
		a.Log.Printf("failed to close database: %v", err)
	}

}

func InitDB(cfg Config) (*pg.DB, error) {
	opt, err := pg.ParseURL(cfg.ConnectionString)
	if err != nil {
		panic(err)
	}
	bdc := pg.Connect(opt)

	return bdc, nil
}

func (a *App) Run() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.startGRPCServer(); err != nil {
			a.Log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.startRESTGateway(); err != nil {
			a.Log.Fatalf("failed to start REST gateway: %v", err)
		}
	}()

	for _, worker := range a.workers {
		wg.Add(1)
		go func(w Worker) {
			defer wg.Done()
			w.Start()
		}(worker)
	}

	a.awaitShutdown()

	for _, worker := range a.workers {
		worker.Stop()
	}

	wg.Wait()
}

func (a *App) startGRPCServer() error {
	lis, err := net.Listen("tcp", ":"+a.Cfg.GRPCPort)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	v1.RegisterAuctionServiceServer(grpcServer, rpc.NewAuctionHandler(a.Auction))

	a.Log.Printf("Starting gRPC server on :%s", a.Cfg.GRPCPort)
	return grpcServer.Serve(lis)
}

func (a *App) startRESTGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	endpoint := "localhost:" + a.Cfg.GRPCPort
	err := v1.RegisterAuctionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}

	a.Log.Printf("Starting HTTP/REST gateway on :%s", a.Cfg.HTTPServer.Port)
	return http.ListenAndServe(":"+a.Cfg.HTTPServer.Port, mux)
}

//go:embed migrations/*.sql
var MigrationFS embed.FS

func RunMigrations(db *pg.DB) error {
	_, _, err := migrations.Run(db, "init")
	if err != nil {
		if err.Error() != "migration table exists" {
			return fmt.Errorf("failed to initialize migration table: %w", err)
		}
	}

	collection := migrations.NewCollection()

	fs := http.FS(MigrationFS)

	entries, err := MigrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read embedded migrations: %w", err)
	}
	fmt.Println("Discovered migration files:")
	for _, entry := range entries {
		fmt.Println(" -", entry.Name())
	}

	collection.DiscoverSQLMigrationsFromFilesystem(fs, "migrations")

	fmt.Println("Migrations to apply:")
	for _, m := range collection.Migrations() {
		fmt.Println(" -", m)
	}

	if _, _, err := collection.Run(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully.")
	return nil
}
