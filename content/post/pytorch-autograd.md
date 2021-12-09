---
layout: post
title:  Pytorch AutoGrad Nedir ve Nasıl Çalışır
date: 2021-02-20
summary: Basit şekilde Pytorch Autograd ile otomatik olarak nasıl türev işlemleri halledilir.
tags: [pytorch, matematik]
---


## TANIM

[PyTorch](https://pytorch.org/) da bulunan `torch.autograd` otomatik türev alma motoru şeklinde çalışır ve bu da nöral ağ eğitimini güçlendirir. Bu yazımızda belirli örnekler vererek konunun daha geniş şekilde anlaşılmasını sağlayacağız. Öncelikle çok kısa bir özetleyici metine bakalım.


### Arka Plan

Nöral ağlar (neural networks) kendisine verilen veriyi belirli fonksiyonlarda işleyen bir bütündür. Bu fonksiyonların her biri bazı *parametrelerden* (ağırlıklar ve önyargı (weights and bias)) oluşur. Bu belirlenen parametrelere Pytorch da `tensor` adlı veri yapılarında tutulur. Bir Nöral ağın eğitilmesi iki kısımdan oluşur. Birinci kısımda sadece ileriye gidilir (forward propagation) ve ikinci kısımda geriye doğru gidilir (backward propagation). Peki bu ileri ve geri gitme işlemleri ne için yapılır onlara bakalım. 


#### **İleriye Gitme (Yayılma)**

Bu kısımda nöral ağ kendisine verilen veriden en iyi tahminini yapmaya çalışır. Bu belirlenen veri, önceden bahsettiğimiz her bir fonksiyondan geçer ve en sonunda bir tahmin ortaya atılmış olur. Daha sonra belirlenen tahmin ve gerçek değer arasından bir kayıp (loss) değeri bulunur ve hatta bu değeri bulan fonksiyona da `loss function` denilir. 


#### **Geriye Gitme (Yayılma)**

Bu kısımda ise nöral ağ, ilk bölümde hesaplanan kayıp veya hata değerini azaltmaya yönelik parametrelerinde iyileşmeye gider. Bunu yaparken de sonuçtan geriye dönük olarak her hata değerinin her bir parametreye bağlı olan türevini (derivative) hesaplar ve bu parametreleri, `gradient descent` kullanarak optimize eder. Ancak bu her bir fonksiyonun parametrelere göre türevini tek tek elimizle alamayız ve bize otomatik bir süreç lazım. 


İşte bu kısımda `pytorch.autograd` devreye giriyor ve bütün yükü alıyor. 


## Autograd'da Türev Alma İşlemleri 

Şimdi `autograd`'ın bütün bu değerleri nasıl kayıt ettiğine bakalım. 

Öncelikle iki tane `a` ve `b` `tensor` oluşturalım. Bu tensorları oluştururken parametre olan `requires_grad` parametresini `True` yapmamız gerekiyor aksi halde otomatik türev alma işlemleri gerçekleşemez çünkü bu parametre Tensorun `grad` adlı attributunda bu değerleri kayıt etmemize yardımcı oluyor. 

```python
import torch

x = torch.tensor([1., 2.], requires_grad=True)
y = torch.tensor([2., 4.], requires_grad=True)
```

Şimdi bu iki tensoru kullanarak yeni bir tensor `z` oluşturalım.

Basit şekilde formül 

$$
z = 6x^2 - 2b^3
$$

```python
z = 6*x**2 - 2*y**3
```

Şöyle varsayalım, `x` ve `y` bizim parametrelerimiz ve `z` bizim **hata fonksiyonumuz** olsun. Nöral ağ eğitiminde, hatanın bu parametreleri bağlı olan gradyantlarını (gradient) isteriz. 

PyTorch'da `.backward()` fonksiyonunu çağırdığımız zaman, `auutograd` her bir parametrenin (x, y) gradyantlarını bulur ve bunları her bir tensorun `.grad` attributunda kayıt eder. Öncelikle şuan `x` ve `y` nin grad değerlerine bakalım. 

```python
print("X.grad = ", x.grad)
print("Y.grad = ", y.grad)
```

```python
#Output
X.grad =  None
Y.grad =  None
```

Ancak şimdi `z` tensorunda `.backward()` çağırdığımız zaman `x` ve `y` nin `.grad` attributularında `z'nin` kendilerine göre türevler yer alacak. Ancak `z.backward()` argümanını çağırabilmek için parametre olarak `gradyant (gradient)`  vermemiz gerekiyor çünkü z bir vektör. `Gradyant` `z` ile aynı boyutlara sahip ve z'nin z ye göre türevini temsil eder. Şimdi `z.backward()` fonksiyonunu çağırabiliriz. 

```python
gradyant_parametre = torch.tensor([1., 1.])
z.backward(gradient=gradyant_parametre)
```

Şimdi `x.grad` ve `y.grad` değerleri oluşacak. Ancak bu değerleri görmeden önce kendimiz basit bir türev alalım. 

$$
\frac{\partial z}{\partial x} = 12x
$$

$$
\frac{\partial z}{\partial y} = -6y^2 
$$

Daha sonra bu kısmi türevlere `x` ve `y` tensorlarını koyduğumuz zaman ortaya çıkacak sonuçların şu şekilde olması lazım. 

```python
print("x için = ", 12 * x)
print("y için = ", -6 * y**2)
```

```
x için =  tensor([12., 24.], grad_fn=<MulBackward0>)
y için =  tensor([-24., -96.], grad_fn=<MulBackward0>))
```

Şimdi basit bir şekilde kontrol edelim.

```python
print("x.grad = ", x.grad)
print("y.grad = ", y.grad)
```

```python
x.grad =  tensor([12., 24.])
y.grad =  tensor([-24., -96.])
```

Gördüğümüz üzere sonuçlar doğru çıkıyor. Üstte gözüken `ggrad_fn=<MulBackward0>` ise bu bu tensorun nasıl bir matematiksel operatör kullanarak oluşturulduğunu söylüyor. Eğer `required_grad=False` olsaydı bu değer `None` olurdu. 

Bütün yazdığımız operasyonlar için `required_grad=True` idi. Şimdi `required_grad=False` yapıp bir de öyle deneyelim.

```python
x = torch.tensor([1., 2.], requires_grad=False)
y = torch.tensor([2., 4.], requires_grad=False)
z = 6*x**2 - 2*y**3

gradyant_parametre = torch.tensor([1., 1.])
z.backward(gradient=gradyant_parametre)

print("X.grad = ", x.grad)
print("Y.grad = ", y.grad)
```
Bu işlemden şöyle bir sonuç alacaksınız.
```
RuntimeError: element 0 of tensors does not require grad and does not have a grad_fn
```

Bu da demek oluyor ki `x` ve `y` nin grad değerleri yok ve bundan dolayı `grad_fn` fonksiyonları da yok. `x` ve `y` nin grad değerleri olmadığı için `z`'nin de grad değeri yok ve bu da hataya yol açıyor. 


Derin öğrenmede genellike önceden belirli datasetler ile eğitilmiş hazır modeller bulunmaktadır ve bunlara *pretrained model* denir. Bu modelleri kullanırken genellikle son katmana kadar olan bütün katmanların parametrelerini eğitmek istemeyiz çünkü bu işlem hem pahalı hem de çok da gerekli olmayan bir işlem. Bu parametreleri optimize etmeye çalışmadığımızdan dolayı bu parametrelerin `required_grad` değerleri `False` olacaktır. Örnek bir kod olarak da 

```python
from torch import nn, optim

model = torchvision.models.resnet18(pretrained=True)

# Freeze all the parameters in the network
for param in model.parameters():
    param.requires_grad = False
```

Burada [Resnet18](https://pytorch.org/hub/pytorch_vision_resnet/) modelinin parametrelini dondurma (freeze) işlemi yapılıyor ve böylece modeli kullanırken `resnet18` modelini parametrelerinde herhangi bir optimize etme durumu söz konusu olmayacak. Ancak, sonradan eklenilen layerlarda optimize çalışması yapılabilir o kadar. 

## REFERENCES
[Pytorch Tutorials](https://pytorch.org/tutorials/beginner/blitz/autograd_tutorial.html#sphx-glr-beginner-blitz-autograd-tutorial-py)


