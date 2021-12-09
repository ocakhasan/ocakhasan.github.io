---
layout: post
title:  Evrişimsel Sinir Ağları  (Convolutional Neural Network) Nedir
summary: Derin Öğrenmede resimler üzerinde kullanılan evrişimsel sinir ağlarını kullanıyoruz.
date: 2020-12-17 
tags: [numpy, pytorch, cnn]
---

Yazıya başlamadan önce belirmek isterim ki, bu tarz derin öğrenme terimlerinin İngilizce ile kullanılması taraftarıyım. Teknik terimlerin Türkçe karşılıkları genelde her zaman duymadığımız kelimeler oluyor ve internette Türkçe pek kaynak yok. Ondan dolayı ben bu terimlerin İngilizce öğrenilip, İngilizce kullanılması taraftarıyım. Herkes global olmaya çalışırken, bizim öyle davranmamamız için hiçbir sebep yok. 


## TANIM
Convolutional sinir ağları genel olarak sıradan sinir ağlarına çok benzerdir. Bu sinir ağları da öğrenebilir ağırlık (weight) ve önyargısı (bias) olan sinirlerden (neuron) oluşur. Her bir nöron bazı inputlar alır, dot product uygular ve bu işlemi lineer olmayan bir yolla devam ettirir. Bütün network hala tek bir ayırt edilebilir skoru açıklar. Network resim pixellerini alıp, sonda bir tahmin üretir. Networkun sonunda belirli  bir kayıp fonksiyonu (loss function) bulunur. 

Peki bu convolutional sinir ağları normal sinir ağlarına bu kadar benziyorsa ne değişiyor? Bu sorunun cevabı ise şu şekildedir: 

Convolutional sinir ağları inputun resimlerden oluştuğunu varsayar, bu varsayım bize bazı özellikleri sisteme entegre etmemize yardımcı olur. 


## YAPISAL GÖZLEM

**Normal Sinir Ağları:** Normal sinir ağları tek bir input alır, onu bazı gizli katmanlardan (hidden layer) geçirir. Her bir hidden layer nöron kümelerinden oluşur,  her bir nöron, bir önceki katmandaki bütün nöronlarla bağlantılıdır ve diğer nöronlardan bağımsız şekilde çalışır. Son katman ise sonuç katmanı (output layer) olarak adlandırılır ve bu katmanda her bir sınıfın olasılığı belli olur. 


Bu normal sinir ağları resimler kullanılınca pek iyi ölçeklenemiyor. Örnek olarak $(32, 32, 3)$ lük boyutlarda resimler kullanırsak, ilk katman $32 * 32 * 3 = 3072$ ağırlığa sahip olacaktır. Bu yük halledilebilir şekilde görülüyor ancak, bu fully-connected yapı büyük resimlere ölçeklenmiyor. Örnek olarak eğer biz boyutları $(200, 200, 3)$ olan resimler kullanırsak, bu sefer nilk nöronlar $200 * 200 * 3 = 120, 000$ ağırlığa sahip olacaklar. Ancak bu büyük numaralı ağırlıklar aşırı uyma (*overfitting*) denilen olaya sebep olacaktır. 


Convolutional sinir ağları ise inputun resimlerden oluşmasınından faydalanır ve buna göre yapıyı daha mantıklı şekilde kurar. Normal sinir ağlarının aksine, Convolutiona sinir ağlarının nöronları 3 boyuta ayarlanmış şekildedir. **genişlik, yükseklik, derinlik**.

Örnek olarak $(32, 32, 3)$ boyutlu resimlerde
* Genişlik = 32
* Yükseklik = 32
* Derinlik = 3
olacaktır. 


## PEKI BU CONVOLUTIONAL SINIR AĞLARI NASIL OLUŞTURULUYOR? 

Bu sinir ağları katman dizilerinden oluşur ve bu katmanlar ise şu şekildedir. 

* Convolutional Katman
* Pooling Katmanı
* Fully-Connected Katmanı

Bu 3 katmandan oluşan katmanları birleştirip bir sinir ağı oluşturacağız. 


#### CONVOLUTIONAL KATMAN

Convolutional katman Convolutional sinir ağlarının büyük ağır işini yapan katmanlardır. 

Conv katmanlar parametreleri öğrenilebilir filtrelerden oluşur. Her bir filtre boyut olarak küçüktür, ancak input derinliği boyunca uzanırlar. Örnek olarak, tipik bir filtre $5 * 5 * 3$ boyutlarında olabilir. İlk 5 genişlik, ikinci 5 yükseklik ve üçüncü 3 ise resimin 3 derinlikli olmasından kaynaklanır. Doğrudan iletme kısmında, her bir filtreyi input resmi üzerinde kaydırıyoruz, bu kaydırma sırasında resimlerde pixeller ile filtredeki sayılar ile dot product alıyoruz. Filtreyi kaydırma işlemi sırasında 2 boyutlu bir aktivite haritası oluşturuyoruz. Bu harita ise bize her bir pozisyondaki cevabı veriyor. Sinir ağı, bu filtreler ne zaman belirli bir görsel özellik, örnek olarak kenar, gördüğü zaman öğrenecek. Her bir filtrenin oluşturduğu haritaları üst üste sıkıştırıp bunu bir sonraki katmana iletiyoruz. 


