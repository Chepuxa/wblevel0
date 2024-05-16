package subscribers

import (
	"encoding/json"
	"fmt"
	"sync"
	"wbtech/level0/internal/api/cache"
	"wbtech/level0/internal/api/models"
	"wbtech/level0/internal/db/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type OrderSubscriberInterface interface {
	Subscribe() stan.Subscription
	Insert([]byte) error
}

type OrderSubscriber struct {
	Sc              stan.Conn
	Subject         string
	OrderRepository *repositories.OrderRepository
	Validate        *validator.Validate
	Cache           *cache.Cache
	Wg              *sync.WaitGroup
	Logger          *zap.SugaredLogger
}

func NewOrderSubscriber(sc stan.Conn, sb string, or *repositories.OrderRepository, va *validator.Validate, ca *cache.Cache,
	wg *sync.WaitGroup, lg *zap.SugaredLogger) *OrderSubscriber {
	return &OrderSubscriber{
		Sc:              sc,
		Subject:         sb,
		OrderRepository: or,
		Validate:        va,
		Cache:           ca,
		Wg:              wg,
		Logger:          lg,
	}
}

func (s *OrderSubscriber) Subscribe() stan.Subscription {
	sub, _ := s.Sc.Subscribe(s.Subject, func(m *stan.Msg) {
		s.Wg.Add(1)

		s.Logger.Infof("Received message: %s", m.Data)
		err := s.Insert(m.Data)

		if err != nil {
			s.Logger.Error(err.Error())
		}

		s.Wg.Done()
	}, stan.DeliverAllAvailable())

	return sub
}

func (s *OrderSubscriber) Insert(body []byte) error {
	var data models.Data

	unmarshallErr := json.Unmarshal(body, &data)

	if unmarshallErr != nil {
		return unmarshallErr
	}

	validateErr := s.Validate.Struct(&data)

	if validateErr != nil {
		return validateErr
	}

	if _, exists := s.Cache.Get(data.OrderUID); exists {
		return fmt.Errorf("order with order_uid %s already exists", data.OrderUID)
	}

	if _, err := s.OrderRepository.GetById(data.OrderUID); err == nil {
		return fmt.Errorf("order with order_uid %s already exists", data.OrderUID)
	}

	s.Cache.Set(data.OrderUID, data, 0)

	err := s.OrderRepository.Create(&data)

	if err != nil {
		return err
	}

	return nil
}
