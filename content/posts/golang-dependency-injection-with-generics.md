---
layout: post
title: Dependeny Injection In Golang
summary: Let's observe how to handle dependency injection with samber/lo package.
date: 2024-08-31
tags: [golang]
---

Dependency injection is used to provide the necessary dependencies to the objects from external sources rather than creating the dependencies inside the objects.

For example, instead of a service creating its own dependencies, you inject those dependencies from outside, allowing you to easily swap them for different implementations or mock objects during testing.

For this blogpost we are going to use a dependency injection tool called [samber/do](https://github.com/samber/do). Why we use this package?

- it is simple
- it is lazily loaded.

{{< admonition note  >}}
All of the code here can be accessed from [ocakhasan/dependency_injection_golang](https://github.com/ocakhasan/dependency_injection_golang)
{{< /admonition >}}

## IMPLEMENTATION

We are going to explore the package by creating a simple service. Now image that we have a service which processes the orders for the website. This service depends on 2 other services called `PaymentProcessor` and `NotificationService`.

Let's follow the best practices and define the interfaces for our dependencies.

```go
package main

type PaymentProcessor interface {
	ProcessPayment(amount float64) error
}

type NotificationService interface {
	SendNotification(orderID string) error
}
```

Now let's have some basic implementations for our services.

```go
package main

import "github.com/rs/zerolog"

type EmailNotification struct {
	logger zerolog.Logger
}

func (n *EmailNotification) SendNotification(orderID string) error {
	n.logger.Printf("Sending email notification for order %s", orderID)

	return nil
}

func NewEmailNotification(logger zerolog.Logger) *EmailNotification {
	return &EmailNotification{logger: logger}
}
```

```go
package main

import "github.com/rs/zerolog"

type StripeProcessor struct {
	logger zerolog.Logger
}

func (p *StripeProcessor) ProcessPayment(amount float64) error {
	p.logger.Printf("Processing payment of $%.2f through Stripe", amount)

	return nil
}

func NewStripeProcessor(logger zerolog.Logger) *StripeProcessor {
	return &StripeProcessor{logger: logger}
}
```

As we can see 
- the `StripeProcessor` implements the `PaymentProcessor` inferface.
- the `EmailNotification` implements the `NotificationService` interface.

Now let's create the order service with best practises

```go
package main

type OrderService struct {
    paymentProcessor  PaymentProcessor
    notificationService NotificationService
}

func NewOrderService(paymentProcessor PaymentProcessor, notificationService NotificationService) *OrderService {
    return &OrderService{
        paymentProcessor:  paymentProcessor,
        notificationService: notificationService,
    }
}

func (s *OrderService) ProcessOrder(orderID string, amount float64) error {
    err := s.paymentProcessor.ProcessPayment(amount)
    if err != nil {
        return err
    }

    return s.notificationService.SendNotification(orderID)
}
```

We have the simple implementation ready. Let's do the dependency injection without and with `samber/do` package.

### Straight Forward Injection

It is pretty straightforward, you just need to initiaze the necessary dependency by yourself to provide the services by you.

Now for the small cases, such as this seems fine, but as we can see it lacks some kind of generalization, structure.

```go
package main

import (
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)

	stripeProcessor := NewStripeProcessor(logger)
	emailNotif := NewEmailNotification(logger)

	orderService := NewOrderService(stripeProcessor, emailNotif)

	// Use the service
	err := orderService.ProcessOrder("12345", 100.00)
	if err != nil {
		panic(err)
	}
}
```

### Structured Injection

It seems a little bit complicated but I am going to explain bit by bit.

```go
package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/samber/do"
)

func main() {
	injector := do.New()

	do.Provide(injector, func(i *do.Injector) (zerolog.Logger, error) {
		return zerolog.New(os.Stdout).Level(zerolog.DebugLevel), nil
	})

	// Bind implementations to the interfaces
	do.Provide(injector, func(i *do.Injector) (PaymentProcessor, error) {
		logger := do.MustInvoke[zerolog.Logger](i)

		return NewStripeProcessor(logger), nil
	})

	// Bind implementations to the interfaces
	do.Provide(injector, func(i *do.Injector) (NotificationService, error) {
		logger := do.MustInvoke[zerolog.Logger](i)

		return NewEmailNotification(logger), nil
	})

	do.Provide(injector, func(i *do.Injector) (*OrderService, error) {
		paymentProcessor := do.MustInvoke[PaymentProcessor](i)
		notificationService := do.MustInvoke[NotificationService](i)

		return NewOrderService(paymentProcessor, notificationService), nil
	})

	// Resolve the OrderService and use it
	orderService := do.MustInvoke[*OrderService](injector)
	err := orderService.ProcessOrder("12345", 100.00)
	if err != nil {
		panic(err)
	}
}
```

With this package, to use something, first you need to provide it. There are multiple way to provide something,

- [Provide](https://pkg.go.dev/github.com/samber/do@v1.6.0#Provide):  It is provided anonymously, and it is loaded lazily. You can see the example, we provided the type we want to inject.
- [ProvideNamed](https://pkg.go.dev/github.com/samber/do@v1.6.0#Provide): This is same as the Provide but we define a tag to identify the object we want to provide. It is mostly useful for the case where you need to provide multiple different instances of the same type (such as multiple different loggers). To define them, a tag is needed. For example

```go
injector := do.New()

do.ProvideNamed(injector, "debugLogger", func(i *do.Injector) (zerolog.Logger, error) {
    return zerolog.New(os.Stdout).Level(zerolog.DebugLevel), nil
})

do.ProvideNamed(injector, "errorLogger", func(i *do.Injector) (zerolog.Logger, error) {
    return zerolog.New(os.Stdout).Level(zerolog.ErrorLevel), nil
})

// To get them use do.MustInvokeNamed
debugLogger := do.MustInvokeNamed[zerolog.Logger](injector, "debugLogger")
errorLogger := do.MustInvokeNamed[zerolog.Logger](injector, "errorLogger")
```

These 2 function are loaded lazily which means that if we provided something but we did not call invoke on them, these instances are not created. It is one of the most benefical features of `samber/do` package.

With this package, to access something we provided, we need to invoke them. You can check the examples from code.

With this approach you can even have a helper function, which returns the injector and you can access each of the objects inside. 

```go
package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/samber/do"
)

func Injector() *do.Injector {
	injector := do.New()

	do.Provide(injector, func(i *do.Injector) (zerolog.Logger, error) {
		return zerolog.New(os.Stdout).Level(zerolog.DebugLevel), nil
	})

	// Bind implementations to the interfaces
	do.Provide(injector, func(i *do.Injector) (PaymentProcessor, error) {
		logger := do.MustInvoke[zerolog.Logger](i)

		return NewStripeProcessor(logger), nil
	})

	// Bind implementations to the interfaces
	do.Provide(injector, func(i *do.Injector) (NotificationService, error) {
		logger := do.MustInvoke[zerolog.Logger](i)

		return NewEmailNotification(logger), nil
	})

	do.Provide(injector, func(i *do.Injector) (*OrderService, error) {
		paymentProcessor := do.MustInvoke[PaymentProcessor](i)
		notificationService := do.MustInvoke[NotificationService](i)

		return NewOrderService(paymentProcessor, notificationService), nil
	})

	return injector
}
```

Now our main function just becomes

```go
package main

import (
	"github.com/samber/do"
)

func main() {
	injector := Injector()

	// Resolve the OrderService and use it
	orderService := do.MustInvoke[*OrderService](injector)
	err := orderService.ProcessOrder("12345", 100.00)
	if err != nil {
		panic(err)
	}
}
```

As we can see we cleaned up most of the stuff from main function which is the ideal case.

## Summary

with the samber/do package

- all the dependency bindings are centralized.
- lazy initialization is used.
- automatically handle the resolution of dependencies in correct order.
- scalability is easier. As your project grows, adding new dependencies are easy.