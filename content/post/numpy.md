---
layout: post
title: Python Numpy ile Sifirdan K Nearest Neighbours Algoritmasini Yazalim 
summary: Bu yazıda Python Numpy ile Sıfırdan Knn yazıyoruz. 
tags: [numpy, knn, makine öğrenmesi]
---

Merhaba bu yazımızda Makine Öğrenmesinde meşhur bir algoritma olan Knn algoritmasını sıfırdan yazacağız. Tabii ki
hazır bir sürü kütüphane var ancak sıfırdan algoritmayı yazabilmek bize algoritmanın nasıl çalışacağını gösterecektir. 
Böylece Knn algoritması bir tahmin yaparken nasıl yapıyor olayın arkasında neler dönüyor bunları anlayabiliyor olacağız. 


## K-Nearest Neighbour Nedir

Öncelikle şunu bilmek gerekir ki K-Nearest-Neighbour adından da anlaşılacağı üzere en yakın k komşu noktalara bakıp en çok hangi label varsa o labelı tahmin(prediction)  olarak verir. 

Peki bu yakınlık uzaklık ilişkisi nasıl kurulur önce ona bakalım. Uzaklığı ölçebilmek için belli başlı algoritmalar vardır. Bunlardan biri `eucledian` diğeri de `manhattan` uzaklığıdır. 

### Eucledian Uzaklığı 

Manhattan uzaklığında aslında iki nokta arasında uzaklığı alırken normal 2 boyutlu denklemde nasıl alıyorsak, bunun n boyutlu formüle döndürülmüş halidir. 

Örnek olarak $a = (x_1, y_1)$ ve $b = (x_2, y_2)$ olsun. Bu noktalar arasında uzaklığı bulurken yaptığımız işlem 

$$d(a, b) = \sqrt{(x_1 - x_2)^2 + (y_1 - y_2)^2}$$

Peki eğer bizim verimiz n boyutlu olursa bu uzaklık nasıl ölçülecek?. Bu durumda ise uzaklık 

$$d(a,b)= \sum_{i=1}^n (a_i - b_i)^2$$

Bu formulu ise Numpy ile şu şekilde yazabiliriz

```python
np.sqrt(np.sum(np.square(a - b), axis=1))
```

### Manhattan Uzaklığı 

Manhattan uzaklığında iki nokta arasındaki uzaklık her bir alt noktanın farkının mutlak değerlerinin toplamı ile bulunur. 


Örnek olarak $a = (x_1, y_1)$ ve $b = (x_2, y_2)$ olsun. Bu noktalar arasında uzaklığı bulurken yaptığımız işlem

 $$d(a, b) = \lvert x_1 - x_2\rvert + \lvert y_1 - y_2 \rvert$$


Peki eğer bizim verimiz n boyutlu olursa bu uzaklık nasıl ölçülecek?. Bu durumda ise uzaklık

 $$d(a,b)= \sum_{i=1}^n \lvert a_i - b_i\rvert$$

Bu formulu ise Numpy ile şu şekilde yazabiliriz
```python
np.sum(np.abs(a - b), axis=1)
```
## Algoritma Akışı

KNN algoritmasında eğitme (training) işlemi aslında sadece verilen veriyi ezberlemekten ibarettir. Bundan dolayı eğitme kısmında bir şey yapmayacağız ancak tahmin etme (prediction) kısmında ise asıl üstteki formuülleri kullanıp işlem yapacağız. 

Bu algoritmayı Python ile Numpy  Kullanarak implement edeceğiz. 
Önceklikle şu komut ile Numpy kütüphanesini import edelim.

```python
import numpy as np
```

Şimdi Bir tane class tanımlayalım. 
Öncelikle kaç tane komşu kullanacağı modelin bir parametresi `k` olacak. Daha sonra hangi uzaklık formülünü kullanacağı da modelin bir parametresi olacak. Hadi başlayalım.

```python
class KNN:
    def __init__(self, k=2, uzaklık="eucledian"):
        self.k = 2
        self.uzaklık = uzaklık
```

Üstteki kod bloğunda yaptığımız işlem aslında modelin parametrelerini `constructor` fonskiyonunda tanımlamak oldu. 

Şimdi modelin eğitme fonksiyonunu yazalım. 

```python
class KNN:
    def __init__(self, k=2, uzaklık="eucledian"):
        self.k = 2
        self.uzaklık = uzaklık

    def fit(self, X, y):
        self.X = X
        self.y = y
```

