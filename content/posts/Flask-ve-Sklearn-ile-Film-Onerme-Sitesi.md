---
layout: post
title: Flask ve Sklearn ile Film Önerme Sitesi Yapalım 
summary: Metin benzerliği benzerliğini kullanarak Flask ile film önerme sitesi yapalım.
date: 2021-03-01
tags: [flask, makine ogrenmesi]
---


Bu yazıdaki bütün kodlar [Bu repodan](https://github.com/ocakhasan/movie-recommender) bulunmaktadır. Eğer demo versiyonunu görmek isterseniz [http://banafilmoner.herokuapp.com/](http://banafilmoner.herokuapp.com/) sitesinden görebilirsiniz.

## Gereksinimler
Bu yazımızda yapacağımız siteyi eğer kendiniz de yapmak istiyorsanız [Flask](https://flask.palletsprojects.com/en/1.1.x/) ve [Scikit-learn](https://scikit-learn.org/) kütüphanelerini yüklemeniz gerekmektedir. Bunları yüklemek için terminalden şu komutları yazabilirsiniz ya da her bir paketin dökümentasyonundan bakabilirsiniz.

```bash
pip install Flask 
pip install scikit-learn
```

## Sitenin Yapısı
Yapacağımız sitede film önerileri metin benzerliği ile olacak. Bu filmlerin açıklama metinlerini ise bir veri kümesinden alacağız. Bu veri kümesine [TMDB 5000 Movies](https://www.kaggle.com/tmdb/tmdb-movie-metadata) sayfasından ulaşabilirsiniz. Bundan dolayı önerebileceğimiz metinler sadece bu veri kümesindekiler olacaktır. Metin benzerliğini ise [kosinüs benzerliği](https://merveenoyan.medium.com/kosin%C3%BCs-benzerli%C4%9Fi-2b4a4c924f27) ile yapacağız. 

## Veri Seti ve Metin Benzerliği
Veri setindeki `title` sütunu filmin başlığını ve `overview` sütunu ise filmi basitçe açıklar.Bu yazıda `overview` sütununu kullanarak metin benzerliğini kuracağız. Bunun için önce `utils.py` diye bir dosya oluşturalım ve indirdiğimiz veri setini de projedeki dosyaya koyalım. Öncelikle filmlerin açıklamalarını kullanarak kosinüs benzerliğini verecek olan bir fonksiyon yazalım.

```python
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import linear_kernel

def get_cosine_similarities(df):

    vectorizer = TfidfVectorizer(stop_words="english")

    tf_idf_mat = vectorizer.fit_transform(df['overview'])

    cosine_sim = linear_kernel(tf_idf_mat, tf_idf_mat)

    return cosine_sim
```

`get_cosine_similarities(df)` fonksiyonu parametere olarak `DataFrame` alır, `DataFrame`i ise veri setini okuduktan sonra  bu fonksiyona parametre olarak vereceğiz. Fonksiyonda kullanılan [`TfidfVectorizer`](https://scikit-learn.org/stable/modules/generated/sklearn.feature_extraction.text.TfidfVectorizer.html) metinlerden bilgi çıkarmamıza yarayan bir algoritmadır. Açılımı **Term frequency (tf) -> (terim sıklığı)** ve **inverse document frequency (ters döküman sıklığı)dır**. Yani terimlerin her bir metinde ne kadar sıklıkla geçtiğine ve bu terimlerin bütün dökümanda ne kadar sıklıkla geçtiğine bakıp, hangi terimlerin cümleleri ayırmada önemli olduğuna karar verir. Bu bize (4803, n) boyutunda bir matrix dönderecektir. **n** ise bu algoritmanın bulduğu belirleyici kelimelerin sayısıdır. Yazdığımız fonksiyonla beraber, her bir cümle için her bir kelimenin ne kadar önemi olduğunu gösteren bir matrix elde edilecek. Daha sonra bu matrixi kullanarak her bir metin arasındaki benzerliği bulmak için [linear_kernel](https://scikit-learn.org/stable/modules/generated/sklearn.metrics.pairwise.linear_kernel.html) kullanıyoruz. Bu algoritma ise bize (4803, 4803) boyutunda her bir metnin diğer 4038 filmin metni ile benzerliğini gösteren bir matrix döndürecek. 
Bu fonksiyondan çıkan sonuç ise şu şekildedir

```
[[1.         0.         0.         ... 0.         0.         0.        ]
 [0.         1.         0.         ... 0.02160533 0.         0.        ]
 [0.         0.         1.         ... 0.01488159 0.         0.        ]
 ...
 [0.         0.02160533 0.01488159 ... 1.         0.01609091 0.00701914]
 [0.         0.         0.         ... 0.01609091 1.         0.01171696]
 [0.         0.         0.         ... 0.00701914 0.01171696 1.        ]]
```
Görüldüğü gibi bazı değerler 0 bazıları 1 (köşegendekiler), bazıları da 0 ile 1 arasında. 
Kısaca
1. Sonucu 0 olanlar arasında hiçbir benzerlik yok,
2. 1 olanlar zaten kendileri ile ölçüldüğü için aynı olarak çıkıyor, örnek olarak 1.film ile 1.film arasındaki benzerlik 1 olacak doğal olarak
3. 0-1 arasındakiler ise iki film arasındaki benzerliği gösteriyor. 

Ne yaptığımızı kısaca yazalım. 
* Veri setini okuduk
* Kosinüs benzerlik matriksini oluşturduk.

Şimdi yapılması gerekenler ise bu matrixi kullanıp film önerileri alabilmek. Bunun için yapılması gerekenler
1. Kosinüs matrixini kullanıp bize verilen film için önerileri döndüren bir fonksiyon yazmak
2. Flask ile web arayüzü oluşturup, kullanıcın girdiği filme öneriler vermek
3. Bu fonksiyonu flask ile kullanabilmek.

### Film Önerme Fonksiyonu
Bu fonksiyona geçmeden önce veriyi okuyalım, ve kosinüs matriximizi alalım. Şunu belirtmem gerekir ki, kullanıcının attığı her requestte veri setini baştan okuyup kosinüs matrixini okumak yük olur. Bundan dolayı, bunu bir kez yapmak adına,  bu işlemleri `if __name__ == "__main__"` altında yapacağız.

Öncelikle bir `app.py` adında bir dosya açalım. Bu dosyada Flask applikasyonumuzun kodları olacak. Diğer `utils.py` dan fonksiyonları çağıracağız. 

`app.py` dosyasına şu kodları girelim. 

```python
from flask import Flask, render_template, request, redirect
import pandas as pd
import utils

app = Flask(__name__)

if __name__== "__main__":
    df = pd.read_csv("data.csv")
    df['overview'] = df['overview'].fillna('')
    df['lower_name'] = df['title'].str.lower()

    titles = pd.Series(df.index, index=df['lower_name']).drop_duplicates()

    cosine_sim = utils.get_cosine_similarities(df)

    app.run()
```

Şuan `app.py` dosyasında yapılan işlemler.
1. Flask uygulaması oluşturuldu.
2. Veri okundu.
3. Kosinüs benzerlik matriksi oluşturuldu.

Main kısmında ***titles*** diye bir değişken oluşturulma sebebi bu değişkenin filmleri önerecek olan fonksiyonda kullanacağımızdan dolayıdır. ***Titles*** değişkeni tip olarak `Series`dir. Konsola yazdırdığımız zaman şöyle bir sonuç çıkacaktır. 

```
lower_name
avatar                                         0
pirates of the caribbean: at world's end       1
spectre                                        2
the dark knight rises                          3
john carter                                    4
                                            ...
el mariachi                                 4798
newlyweds                                   4799
signed, sealed, delivered                   4800
shanghai calling                            4801
my date with drew                           4802
Length: 4803, dtype: int64
```

Şimdi filmleri önerecek fonksiyonu yazmaya başlayabiliriz. Bunu `utils.py` dosyasında yazalım. 

```python
"""
movie_title = istenilen filmin ismi
cosine_similarity = kosinüs benzerlik matriksi
titles= az önce oluşturduğumuz filmin isimlerine sahip olan `Series`
df = bütün filmleri barındıran dataframe
"""
def get_recommendations(movie_title, cosine_similarity, titles, df):

    index_movie = titles[movie_title]                   #istenilen filmin indexini bul
    name_of_movie = df.iloc[index_movie]['title']       #daha sonra dataframeden filmin adını bul. 
                                                        #istenilen isim küçük harfli olabilir, biz
                                                        #dataframde nasılsa onu almak için yapıyoruz.
    
    similarities = cosine_similarity[index_movie]       #daha sonra girilen filmin kosinüs benzerlik 
                                                        #arrayini al, diğer filmlerle benzerlik arrayi

    similarity_scores = list(enumerate(similarities))   #işlem kolaylığı için her bir benzerliğin indexini
                                                        #alabilmemiz lazım. yani (0, 0.2), (1, 0.4), (2. 0.7) ... gibi.
    similarity_scores = sorted(similarity_scores , key=lambda x: x[1], reverse = True) #bütün benzerlik skorlarını sırala
    similarity_scores = similarity_scores[1:11]         #en benzer 10 filmi al
    similar_indexes = [x[0] for x in similarity_scores] #benzer filmlerin indexlerini al
    
    return df.iloc[similar_indexes], name_of_movie      #benzer filmlerin bilgilerini almak için indexlerini kullan.
```



### HTML Arayüz
Bu fonksiyon da yazıldığına göre şimdi Flask ile bağlayabiliriz. Ama öncelikle bir arayüzümüz olması gerekiyor. Bunun için aynı klasörde `templates` diye bir klasör oluşturun ve içine `index.html` adında bir dosya oluşturun. Bu dosya bizim kullanıcıdan arayüzü almamızı sağlayacak olan `HTML` kodunu içerecek. `HTML` kısmını anlatmayacağım. Basit şekilde `Flask` bildiğinizi varsayıyorum. 

`index.html` dosyasına [buradaki arayüz kodunu](https://github.com/ocakhasan/movie-recommender/blob/master/templates/index.html) yapıştırın. HTML kısmı şuan çok ilgi alanımız değil, eğer arayüz nasıl görünüyor diye merak ediyorsanız, [buradan](https://banafilmoner.herokuapp.com) bakabilirsiniz.
### Flask Endpointleri halletme
Bu kodda dikkatinizi çekmek istediğim bir nokta var. `FORM` bir '/' yoluna *POST* request yapıyor. Flask uygulamasında '/' adresine bir *POST* request yapılacak. Ayrıca websitesinin giriş sayfası da bu adrese *GET* request yapılarak alınacak. Şimdi `app.py` dosyasında bu koşulları sağlayan kodumuzu yazalım. 

```python
from flask import Flask, render_template, request, redirect, flash, url_for
import pandas as pd
import utils

app = Flask(__name__)



@app.route('/', methods=['GET', 'POST'])
def hello():
    length = 0
    movie_name = ""
    context = {                     #Bu dictionary önerilen filmlerin bilgilerini tutuyor.
        'movies': [],               #isimler
        'urls': [],                 #filmlerin sayfaları
        'release_dates': [],        #filmlerin yayınlanma tarihleri
        'runtimes': [],             #filmlerin süreleri
        'overviews': []             #filmleri anlatan metinler
    }

    if request.method == "POST":                    #Kullanıcı bir input girdiyse
        text = request.form['fname'].lower()
        print("request text", text)
        try:
            recommended_df, movie_name = utils.get_recommendations(
                text, cosine_sim, titles, df)                               #girilen inputtan filmleri al
            context['movies'] = recommended_df.title.values
            context['urls'] = recommended_df.homepage.values
            context['release_dates'] = recommended_df.release_date.values
            context['runtimes'] = recommended_df.runtime.values
            context['overviews'] = recommended_df.overview.values

            length = len(context['movies'])
        except:
            return render_template('index.html', error=True)            #filmi bulamadıysak error döndür.

    return render_template('index.html', length=length, context=context, movie_name=movie_name, error=False) 


if __name__ == '__main__':
    df = pd.read_csv("data.csv")
    df['overview'] = df['overview'].fillna('')

    titles = pd.Series(df.index, index=df['lower_name']).drop_duplicates()

    cosine_sim = utils.get_cosine_similarities(df)

    app.run()

```
Render templatede gönderilen `context` değişkeni `HTML` dosyasında parse ediliyor ve bilgiler güzel bir şekilde gösteriliyor. Dediğim gibi basit şekilde *Flask* bildiğiniz düşünüyorum.

Beğendiyseniz paylaşırsanız çok sevinirim. İyi öğrenmeler.