package subscribers

import (
	"sync"
	"wbtech/level0/internal/api/cache"
	"wbtech/level0/internal/db/repositories"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type Subscribers struct {
	OrderSubscriber *OrderSubscriber
}

func NewSubscribers(sc stan.Conn, sb string, or *repositories.OrderRepository, va *validator.Validate, ca *cache.Cache,
	wg *sync.WaitGroup, lg *zap.SugaredLogger) *Subscribers {
	return &Subscribers{
		OrderSubscriber: NewOrderSubscriber(sc, sb, or, va, ca, wg, lg),
	}
}

func SubscribeAll(subs *Subscribers) []stan.Subscription {
	subscriptions := []stan.Subscription{}
	subscriptions = append(subscriptions, subs.OrderSubscriber.Subscribe())
	return subscriptions
}

func UnsubscribeAll(s []stan.Subscription) {
	for _, v := range s {
		v.Unsubscribe()
	}
}
