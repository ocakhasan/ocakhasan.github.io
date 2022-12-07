---
layout: post
title: Go Dilinde Concurrency Üzerinde Çalışma
summary: Kod parçasındaki hataları bulup düzelteceğiz
date: 2022-09-14
tags: [golang, concurrency, refactor]
---

# Özet

Bu yazıda basit bir kod parçasındaki bütün hataları bulup refactor edeceğim. Bunu yaparken de Go dilindeki temel unsurları açıklayarak yapacağım.

***Bu yazı [Concurrency Made Easy](https://www.youtube.com/watch?v=DqHb5KBe7qI&t=643s) videosundan ağır şekilde esinlenmiştir.***

Go dilinde `concurreny` baya öne çıkan bir unsur ancak doğru kullanmayı bilmek daha da önemli. Kendim de bu konuda mükemmel sayılmam ancak hala öğreniyorum.

## Elimizdeki Fonksiyon

Elimizdeki fonksiyon sadece bir parametre `websites` alıyor. Bu websiteler üzerinde gezinirken `handle` diye error döndüren bir fonksiyon alıyor ve `handle` fonksiyonu herhangi bir error döndürdüğü anda ise bu erroru döndürmek istiyor.

```go
func handleWebsites(websites []string) error {
	errChan := make(chan error, 1)
	semaphores := make(chan struct{}, 5) // aynı anda 5 iş çalıştır
	var wg sync.WaitGroup
	wg.Add(len(websites))
	for _, website := range websites {
		semaphores <- struct{}{}    // semaphore acquire et
		go func() {
			defer func() {
				wg.Done()
				<-semaphores
			}()
			if err := handle(website); err != nil {
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(semaphores)
	close(errChan)
	return <-errChan
}
```

## Sorunlar

### Semaphore ve WaitGroup Kısmı

```go {linenos=table,hl_lines=["9-12","18-19"]}
func handleWebsites(websites []string) error {
	errChan := make(chan error, 1)
	semaphores := make(chan struct{}, 5) // aynı anda 5 iş çalıştır
	var wg sync.WaitGroup
	wg.Add(len(websites))
	for _, website := range websites {
		semaphores <- struct{}{}    // semaphore acquire et
		go func() {
			defer func() {
				wg.Done()
				<-semaphores
			}()
			if err := handle(website); err != nil {
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(semaphores)
	close(errChan)
	return <-errChan
}
```

Bu kısımlar kodumuzda bir panic oluşturmuyor, ancak aşağıdaki 2 durumdan birisi oluşuyor. 
1. `<-semaphores` işlemi `close(semaphores)` işleminden önce oluşabilir ve bu durumda zaten kanaldan bir değer okur. 
2. `close(semaphores)` işlemi daha önce gerçekleşir ve `<-semaphores` ise zero value alır. Önce `wg.Done()` operasyonu `wg.Wait()` fonksiyonun bitmesine ve `close(semaphores)` satırının çalışmasına yol açabilir. 

Her iki durumda da bir sıkıntı yok ancak bu kod fonksiyonun takibini daha zor yapıyor. Bunu go dilindeki şu tavsiyeyle çözebiliriz.


> Release locks and semaphores in the reverse order you acquired them.


Anlamı ise **locklar ve semaphoreları onları aldığınız sıranın tersinde bırakın.** Bu durumda kodumuz şu hale geliyor ve daha basit bir duruma dönüşüyor. 

```go {linenos=table,hl_lines=["9-12","18-19"]}
func handleWebsites(websites []string) error {
	errChan := make(chan error, 1)
	semaphores := make(chan struct{}, 5) // aynı anda 5 iş çalıştır
	var wg sync.WaitGroup
	wg.Add(len(websites))
	for _, website := range websites {
		semaphores <- struct{}{}    // semaphore acquire et
		go func() {
			defer func() {
				<-semaphores
				wg.Done()
			}()
			if err := handle(website); err != nil {
				errChan <- err
			}
		}()
	}
	wg.Wait()
	close(semaphores)
	close(errChan)
	return <-errChan
}
```

Şimdi ise sadece tek bir durum gerçekleşebilir o da `<-semaphores` işlemi channel kapanmadan okuma işlemlerini yapabilir çünkü `wg.Wait()` işlemi ancak ve ancak bütün `semaphores` kanalından okuma işlemleri gerçekleştikten sonra gerçekleşebilir.

### Semaphoreların Kullanımı

Semaphoreların kullanıldığı kısıma biraz daha yakından bakalım.

```go {linenos=table,hl_lines=[2]}
for _, website := range websites {
	semaphores <- struct{}{} // semaphore acquire et
	go func() {
            defer func() {
                    <-semaphores
                    wg.Done()
            }()
            if err := handle(website); err != nil {
                    errChan <- err
            }
        }()
}
```

`semaphores` channelı 5 uzunluklu bir channel olduğundan dolayı 5 goroutine çalıştıktan sonra 6. taska geldiğinde fonksiyon 2. satırda duracak ve bu `handle(website)` fonksiyonu bitene kadar durmayacak. Halbuki şöyle bir durum daha mantıklı olabilir.

Aynı anda 5 kez `handle(website)` fonksiyonu çalışsın, bir diğer deyimle goroutineler yaratılsın ve hazırda beklesin. Bunun için şu motto ile hareket edebiliriz.

> Acquire semaphores when you're ready to use them.

Anlamı ise `semaphoreları ne zaman kullanmaya hazırsan o durumda acquire et`.

```go
for _, website := range websites {
	go func() {
            semaphores <- struct{}{} // semaphore acquire et
            defer func() {
                    <-semaphores
                    wg.Done()
            }()
            if err := handle(website); err != nil {
                    errChan <- err
            }
        }()
}
```

Bu değişiklikten sonra artık bütün goroutineler yaratılır ve aynı anda ancak 5 tanesi sadece `handle(website)` fonksiyonunu çalıştırabilir.

### For Loop

`For-range` loop da yeni bir değişken `website` yaratıyoruz. Bir goroutine bu değişkeni updatelerken diğer goroutineler bu değişken üzerinden işlem yapıyor. Bundan dolayı burada bir data race var. Onun yerine 2 şekilde halledebiliriz.

#### Functiona parametre olarak verme

```go
for _, website := range websites {
	go func(website string) {
            semaphores <- struct{}{} // semaphore acquire et
            defer func() {
                    <-semaphores
                    wg.Done()
            }()
            if err := handle(website); err != nil {
                    errChan <- err
            }
        }(website)
}
```

#### Yeni Değişken Olarak Tanımlama

```go
for _, website := range websites {
        website := website
        go func() {
            semaphores <- struct{}{} // semaphore acquire et
            defer func() {
                    <-semaphores
                    wg.Done()
            }()
            if err := handle(website); err != nil {
                    errChan <- err
            }
        }()
}
```

Bundan ayrı olarak da genelde goroutineleri ayrı fonksiyonlara almak önerilir. Bu kod parçasını

```go
go func() {
        semaphores <- struct{}{} // semaphore acquire et
        defer func() {
                <-semaphores
                wg.Done()
        }()
        if err := handle(website); err != nil {
                errChan <- err
        }
}()
```

şu şekilde refactor edebiliriz.

```go
func handleWebsites(websites []string) error {
	errChan := make(chan error, 1)
	semaphores := make(chan struct{}, 5) // aynı anda 5 iş çalıştır
	var wg sync.WaitGroup
	wg.Add(len(websites))
	for _, website := range websites {
		go worker(website, semaphores, &wg, errChan)
	}
	wg.Wait()
	close(semaphores)
	close(errChan)
	return <-errChan
}

func worker(website string, sem chan struct{}, wg *sync.WaitGroup, errChan chan err) {
        semaphores <- struct{}{} // semaphore acquire et
        defer func() {
                <-semaphores
                wg.Done()
        }()
        if err := handle(website); err != nil {
                errChan <- err
        } 
}
```

### Error Channele Yazma

Bütün bu işlemleri yaptık ancak hala kodumuzda bir sorun var. Herhangi bir goroutine `errChan <- err` işlemini yaptığında diğer bütün error goroutineler bu kanala yazarken sonsuza kadar bekleyecekler ve bu da deadlock yaratacak. Bekleme sebebi `errChan` kanalının 1 uzunlukta bir channel olmasından dolayıdır. 

> Bir goroutine başlatmadan önce ne zaman ve nasıl duracağını bilmek gerekir.

Bunun yerine `select` ve `case` kullanarak sorunu halletmiş oluruz. 

```go {linenos=table,hl_lines=["22-25"]}
func handleWebsites(websites []string) error {
	errChan := make(chan error, 1)
	semaphores := make(chan struct{}, 5) // aynı anda 5 iş çalıştır
	var wg sync.WaitGroup
	wg.Add(len(websites))
	for _, website := range websites {
		go worker(website, semaphores, &wg, errChan)
	}
	wg.Wait()
	close(semaphores)
	close(errChan)
	return <-errChan
}

func worker(website string, sem chan struct{}, wg *sync.WaitGroup, errChan chan err) {
        semaphores <- struct{}{} // semaphore acquire et
        defer func() {
                <-semaphores
                wg.Done()
        }()
        if err := handle(website); err != nil {
                select {
                    case errChan <- err:
                    default:
                }
        } 
}
```

Bu durumda eğer herhangi bir goroutine `errChan`e yazabilirse yazacak ve yazamazsa `default` case çalışacak. Hiçbir goroutine bloklanmayacak. Select Case ile blocking çağrıları non-blocking olarak değiştirebiliriz.

### REFERENCES

- [Concurrency Made Easy From Dave Chevey](https://www.youtube.com/watch?v=DqHb5KBe7qI&t=643s)