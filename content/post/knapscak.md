---
layout: post
title: Dinamik Programlama ile Knapsack Problemi Nasıl Çözülür
summary: Meşhur Knapscak problemini dinamik programlama ile çözüyoruz. 
tags: [algoritmalar, python]
---


## Problem Tanımı

**Knapsack problemi** bilgisayar biliminde çok meşhur bir problemdir. Bu problemdeki amaç verilen ağırlık ve değerlerle en fazla değer toplayacak şekilde verilen ağırlık limitini aşmadan hangi itemlerin seçileceğidir. 




Knapscak problemi bir yüzyıldan fazla bir süredir, 1897 e kadar çalışmalar vardır. İsmini matematikçi [Tobias Dantzig](https://en.wikipedia.org/wiki/Tobias_Dantzig) adlı matematikçinin eski çalışmalarından alır. 

Buradaki problemimiz için birden fazla yöntem vardır. Biz dinamik programlama ile bu problemin nasıl hallediliğine bakacağız. 

Şimdi problemdeki input ve istenilen outputa bakalım

### INPUT
* Maksimum ağırlık limiti W ve elimizdeki paket sayısı *n*
* Ağırlıkların bulunduğu *w[i]* ve buna eş değer olan değer *v[i]*


### OUTPUT
* Maksimum değer
* Hangi paketlerin alındığı


### İMPLEMENTASYON
Bu problemi analiz ederken algoritmanın hangi değerlere bağlı olacağını bulmaktır. Buradaki algoritmamız 2 ayrı değişkene dayanır. Bunların birincisi kaç tane paket taşıyacağımız ve elimizde kalan ağırlık limiti.

Evet algoritmamızı iki değişkene bağlı şekilde yazacağız. Örnek olarak ilk 3 elemanı alarak, j maksimum limitli bir prpblemde optimum değer kaçtır. Buradaki ilk 3 eleman, hangi elemanları seçeceğimiz değişkenine örnektir. `J` limit ise ne kadar ağırlık limitimizin olduğudur. 

Bundan dolayı bir matrix oluşturup, her alt problemdeki optimum çözümü yazarsak bu şekilde istenilen sonuca ulaşabiliriz.

Bundan dolayı `[n+1][W+1]` boyutlarında bir matrixte elimizdeki her alt alt problem için çözümleri saklayacağız.

`K[i][j]` deki değer şu anlama geliyor: İlk `i` elemanı alarak *j* ağırlık limitli bir problemdeki optimum çözüm nedir.
Peki bizim soruda ne isteniyordu? `n` elemanı kullanarak W limitli bir problemdeki çözüm nedir. Bundan dolayı bizim istediğimiz sonuç ise matrixin en son elemanı olan `K[n][W]` dir. 

Peki asıl soru olan her bu `K[i][j]` nasıl bulacağız. 

Öncelikle matriximizin ilk satırının hepsi 0 olacak. Bunun nedeni ise 0 item kullanırsak elde edebileceğimiz maksimum değer 0 dır. 

Diğer satırlarda ise durum farklıdır. Bundan dolayı her $1\leq i \leq n$ ve her $0 \leq j \leq W$ için bir durumu kontrol etmemiz gerekiyor. Kontrol etmemiz gereken değişken şuanki durumda yani `item i` ağırlığı şuanki `j` (ağırlık limiti) büyük mü. Çünkü eğer bizim şuanki itemimizin ağırlığı ağırlık limitimizden büyükse o itemi basitçe alamayız. Bu itemi alamadığımız için `K[i][j] == K[i-1][j]` olacak. Nedeni ise bu itemi seçmediğimiz için o durumdaki en optimum çözüm, o iteme kadar olanki en optimum çözüme eşittir. 

``` python
if w[i] > j:
    K[i][j] = K[i-1][j]
```

Eğer bu durum gerçekleşme ise elimizde iki seçenek var. 

* Birinci seçenek şuanki itemi almamak. Bu durum yukarda bahsettiğimizin aynısı yani 

```python
K[i][j] = K[i-1][j]
```

* İkinci seçenek ise bu itemi almak. Bu durumda elimizde olan 
optimum çözüm şu anlama geliyor.
  * Şuanki itemin değeri `v[i]`
  * Bu itemi aldığımız için geriye kalan `j - w[i]` ağırlık limitli ve ilk `i-1` item arasındaki optimum çözüm, yani K[i-1][j-w[i]]

Bu durumda ise optimum çözüm şu anlama geliyor. 
```python
K[i][j] = v[i] + K[i-1][j - w[i]]
```

Elimizde iki seçenek var. Peki hangisini seçeceğiz. Çok basit, en yüksek olan hangisi ise bunu seçeceğiz. Yani 

```python
K[i][j] = max(K[i])[j], v[i] + K[i-1][j - w[i]])
```

Eğer bir basit bir kod yazmak istersek 

```python
for j in range(W+1)
    K[0][j] = 0             //Yani İlk satırı 0 yap

for i in range(1, n+1)
    for j in range(0, W + 1)
        if w[i] > j         // Eğer şuanki ağırlığımız şuanki limitten büyükse
            K[i][j] = 0
        else
            K[i][j] = max(K[i-1][j], v[i] + K[i-1][j - w[i]])
```

Bizim çözümümüz ise `K[n][W]` deki değerdir.


BackTracking kısmını sonra ekleyeceğim.





