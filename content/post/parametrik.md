---
layout: post
title: Makine Ogrenmesinde Parametrik Metodlar  
summary: İstatiksel biçimde parametrik metodları inceliyoruz.
tags: [makine öğrenmesi, matematik]
---

## TANIM
*İstatistik*, verilen bir örneklemden(sample) elde edilen herhangi bir değer demektir. İstatistiksel öğrenmede, verilen sampledan sağlanan bilgi ile karar verilir. İlk yaklaşımımız, sample'ın belirli bir dağılımdan (distribution) geldiğini farz ederek yapmak olacaktır. Bu dağılıma örnek olarak *Gaussian dağılım* verilebilir. Bu durumun avantajı ise, parametre sayısının azaltılması olacaktır. Tüm parametrelerimiz ortalama değer (mean) ve varyans (variance) olacaktır. Bu parametreleri sample tarafından elde ettikten sonra, bütün dağılımı biliyor olacağız. Bu parametreleri verilen sample üzerinden öğrenip, daha sonra bu bulduğumuz ortalama ve varyans değerlerini modele entegre ederek, tahmini bir dağılım elde edeceğiz. Daha sonra bu dağılımı da karar vermek için kullanacağız. 


Öncelikle olasılık kavramı diğer ismiyle density estimation (yoğunluk tahmini) anlamına gelen $p\left(x\right)$ kavramı ile başlıyoruz. Bu kavramı, Naive Bayesde de olduğu gibi tahmini olasılıkların $p(x \mid C_{i})$, ve prior olasılık olan $P\left(C_{i}\right)$ olduğu ve bu olasılıkların daha sonra asıl amaç olan $P\left(C_{i} \mid x\right)$'i tahmin ederek sınıflandırma işlemi yapılması için kullanıyoruz. Peki bu parametreleri nasıl öğreneceğiz. Maksimum Likelihood Estimation kullanarak yapacağız. 

### Maximum Likelihood Estimation (Maksimum Olasılık Tahmini)

Elimizde birbirinden bağımsız ve aynı şekilde dağıtılmış olan bir sample var. Bu sample'ı $X = \\{ x^{t} \\}_{i=1}^{N}$ şeklinde gösterebiliriz. Bu sampledan çekilen her bir $x^{t}$ örneğin, bilinen bir olasılık dağılımına ait olduğunu varsayıyoruz. Bu olasılık dağılımını da $p\left(x  \mid \theta \right)$ gösteriyoruz. 

$$
x^{t} \sim p(x|\theta)
$$

Bizim buradaki amacımız bize en yüksek olasılığı $p\left(x \mid \theta \right)$ verecek olan $\theta$ değerini bulmak.  Bütün örnekler $x^{t}$ birbirinden bağımsız olduğundan parametre $\theta$ nın olasılık fonksiyonu bütün verilen sampleların olasılıklarının çarpımına eşittir. 

$$
l(\theta | X) = p(X|\theta) = \prod_{t=1}^{N}p(x^{i}|\theta)
$$

Maksimum olasılık tahmininde, bu değeri maksimum yapan $\theta$ değerini bulmak istiyoruz. Bunu bulmak için önce logaritma alıp, daha sonra nerede maksimum yaptığına bakabiliriz. Logaritma alma sebebimiz ise logaritmanın çarpım sembolünü toplama çevirmesi ve başka kolaylıklar sağlaması dolayısıyladır. Log olasılık ise şöyle tanımlanır. 

$$
l(\theta | X)  \equiv \log l(\theta | X) = \sum_{t = 1}^{N}\log p(x^{t}|\theta)
$$

Yazımızın başında bu her sample ın belirli bir dağılımdan geldiğini söylemiştik. Bunun için bir sürü seçenek olabilir. *Bernouilli, Multinomial ve Gaussian(Normal)* dağılımlar olabilir. Ancak biz burada sadece **Gaussian(Normal)** dağılım ile ilgilineceğiz. 

### NORMAL DAĞILIMDA MAXİMUM LIKELIHOOD ESTIMATION

X, ortalama yani $E[X] \equiv \mu$ ve varyans $Var(X) \equiv \sigma^{2}$ değerlerine sahip normal dağılımla dağıtılmış bir random variable olsun. O zaman density (yoğunluk) fonksiyonu şu şekilde

$$
N(\mu , \sigma^{2}) = p(x) = \frac{1}{\sqrt{2\pi}\sigma}e^{-\frac{(x - \mu)^2}{2\sigma^{2}}}
$$

O zaman verilen sampleın $X = \\{ x^t \\}_{t=1}^{N}$ log likelihood değeri de şu şekilde olur.


$$
l(\mu, \sigma | X) = -\frac{N}{2}\log(2\pi) - N \log(\sigma)  - \frac{\sum_{t}(x^t - \mu)^{2}}{2\sigma^{2}}
$$

Daha sonra sırayla bu fonksiyonun ortalama değer ve varyans için partial türevlerini alıp sıfıra eşitlediğimizde ortaya şöyle bir sonuç çıkıyor.


$$
m = \frac{\sum_{t}x^t}{N}
$$

$$
s^2 = \frac{\sum_{t}(x^t - m)^2}{N}
$$

Burada $m$ `ortalama değer` için maximum likelihood estimate oluyor ve $s^2$ ise `varyans` için maximum likelihood estimate oluyor. Bu durumda istenilen parametreleri bulmuş olduk. Bundan sonraki yazıda ise bias(önyargı) ve Varyans(Variance) konuları işleyeceğiz.

Sonraki yazılarda görüşmek üzere. 





