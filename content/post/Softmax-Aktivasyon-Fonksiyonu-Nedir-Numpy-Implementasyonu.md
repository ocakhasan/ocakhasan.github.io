---
layout: post
title: Softmax Aktivasyon Fonksiyonu Nedir ve Numpy ile Nasıl Implement Edilir 
date: 2020-08-12
summary: Derin Öğrenmede kullanılan softmax aktivasyon fonksiyonunu inceliyoruz.
tags: [numpy, derin ogrenme, matematik]
math: true
---


## TANIM

Softmax fonksiyonu modelden çıkan sonuçların olasılıksal şekilde ifade edilmesi için kullanılan bir fonksiyondur. Genellikle nöral ağlarda (neural network) ağın sonucunu sınıflara olasılık değerleri vermek için kullanılır. 

Softmax fonksiyonu input olarak $K$ boyutlu uzaydan vektör $z$ alır. Bu vektörü $K$ olasılık değerlerinden oluşan bir olasılık dağılımına çevirir. Bu olasılıkların her biri exponentialları ile doğru orantılıdır. Softmax fonksiyonu uygulamadan önce bu $z$ vektöründeki bazı değerler negatif de olabilir 0 da olabilir, pozitif de olabilir. Softmax fonksiyonunu uyguladıktan sonra ise bütün değerler $(0, 1)$ aralığında değer alır ve bütün değerlerin toplamı 1 olur. 

Standart softmax function tanımı şu şekildedir. $\sigma : \mathbb{R^{K}} \rightarrow  \mathbb{R^K}$ 

$$
\sigma(z)_{i} = \frac{e^{z_i}}{\sum_{j=1}^{K} e^{z_j}} her \hspace{1mm} i = 1, 2, 3, ...,  K ve  \hspace{1mm} z = (z_1, z_2, ... , z_k) \in \mathbb{R^K}
$$

Bir diğer deyişle bizim yaptığımız işlem her bir değerin exponential fonksiyonunu almak ve bunu toplama bölmek. Böylece normalize etmiş oluyoruz ve bütün değerleri topladığımız zaman sonuç 1 ediyor. 

Örnek olarak vektör $k = [1, 1, 1] \in \mathbb{R^3}$ olsun. O zaman,

$$
\sigma(k) = [\frac{1}{3}, \frac{1}{3}, \frac{1}{3}]
$$

Peki bunu nasıl bulduk. Öncelikle toplamı hesaplayalım

$$
\sum_{j=1}^{K}e^{k_i} = e^1 + e^1 + e^1 = 3e
$$

Toplam $3e$ çıktı. Şimdi her bir değeri exponential fonskiyona input olarak verirsek çıkacak sonuç $e^1 = e$ olur. Bizim softmax fonksiyonunda yaptığımız işlem ise bu değerleri alıp toplama bölmek. Yani,

$$
\frac{e}{3e} = \frac{1}{3}
$$


## NUMPY İLE NASIL İMPLEMENT EDİLİR

Numpy fonksiyonunda arrayin direk exponential fonksiyonunu alabiliriz. Bunun için for loop açmamıza gerek yok

```python
import numpy as np

arr = np.array([1, 3, 2])

exponential_arr = np.exp(arr)

print("Array:   {} \nExponential Array: {} \n".format(arr, exponential_arr))
```

```python
Array: [1 2 3] 
Exponential Array : [ 2.71828183  7.3890561  20.08553692] 
```

Arrayin direk üstel şekilde toplamını da alabiliriz. 

```python
sum_of_exponentials = np.sum(exponential_arr)
print("Exponential Array Toplamı: ", sum_of_exponentials)
```

```python
Exponential Array Toplamı:   30.19287485057736
```

Şimdi softmax implement etmek için her şeye sahibiz. Fonksiyon şeklinde implement edebiliriz.

```python
def softmax(arr):
    exp_array = np.exp(arr)
    exp_toplam = np.sum(exp_array)
    return exp_array/exp_toplam
```

Şimdi fonksiyonumuzu test edelim

```python
arr = np.array([1, 1, 1])
softmax_array = softmax(arr)
print("Array: {} \nSoftmax Array: {}".format(arr, softmax_array))
```

```python
Array: [1 1 1] 
Softmax Array: [0.33333333 0.33333333 0.33333333]
```

Gördüğümüz üzere softmax fonksiyondan çıkan arrayin toplamı 1 e eşit oluyor
```python
np.sum(softmax_array) #Sonuç 1 çıkıyor. 
```

Softmax fonksiyonu bu kadar. Bir sonraki yazıda görüşmek üzere. 

## REFERENCES

[wikipedia-softmax](https://en.wikipedia.org/wiki/Softmax_function)