![Convolutional Sinir Ağları Örnek](../../images/cnn.png)


**BOYUTSAL AYARLAMA**

Her bir nöronun nasıl bağlı olduğunu anlattık ancak output hacminde kaç tane nöron olduğundan bahsetmedik. Output hacmini belirleyen 3 ayrı parametre vardır. 
* **DERİNLİK:** Bu parametre kaç tane nöron kullandığınıza işaret eder. Örnek olarak ilk convolutional katman input olarak resmi alırken, farklı nöronlar bu resimde farklı detayları fark edebilir. 
* **STRIDE (KAYDIRMA ADIMI):** Bu parametre ise filtreyi kaç pixel kaydıracağımıza işaret eder. Eğer `stride` bir ise, filtreleri bir pixel kaydıracağımız anlamına gelir. 
* **ZERO-PADDING: (SIFIRLARLA DOLDURMA)** Bazı durumlarda inputun etrafını sıfırlarla doldurmak uygun olmaktadır. Bu işlemin güzel bir tarafı ise, bize output boyutunu kontrol altında tutma olanağı vermesidir. Örnek olarak daha yüksek boyutlu outputlar istersek, inputu filtre boyutu kadar sıfırlarla doldurup, bir sonraki katmana aktarılacak outputun boyutunu, şuanki katmandaki input boyutuna eşit tutabiliriz. 


Output hacminin boyutunu şu şekilde hesaplayabiliriz. 
* Input Boyutu = $W$
* Convolutional katman nöronları filtre boyutu = $F$
* Stride = $S$
* Zero-Padding = $P$

Output hacmi boyutu formülü = $(W - F + 2P) / S + 1$. 

Örnek olarak eğer elimizde $10 * 10$ boyutlu bir resim varsa ve bizim filtre boyutumuz $3 * 3$, stride = $1$ ve padding = $0$ ise

$$
Output Boyutu = (10 - 3 + 2*0) / 1 + 1 = 8 * 8 
$$

Şimdi bu boyut tek bir nörondan çıkan sonuç. Eğer elimizde $n$ tane nöron varsa, bu katmandan çıkan sonucun boyutu $8 * 8 * n$ olacaktı. 

![Convolutional Sinir Ağları Örnek](../../images/cnn_filter.png)

Yukarıdaki örnekten de görüleceği üzere filtre boyutumuz $3 *3$, bundan dolayı resimde de (3*3) lük alanlar alıp, bu aldığımız alanla filtre arasında bir dot product işlemi uyguluyoruz. Peki resimdeki $31$ sayısına nasıl ulaştık onu inceleyelim.

$$
(1 * 1) + (0 * 2) + (1 * 3) + (0 * 4) + (1 * 5) + (1 * 6) + (1 * 7) + (0 * 8) + (1 * 9)
$$

$$
1 + 3 + 5 + 6 + 7 + 9 = 31
$$

Özetlemek gerekirse
* Conv layer $W_1  * H_1 * D_1$ boyutlarında input alır
* 4 parametreye ihtiyaç duyar
  * Filtre sayısı $K$
  * Filtrenin boyutları $F$
  * Stride $S$
  * Zero padding sayısı $P$
* $W_2 * H_2 * D_2$ boyutlarında output üretir.
  * $W_2 = (W_1 - F + 2P)/S + 1)$
  * $H_2 = (H_1 - F + 2P)/S + 1$
  * $D_2 = K$

### PYTHON İLE UFAK BİR GÖSTERİM

Şimdi `tensorflow` ile basit bir gösterim yapıp bu boyutları daha iyi anlayalım.

```python
import tensorflow as tf
# The inputs are 28x28 RGB images with `channels_last` and the batch
# size is 4.
input_shape = (4, 28, 28, 3)
x = tf.random.normal(input_shape)
y = tf.keras.layers.Conv2D(
2, 3, activation='relu', input_shape=input_shape[1:])(x)
print(y.shape)
```

```python
(4, 26 , 26, 2)
```

Burada olan işlemler şu şekildedir
* `input_shape` Conv layer'a verilecek olan inputun boyutlarıdır. (4, 28, 28, 3) şu anlama gelmektedir. Bizim elimizde 4 adet resim var, ve bu resimlerin boyutları (28, 28, 3)tür. 
* `Conv2D` ' e verilen parametreler ise şu şekildedir. İlk verilen parametre `2` kaç adet filtre kullanacağımızı gösterir. İkinci parametre `3` ise filtre boyutunu vermektedir. Yani filtre boyutumuz $(3, 3)$ olacaktır. 

