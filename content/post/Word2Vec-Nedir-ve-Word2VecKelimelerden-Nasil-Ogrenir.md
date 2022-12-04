---
layout: post
title: Word2Vec Nedir ve Word2Vec Kelimelerden Nasıl Öğrenir
date: 2020-12-14
summary: Natural Language Processing'de kullanılan Word2Vec modelini inceliyoruz.
tags: [makine ogrenmesi, nlp]
math: true
---


Makine öğrenmesinde modellerin veriyi görme şekli biz insanlardan farklıdır. Biz kolayca ***Kırmızı arabayı görüyorum.*** cümlesini anlayabilirken, model bu kelimeleri anlayacak vektörlere ihtiyaç duyar. Bu vektörlere `word embeddings` denir. 

## WORD VECTORLERİ NASIL ÇALIŞIR - Tablodan Bak

Her kelimemiz için belirli bir boyutta vektörümüz olacak ve bu vektörleri kelimeyi isteyerek alabiliriz.

Buna key-value pair örneği verilebilir. 
* key: kelime 
* value: vektör

Bundan dolayı herhangi bir kelimenin vektörüne bakmak için dictionaryden kelimeyi istediğimiz zaman vektöre ulaşmış olacağız.


## Word2Vec: Tahmin Bazlı bir Metod

Ana amacımız kelimelerden, kelime vektörleri oluşturmak. 

Word2Vec parametreli word vektörleri olan bir modeldir. Bu parametreler itaretive yöntemle, objective function(küçültmeye çalıştığımız fonksiyon) kullanarak optimize edilir. 

Peki bunu nasıl yapacağız. 

Unutmadan:
* amaç : her bir vektörü kelimenin içeriğini bilecek şekilde kodlamak
* nasıl yapılacak: vektörleri kelimelerden olası içerik tahmin edecek şekilde eğitmek.


`Word2Vec` iterative bir metottur. Ana fikirleri kısaca şöyledir. 
* büyük bir text corpusu alır
* texti,  belirli bir sliding window(kayan pencere) kullanarak, her seferinde bir kelime ilerleyecek şekilde ilerlemek. Her bir adımda, bir tane `central word (merkezi kelime)` ve `context words(içerik kelimeleri)` -> penceredeki diğer kelimeler. 
* merkezi kelime için, içerik kelimelerinin olasılıklarını hesapla.
* vektörleri olasılıkları artıracak şekilde ayarla

![word2vec-nedir](../../images/training_data.png)

Resimde de görüleceği üzere her seferinde arkası mavi olan `merkezi kelime` ve diğerleri de `içerik kelimeleri`. 


### Objective Function (Amaç Fonksiyonu)

Text corpusundaki her bir $ t = 1, ... , T$ pozisyon için, Word2Vec merkezi kelimesi $w_{t}$ verilmiş m-boyutlu penceredeki içerik kelimelerini tahmin eder. 

$$
Likelihood = L(\theta) = \prod_{t=1}^{T} \prod_{-m \leq j \leq m, j \neq 0}  P(w_{t + j} \mid w_t, \theta) 
$$

Bu fonksiyonda $\theta$ optimize edilecek bütün parametrelerdir. Amaç ve Kayıp Fonksiyonu $J(\theta) ise ortalama negatif log olabilirlik fonksiyonudur. (Negative log-likelihood)

$$
J(\theta) = -\frac{1}{T} \log L(\theta) = -\frac{1}{T} \sum_{t=1}^{T} \sum_{-m \leq j \leq m, j \neq 0 } \log P(w_{t + j} \mid w_t, \theta)
$$


Bu formüldeki parçalara ayıralım.
* $\sum_{t=1}^{T}$ Bu kısım bütün text üzerinde gezinir. 
* $\prod_{-m \leq j \leq m, y \neq 0}$ bu ise kayma penceresini(sliding window) temsil eder.
* $\log P(w_{t + j} \mid w_t, \theta)$ : bu ise merkezi kelimesi verilen içeriğin olasılığını hesaplar.

Peki asıl sorulması gereken soru bu olasılıklar nasıl hesaplanacak?

### Olasılıkları Nasıl Hesaplayacağız?
Hesaplamak istediğimiz olasılık 
$$
P(w_{t + j} \mid w_t, \theta)
$$

Verilen her kelime $w$ için, iki adet vektörümüz var.

* $v_w$ -> kelimenin merkezi kelime (central word) olduğu zaman
* $u_w$ -> kelimenin içerik kelime (context word) olduğu zaman


Vektörler train edildikten sonra, genel olarak içerik vektörlerini $u_w$ atar ve sadece merkezi kelime vektörlerini $v_w$ kullanılır.

Bundan sonra verilen **merkezi kelime** $c$ ve **içerik kelimesi** $o$ kelimeleri için olasılık: 

$$
P(o \mid c) = \frac{exp(u_{o}^{T})}{\sum_{v \in V} exp(u_{w}^{T} v_c)}
$$

