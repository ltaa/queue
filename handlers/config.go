package handlers

import (
	"database/sql"
	"log"
	"github.com/streadway/amqp"
	"log/syslog"
	"runtime"
)

var db *sql.DB

const (
	defaultAmqpChannel = "queueTest"
	defaultAmqpUrl = "amqp://guest:guest@localhost:5672/"
	defaultDbSource = "postgres://postgres:postgres@localhost:5432/queue?sslmode=disable"
)

type Config struct {
	AmqpUrl string
	AmqpChannel string

	DbUrl string

	db *sql.DB
	msgs <- chan amqp.Delivery
	amqpCon *amqp.Connection
	amqpChan *amqp.Channel

	logger *log.Logger
	GoroutinesNum int

}

var cfg Config

func NewConfig() *Config {
	conf := Config{
		AmqpUrl: defaultAmqpUrl,
		AmqpChannel: defaultAmqpChannel,
		DbUrl: defaultDbSource,
		GoroutinesNum: runtime.NumCPU(),

	}

	return &conf
}


func (c *Config) Init() {

	if c.AmqpChannel == "" || c.AmqpUrl == "" || c.DbUrl == ""{
		return
	}

	if cfg.amqpCon != nil {
		cfg.amqpCon.Close()
		cfg.amqpChan.Close()

		///!!! close amqp messages chan
		//close(amqp.)}cfg.msgs)

	}

	if cfg.db != nil {
		cfg.db.Close()
	}

	cfg = *c

	logger, err := syslog.NewLogger(syslog.LOG_DAEMON | syslog.LOG_EMERG, -1)
	if err != nil {
		panic(err)
	}
	cfg.logger = logger

	if cfg.GoroutinesNum <= 0 {
		cfg.GoroutinesNum = 1
	}


	err = amqpInit(cfg.AmqpUrl, cfg.AmqpChannel)
	if err != nil {
		cfg.logger.Print(err)
		panic(err)
	}

	err = dbInit(cfg.DbUrl)
	if err != nil {
		log.Print(err)
		panic(err)
	}


}

func dbInit(dataSouceName string) error {
	db, err := sql.Open("postgres", dataSouceName)

	if err != nil {
		return err
	}

	cfg.db = db

	return nil
}

func amqpInit(url string, amqpChannel string) error {
	con, err := amqp.Dial(url)

	if err != nil {
		return err
	}

	cfg.amqpCon = con
	//defer con.Close()

	ch, err := con.Channel()

	if err != nil {
		return err
	}
	cfg.amqpChan = ch
	//defer ch.Close()


	q, err := ch.QueueDeclare(amqpChannel,false,false, false, false, nil)

	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	cfg.msgs = msgs

	return nil
}