Şimdi burada oluşan outputun nasıl oluştuğuna bakalım. Yukarıda özetlediğimiz gibi her şeyi tek tek yazalım

Input boyutları $W_1 * H_1 * D_1$ şeklindeydi. Bundan dolayı
* $W_1 = 28$
* $H_1 = 28$
* $D_1 = 3$

Daha sonra filtre sayımız $K = 2$, filtre boyutumuz ise $F = 3$, stride ise default olarak $S = 1$dir. Padding ise default olarak $P = 0$dır. 

O zaman şimdi output boyutlarımızı $(W_2 * H_2 * D_2)$ hesaplayabiliriz.

* $W_2 = (28 - 3 + 2 * 0) / 1 + 1  = 25 + 1 = 26$
* $H_2 = (28 - 3 + 2 * 0) / 1 + 1 = 25 + 1 = 26$
* $D_2 = K = 2$

Her bir resim için oluşturulan output boyutları $(26, 26, 2)$. Elimizde 4 adet resim var ve bundan dolayı çıkan output boyutu $(4, 26, 26, 2)$

## POOLING LAYER

Convolutional sinir ağlarında convolutional katmanlar arasına *Pooling* katmanları eklemek çok yaygındır. *Pooling* katmanının görevi verilen inputun boyutlarını kademeleri olarak azaltarak parametrelerin ve ağın işlem yükünün azaltımasını sağlamak. Bu şekilde aşırı uyma (overfitting) kontol altına alınmış olur. Pooling katmanı, bağımsız olarak çalışır ve her bir inputu `Max` operasyonu kullanarak boyutlarını azaltır. En yaygın *Pooling* katmanı, filtreleri $(2 * 2)$ boyutlarında olan ve inputu hem boydan ve hem enden ikiye bölenlerdir. Her bir `Max` operasyonu input olarak $(2 * 2)$ lik bir bölüm alacak ve bu 4 sayıdan en büyüğünü gönderecektir. Özetlemek gerekirse,

* Pooling katmanı input olarak $W_1 * H_1 * D_1$ boyutlarını kabul eder. 
* İki parametreye ihtiyaç duyar
  * Boyut $F$
  * Stride $S$
* Boyutları $W_2 * H_2 * D_2$ olan output çıkarır. 
  * $W_2 = (W_1 - F)/S + 1$
  * $H_2 = (H_! - F)/S + 1$
  * $D_2 = D1$



![Convolutional Sinir Ağları Pooling Layer Örneği](../../images/pooling.png)

Resimde de görüleceği üzere her bir $(2 * 2)$ lik bölümden en büyük sayılar alınıp yeni bir örnek elde ediliyor.

Şimdi bunu Python ile kodlamaya çalışalım. 

### Pooling Layer Python İle İmplementasyonu

Bu layerı hem sıfırdan hem de kütüphane kullanarak kodlayabiliriz. Önce kütüphane kullanarak gösterelim. 


```python
x = tf.constant([[1., 2., 3.],
                 [4., 5., 6.],
                 [7., 8., 9.]])
x = tf.reshape(x, [1, 3, 3, 1])
max_pool_2d = tf.keras.layers.MaxPooling2D(pool_size=(2, 2),
   strides=(1, 1), padding='valid')
max_pool_2d(x)
```

Bu koddan çıkacak output ise 

```python
<tf.Tensor: shape=(1, 2, 2, 1), dtype=float32, numpy=
  array([[[[5.],
          [6.]],
          [[8.],
          [9.]]]], dtype=float32)>
```

Çıkan sonucun nasıl çıktığını bence rahatlıkla yapabilirsiniz. 

Şimdi kendimiz sıfırdan bu layerı basit bir şekilde implement edelim.

```python
import numpy as np

def pool2d(X, pool_size, mode='max'):
    p_h, p_w = pool_size        #pool size ı al
    Y = torch.zeros((X.shape[0] - p_h + 1, X.shape[1] - p_w + 1)) #Outputu oluştur
    for i in range(Y.shape[0]):
        for j in range(Y.shape[1]):
            Y[i, j] = X[i: i + p_h, j: j + p_w].max()   #Her bir pool size kadar pixelin max'ını al
            
    return Y
```

Şimdi kodumuzu yukarıda yazdığımız `x` arrayi ile test edersek, yine aynı sonucun çıkacağını göreceğiz. 


Bu yazımızda konuşulacaklar bu kadar. Beğendiyseniz paylaşmayı unutmayın.

## REFERENCES
* https://anhreynolds.com/blogs/cnn.html
* https://cs231n.github.io/convolutional-networks/
* https://cezannec.github.io/Convolutional_Neural_Networks/
* https://www.tensorflow.org/api_docs/python/tf/keras/layers/Conv2D
* https://www.tensorflow.org/api_docs/python/tf/keras/layers/MaxPool2D
* https://medium.com/ai-in-plain-english/pooling-layer-beginner-to-intermediate-fa0dbdce80eb



