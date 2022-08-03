# corder

Recorder for [Colly](https://github.com/gocolly/colly) program.

Record `Start Time`, `Request Count`, `Response Count` and `Error List`.

## Usage

```go
c := colly.NewCollector()
c.OnHTML(...)
...

corder := NewCorder(c)

c.Visit(...)

corder.Print(os.Stdout)
```

output:

```
--------- Colly Corder ---------
          cost: 2.383962292s
 request count: 27
response count: 25
   error count: 2

error detail:
> [ Not Found ] (count: 1)
  - http://www.baidu.com/duty/mianze-shengming.html
> [ Get "http://help.baidu.com/question": Not following redirect to help.baidu.com because its not in AllowedDomains ] (count: 1)
  - http://www.baidu.com/search/jiqiao.html
--------------------------------
```