**NOT:** Bu bir `softmax fonksiyonudur. ` Softmax ile alakalı yazıma [bu yazımdan](https://ocakhasan.github.io/blog/Softmax-Aktivasyon-Fonksiyonu-Nedir-Numpy-Implementasyonu/) ulaşabilirsiniz.


Şimdi bu olasılıkları nasıl hesaplayacağımız gördüğümüze göre, vektörleri nasıl eğiteceğimizi görelim.

### NASIL EĞİTİLİR

Kısaca bu sorunun cevabı *Gradient Descent* ile her seferinde bir kelime alarak gerçekleşir. 

Parametrelerimiz $\theta$ bütün kelimelerin $v_w$ ve $u_w$ vektörleri olduğunu hatırlayalım. Bu vektörleri gradient descent kullanarak optimize edeceğiz. 

$$
\theta^{new} = \theta^{old} - \alpha \nabla_{\theta} J(\theta)
$$

Bu parametre güncellemerini her seferinde bir kelime kullanarak yapıyoruz. Her bir güncelleme bir merkez kelime ve içerik kelimesi ikilileriyle yapılır. Tekrardan kayıp fonksiyonuna bakalım.


$$
J(\theta) = -\frac{1}{T} \log L(\theta) = -\frac{1}{T} \sum_{t=1}^{T} \sum_{-m \leq j \leq m, j \neq 0 } \log P(w_{t + j} \mid w_t, \theta)
$$

Merkezi kelime $w_t$ için, kayıp fonksiyonu ayrı bir terimi her bir içerik kelimesi $w_{t + j}$ (sliding window içerisindeki)  $J_{t,j}(\theta) = -\log P(w_{t + j} \mid w_t, \theta)$ 

Bir örnek vererek bu durumu daha iyi anlayalım. Şu cümleyi ele alalım. 

Bugün bahçede <span style="color: green">bir</span> top gördüm. 

Yeşil renkli *bir* kelimesi burada bizim merkezi kelimemizdir. Her seferinde bir kelimeye bakacağımız için, bir tane içerik kelimesi seçeceğiz. Örnek olarak *top* kelimesini ele alalım. Bundan sonra bu iki kelime için kayıp fonksiyonu


$$
J_{t,j}(\theta) = -\log P(top \mid bir) = -log \frac{exp(u^{T}_{top} v_{bir})}{\sum_{w \ in V} exp(u^{T}_{w} v_{bir})} = -u_{top}^T v_{bir} + log \sum_{w \in V} exp(u^{T}_{w} v_{bir})
$$ 

Buradaki $V$ kümesi sliding windowu kapsayan kelimelerden oluşur. Loss (kayıp) fonksiyonumuzu aldığıma göre, şimdi vektörler üzerinde güncelleme yapalım. 


Burada hangi parameterlerin olduğuna göz atalım.
* merkezi kelime vektörlerinden sadece $v_{bir}$
* içerik kelime vektörlerinden ise sliding window içerisindeki bütün kelimeler $u_w \forall w \in V$

Şuanki adımda sadece bu parametreler güncellenecek. 

$$
v_{bir} := v_{bir}  - \alpha \frac{\partial J_{t, j}(\theta)}{\partial v_{bir}}
$$

$$
u_w = u_w - \alpha \frac{\partial J_{t, j}(\theta)}{\partial u_{w}} \forall w \in V
$$

Kayıp fonksiyonunu azaltacak şekilde yaptığımız her bir güncelleme, parametreler arasındaki benzerliği $v_{bir} \hspace{1mm} ve \hspace{1mm}  u_{top}$ dot product'ını artırıyor ve aynı zamanda diğer her bir diğer $u_w$ ile $v_{bir}$ arasındaki benzerliği de azaltıyor. 

Bu biraz garip gelebilir ancak neden *bir* kelimesinin *top* kelimesinden hariç diğer kelimelerle benzerliğini azaltmaya çalışıyoruz. Diğerleri de mantıklı, içerik verecek kelimeler olabilir. Ancak bu bir sorun değil! Biz bu güncellemeyi her kelime için tek tek yaptığımızdan dolayı, yani her kelime bir kez merkezi kelime olacak,  vektörler üzerindeki bütün güncellemelerin ortalaması metin içeriğininin dağılımını öğrenecektir. 

Bu yazıda partial derivative kısımlarına girilmemiştir. Ancak ben genel olarak **Word2Vec** modelinin nasıl çalıştığını anlatabildiğimi düşünüyorum. Eğer denemek isterseniz partial derivative kısımlarını kendiniz deneyebilirsiniz. Diğer yazılarda görüşmek üzere. 

Eğer yazıyı beğendiyseniz paylaşmayı unutmayın ki diğer insanlar da yararlansın. 

## REFERENCES
* https://lena-voita.github.io/nlp_course/word_embeddings.html