Yukarıda da bahsettiğimiz gibi Knn algoritması aslında sadece eğitme verisini ezberler. Bütün işlemler prediction kısmında yapılır. Bundan dolayı eğitme verisini modelin eğitim seti olarak değiştirebiliriz. 
Şimdi en önemli konu olan prediction kısmına gelelim. 

### Tahmin Etme (Prediction) Kısmı Nasıl Olacak?

1. Öncelikle uzaklıklar hesaplanacak
2. Daha sonra en yakın k tane nokta alınacak
3. Daha sonra bu en yakın k noktanın label sayılarını belirleyeceğiz. Yani hangi labeldan kaç tane olduğunu belirleyeceğiz.
4. Daha sonra en çok sayısı olan label bizim tahminimiz olarak

Diyelim ki bizim elimizde bir tane a verisi olsun ve eucledian uzaklık kullanıyor olalım. Bir nokta için nasıl bir işlem gerekiyor ona bakalım.
```python
#Uzaklık Hesaplama
uzaklıklar = np.sqrt(np.sum(np.square(X - a), axis=1))

#En yakın k noktanın indexlerini bulalım
en_yakın_k_index = np.argsort(uzaklıklar)[:k]

#Şimdi bu en yakın k indexin hangi labellara ait olduğunu bulalım.
en_yakın_labellar = y[en_yakın_k_index]

#Daha sonra bu en yakın labellarda her labeldan kaç tane olduğunu bulalım
labellar, adetler = np.unique(en_yakın_labellar, return_counts=True)

#Daha sonra en çok hangi labelın sayısı var bunu bulalım
max_label_index = np.argmax(adetler)

#Daha sonra en çok sayısı olan label döndürelim
return labellar[max_label_index]
```
Şimdi burada bir sürü numpy fonksiyonu kullandık, bunlar kafa karıştırmış olabilir. Bundan dolayı bu fonksiyonlar ne işe yarıyor kısaca anlatayım. 

`np.argsort(array)`: Bu fonksiyon parametre olarak aldığı `array`i sortlayacak indexleri verir. Aslında arrayi sortlamaz, ancak hangi indexler arrayi sortlar onu verir. 

`np.unique(array, return_counts=True)`: Bu fonksiyon ise `array` içerisindeki unique elemanları dönderir. Eğer `return_counts=True` ise o zaman bu unique elemanlardan kaç tane var onu da gösterir. Örnek olarak
```python
array = [2, 3, 4, 3, 2, 10, 2]
labellar, sayılar = np.unique(array, return_counts)
print(labellar) #(2, 3, 4, 10)
print(sayılar)  #(3, 2, 1, 1)
```
Labellar bizim arrayimizde hangi unique label var onları dönderir. Sayılar ise hangi labeldan kaç tane var onu dönderir. Örnek olarak 2 den 3 tane var. 10'dan 1 tane var. 

`np.argmax(array)`: Bu fonksiyon array içerisindeki maximum elemanın indexini verir. Mesela yukarıdaki arrayde maximum eleman 10 ve 10'un indexi 5 dir. `np.argmax()` bu 5 indexini dönderir. 
 
Şimdi bu bir nokta içindi. Bunu test verimizde her bir nokta için yapalım. `KNN Classının` içerisinde implement edelim.


```python
class KNN:
    def __init__(self, k=2, uzaklık="eucledian"):
        self.k = 2
        self.uzaklık = uzaklık

    def fit(self, X, y):
        self.X = X
        self.y = y
    
    def predict(self, X_test):
        predictions = []
        for point in X_test:
            if self.uzaklık == "eucledian":
                uzaklık = np.sqrt(np.sum(np.square(self.X - point), axis=1))
            elif self.uzaklık == "manhattan":
                uzaklık = np.sum(np.abs(self.X - point), axis=1)
        
            indices = np.argsort(uzaklık)[:self.k]                        
            near_labels = self.y[indices]                                   
            labels, values =  np.unique(near_labels, return_counts=True)      
            max_ind_label = np.argmax(values)                               
            prediction = labels[max_ind_label]    
            predictions.append(prediction)
        
        return np.array(predictions) 

```
Yaptığımız işlemler

* her nokta için uzaklığı hesapladık
* en yakın k noktanın labellarını aldık
* bu labelların sayılarını öğrendik
* en çok label kimdeyse onu tahmin olarak öne sürdük


Bu `KNN` modelini istediğiniz veride kullabilirsiniz. Daha kapsamlı koda bakmak isterseniz bu [Github Repo'ya](https://github.com/ocakhasan/machine-learning-from-scratch) bakabilirsiniz. 

Bir sonraki yazıda görüşmek üzere. 