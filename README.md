# idm

idm (it doesn't matter) is a toy implementation of an APL interpreter.

I want to build a toy implementation on an APL-like interpreter. I got this idea after watching [this](https://www.youtube.com/watch?v=PXoG0WX0r_E) video from Rob Pike.

**can do:**

    ./idm
        1
    1
        1 + 1
    2
        a = 2
    2
        a
    2
        a + 1
    3
        b = 1
    1
        a + b
    3

**todo:**

    ./idm
        1 + 2 + 3
    6
        1 + 2 * 3
    9
        2 * 3 + 9
    24
    ...
    ...

Ressources
=====
* [Implementing a bignum calculator](https://www.youtube.com/watch?v=PXoG0WX0r_E)
* [Implementing a bignum calculator - slides](http://go-talks.appspot.com/github.com/robpike/ivy/talks/ivy.slide#1)
* [apl](http://en.wikipedia.org/wiki/APL_%28programming_language%29)
* [synthax and symbols](http://en.wikipedia.org/wiki/APL_syntax_and_symbols)
* [ivy](http://godoc.org/robpike.io/ivy)
* [Handwritten Parsers & Lexers in Go](http://blog.gopheracademy.com/advent-2014/parsers-lexers/)
* [sql-parser](https://github.com/benbjohnson/sql-parser)
* [APL demonstration 1975](https://www.youtube.com/watch?v=_DTpQ4Kk2wA&list=WL&index=13)
* [try apl](http://tryapl.org/)
* [mastering dyalog APL](http://www.dyalog.com/uploads/documents/MasteringDyalogAPL.pdf)
