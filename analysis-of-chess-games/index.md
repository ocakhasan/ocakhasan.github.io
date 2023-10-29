# Analysis Of My Lichess Bullet Games


Hello Guys,

In this post we will analyze some of my bullet games on [Lichess](https://lichess.org). I love playing chess, and I try to play bullet games almost every day, at least 3-4 games. I thought, it is time to analyze some of my games with `Python`.  To be able to do it, we need to use the [Lichess API](https://lichess.org/api) to export my games. There is already a client package implemented in Python3 called [Berserk](https://berserk.readthedocs.io/en/master/index.html), so I will be using it.

We will use `pandas` to manipulate the data and `matplotlib` to plot some charts (hopefully we will get some meaning based on them).

I have more than 5000 games as of the data 16th Oct 2023. We will be analyzing last 1000 games. 

Let's begin.

I got help from [ChatGPT](https://chat.openai.com/) while creating this notebook.


{{< figure src="/images/lichess.webp" title="Lichess Server Crash Image" >}}

## Fetch the Games


```python
import berserk
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt 
```


```python
with open('./lichess.token') as f:
    token = f.read()
session = berserk.TokenSession(token)
client = berserk.Client(session)
```


```python
client.games.export_by_player('IsaacNewton29', opening=True, perf_type="bullet", max=1000)
```




    <generator object Games.export_by_player at 0x1455dbcf0>




```python
games = list(_)
```


```python
df = pd.DataFrame(games)
df = df.drop(columns=['id', 'rated', 'variant', 'speed', 'perf', 'clock', 'lastMoveAt'])
df.head(2)
```




<div>
<style scoped>
    .dataframe tbody tr th:only-of-type {
        vertical-align: middle;
    }

    .dataframe tbody tr th {
        vertical-align: top;
    }

    .dataframe thead th {
        text-align: right;
    }
</style>
<table border="1" class="dataframe">
  <thead>
    <tr style="text-align: right;">
      <th></th>
      <th>createdAt</th>
      <th>status</th>
      <th>players</th>
      <th>winner</th>
      <th>opening</th>
      <th>moves</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <th>0</th>
      <td>2023-10-16 17:34:35.325000+00:00</td>
      <td>mate</td>
      <td>{'white': {'user': {'name': 'Chessington008', ...</td>
      <td>white</td>
      <td>{'eco': 'B01', 'name': 'Scandinavian Defense: ...</td>
      <td>e4 d5 exd5 Qxd5 Nc3 Qd8 Bc4 Nc6 Nf3 Nf6 d4 e6 ...</td>
    </tr>
    <tr>
      <th>1</th>
      <td>2023-10-16 14:06:40.830000+00:00</td>
      <td>outoftime</td>
      <td>{'white': {'user': {'name': 'IsaacNewton29', '...</td>
      <td>black</td>
      <td>{'eco': 'C00', 'name': 'French Defense', 'ply'...</td>
      <td>e4 e6 Bc4 d5 exd5 exd5 Be2 c6 Nf3 Nf6 d4 Bd6 N...</td>
    </tr>
  </tbody>
</table>
</div>




```python
df.createdAt.min().date(), df.createdAt.max().date()
```




    (datetime.date(2022, 12, 15), datetime.date(2023, 10, 16))



It seems like I played the last 1000 games between the dates 12th December 2022 and 16th October 2023.

Let's extract the white and black players for each game. It will be useful


```python
df["white"] = df["players"].apply(lambda x: x["white"]["user"]["name"])
df["black"] = df["players"].apply(lambda x: x["black"]["user"]["name"])
```


```python
df[df["white"] == "IsaacNewton29"].shape
```




    (490, 8)




```python
df[df["black"] == "IsaacNewton29"].shape
```




    (510, 8)




```python
df[((df["white"] == "IsaacNewton29") & (df["winner"] == "white")) | ((df["black"] == "IsaacNewton29") & (df["winner"] == "black"))].shape
```




    (495, 8)




```python
df[(df["black"] == "IsaacNewton29") & (df["winner"] == "black")].shape
```




    (224, 8)




```python
# Count the occurrences of each category
category_counts = df['winner'].value_counts()

# Plot a bar chart with counts displayed on top of each bar
plt.figure(figsize=(8, 6))
bars = plt.bar(category_counts.index, category_counts.values, color='skyblue')

# Add counts as labels on top of the bars
for bar in bars:
    yval = bar.get_height()
    plt.text(bar.get_x() + bar.get_width()/2, yval, int(yval), ha='center', va='bottom')

plt.xlabel('Player')
plt.ylabel('Count')
plt.title('Winner Counts Based on Color')
plt.xticks(rotation=0)  # Rotate x-axis labels if necessary
plt.show()
```


    
{{< figure src="/images/output_15_0.png"  >}}
    


Let's create a new column to check if I won that game. To get the value

1. If the winner is white check if I play white
2. If the winner is black check if I play black


```python
df["did_i_win"] = ((df["white"] == "IsaacNewton29") & (df["winner"] == "white")) | ((df["black"] == "IsaacNewton29") & (df["winner"] == "black"))
df["did_i_win"].value_counts()
```




    did_i_win
    False    505
    True     495
    Name: count, dtype: int64



Let's see on the chart to understand it better.


```python
# Count the occurrences of each category
category_counts = df['did_i_win'].value_counts()

# Plot a bar chart with counts displayed on top of each bar
plt.figure(figsize=(8, 6))
bars = plt.bar(category_counts.index.values, category_counts.values, color='skyblue')

# Add counts as labels on top of the bars
for bar in bars:
    yval = bar.get_height()
    plt.text(bar.get_x() + bar.get_width()/2, yval, int(yval), ha='center', va='bottom')

plt.xlabel('Is the Game a Win')
plt.ylabel('Count')
plt.title('Did I Win the Game')
plt.xticks(range(len(category_counts.index)), category_counts.index)  # Set x-tick labels
plt.show()
```


{{< figure src="/images/output_19_0.png"  >}}
    


It seems like I lost 505 games whereas I won the 495 games, it is almost 50/50. No improvement at all :smile:

It will be continued...

