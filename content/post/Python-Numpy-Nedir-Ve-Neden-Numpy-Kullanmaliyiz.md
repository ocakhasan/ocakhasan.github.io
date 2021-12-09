---
layout: post
title: Kapsamli Şekilde Python Numpy Öğrenelim
summary: Bu yazıda Python Numpy Kütüphanesini inceliyoruz. 
date: 2020-11-09
tags: [numpy]
---
# NUMPY 

![Numpy Logo](../../images/numpy_logo.png)
```python
import numpy as np
```

## NUMPY ARRAY

Python'daki listelere çok benzerdir. Numoy arrayleri sadece aynı veri türüne sahip listeleri barındırabilir. Arrayler daha az hafızada yer kaplar.

Array oluşturmak için yapmamız gereken


```python
a = np.array([1, 2, 3, 4])
type(a)
```

```python
numpy.ndarray
```


Buna ek olarak da, arraylere tuple ekleyebiliriz. Örnek olarak


```python
a = np.array([[1, 2, 3], [4, 5, 6], [7, 8, 9]])
print("Birinci eleman {}".format(a[0]))
print("Birinci elemanın birinci elemanı {}".format(a[0][0]))
```

```python
Birinci eleman [1 2 3]
Birinci elemanın birinci elemanı 1
```

### ARRAY ÖZELLİKLERİ



*   Arrayin boyutlarını öğrenmek için yapmamız gereken `**array.shape**` yazmak olacaktır

* Arrayin rankini öğrenmek için yapmamız gereken` **array.ndim**` yazmak olacaktır

* Arrayin boyutunu öğrenmek için yapmamız gereken `**array.size**` yazmak olacaktır

* Arrayin barındırdığı veri tipini  öğrenmek için yapmamız gereken `**array.dtype**` yazmak olacaktır



```python
print("Arrayin boyutları {}".format(a.shape))
print("Arrayin ranki {}".format(a.ndim))
print("Arraydeki eleman sayısı {}".format(a.size))
print("Arrayin veri tipi {}".format(a.dtype))
```
```python
Arrayin boyutları (3, 3)
Arrayin ranki 2
Arraydeki eleman sayısı 9
Arrayin veri tipi int64
```

### ARRAY FONKSİYONLARI

**Arraylerden Sıfır Oluşturmak** için yapmamız gereken shape yerine istediğimiz boyutları girmek. Numpy bizim için gerekli arrayi oluşturacaktır. 


```python
shape = (2, 2)
zeros = np.zeros(shape)
zeros
```

```python
array([[0., 0.],
        [0., 0.]])
```


**Arraylerden Bir Oluşturmak** için yapmamız gereken shape yerine istediğimiz boyutları girmek. Numpy bizim için gerekli arrayi oluşturacaktır. 


```python
shape = (2, 2)
ones = np.ones(shape)
ones
```

```python
array([[1., 1.],
    [1., 1.]])
```


**İstediğimiz bir değerle istediğimiz boyutta bir array oluşturmak için ise** `np.full` kullanıyoruz.




```python
a = np.full((6,5), 29) 
print(a)
```

```python
[[29 29 29 29 29]
[29 29 29 29 29]
[29 29 29 29 29]
[29 29 29 29 29]
[29 29 29 29 29]
[29 29 29 29 29]]
```

**İdentity Matrix (Birim Matris)** oluşturmak için ise `np.eye`yazmak olacaktır.


```python
np.eye(3)
```

```python
array([[1., 0., 0.],
        [0., 1., 0.],
        [0., 0., 1.]])
```

Eğer belirli bir aralıkta belirli sayılarla artan bir array oluşturmak istiyorsak `np.arange` kullanmalıyız.


```python
rangearray = np.arange(10,100,10, dtype=float)
rangearray
```

```python
array([10., 20., 30., 40., 50., 60., 70., 80., 90.])
```


Eğer yine belirli bir aralıkta değerler oluşturmak istiyorsak ve kaç tane oluşturacağımızı biliyorsak `np.linspace` kullanabiliriz.


```python
linarray = np.linspace(10, 100, 5) 
linarray
```


```python
array([ 10. ,  32.5,  55. ,  77.5, 100. ])
```


`np.arange` de 100 dahil değildi. Ancak `np.linspace` de dahil. Bunu da gözden kaçırmamak gerekir. 

Şimdi ise çok önemli bir fonksiyon olan `np.reshape` fonksiyonuna bakalım. Bu fonksiyon ile arraylerimizi istediğimiz formata çevirme şansımız var.


```python
array = np.arange(20)
print("Önceki hali : \n", array)
new_array = np.reshape(array, (4,5))
print("Sonraki hali : \n", new_array)
```

```python
Önceki hali : 
    [ 0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19]
Sonraki hali : 
    [[ 0  1  2  3  4]
    [ 5  6  7  8  9]
    [10 11 12 13 14]
    [15 16 17 18 19]]
```



    

### ARRAY INDEKSLEME (ARRAY INDEXING)

Numpy arrayleri indekslemede çok kolaylık sağlıyor. 

#### Bir Boyutlu Array


```python
a1 = np.array([1, 3, 4, 5, 2, 10])
```


```python
a1[0]
```


```python
1
```



```python
a1[4]
```

```python
    2
```

```python
a1[-1]
```

```python
10
```


```python
a1[-3]
```

```python
5
```


##### ÇOK BOYUTLU ARRAY


```python
a2 = np.array([[3, 4, 5, 6],
               [1, 3, 7, 2],
               [8, 4, 5, 10],
               [12, 124, 125, 126]])
```

```python
a2[0]
```


```python
array([3, 4, 5, 6])
```

```python
a2[0, 0]
```


```python
3
```





```python
a2[2, -1]  # 2.indexin son elemanı
```
```python
10
```

```python
a2[2, 0] = 100 # 2.elemanın 0.elemanını 100 yap
a2
```

```python
array([[  3,   4,   5,   6],
        [  1,   3,   7,   2],
        [100,   4,   5,  10],
        [ 12, 124, 125, 126]])
```



```python
a2[[0, 0, 2, 1]] # 0.index, 0.index, 2.index, 1.index
```


```python
array([[  3,   4,   5,   6],
        [  3,   4,   5,   6],
        [100,   4,   5,  10],
        [  1,   3,   7,   2]])
```






```python
a2[:2, ::2] # 2.satıra kadar 0 ile 2. indexler
```


```python
array([[3, 5],
        [1, 7]])
```

```python
a2[::-1, ::-1] # Arrayi ters çevir
```

```python
array([[126, 125, 124,  12],
        [ 10,   5,   4, 100],
        [  2,   7,   3,   1],
        [  6,   5,   4,   3]])
```


```python
a2[:, 0] # İlk sutün
```


```python
array([  3,   1, 100,  12])
```

```python
a2[0, :] # İlk satır
```

```python
array([3, 4, 5, 6])
```

I will ad other features to see how it is going