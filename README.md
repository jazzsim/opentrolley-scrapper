REQUEST
---
POST REQUEST to http://localhost:8080/ with JSON body

**start_id** - determine the starting book's id \
**file_name** - the file name of author json file

```
{
  "start_id": 20,
  "genre": "Comics & Graphic Novels",
  "url": "https://opentrolley.com.my/Search.aspx?category=comics-graphic-novels&page=1&pgsz=3&sorttype=2",
  "file_name": "cgn-authors.json"
}
```

Change ```pgsz``` query to determine the amount of books returned

 ---
RESPONSE
---
**API response**
```
[
  {
    "id": 20,
    "imageUrl": "https://otimages.com/Bookcover/6581/9780062976581.jpg",
    "title": "The Boy, the Mole, the Fox and the Horse: A Great Gift for Book Lovers",
    "genre": "Comics & Graphic Novels",
    "bindingDescription": "Hardcover",
    "language": "English",
    "description": "From a revered British illustrator comes a modern fable for all ages that explores life's universal lessons, featuring 100 color and black-and-white drawings....A Great Gift for Book LoversCharlie Mackesy's beloved The Boy, the Mole, the Fox and the Horse has been adapted into an Academy Award(R) winning animated short film, now available to stream on Apple TV+#1 NEW YORK TIMES BESTSELLER - WALL STREET JOURNAL BESTSELLER  - USA TODAY BESTSELLER\"The Boy, the Mole, the Fox and the Horse is not only a thought-provoking, discussion-worthy story, the book itself is an object of art.\"- Elizabeth Egan, The New York TimesFrom British illustrator, artist, and author Charlie Mackesy comes a journey for all ages that explores life's universal lessons, featuring 100 color and black-and-white drawings.\"What do you want to be when you grow up?\" asked the mole.\"Kind,\" said the boy.Charlie Mackesy offers inspiration and hope in uncertain times in this beautiful book, following the tale of a curious boy, a greedy mole, a wary fox and a wise horse who find themselves together in sometimes difficult terrain, sharing their greatest fears and biggest discoveries about vulnerability, kindness, hope, friendship and love. The shared adventures and important conversations between the four friends are full of life lessons that have connected with readers of all ages.",
    "price": 125.74,
    "discountPrice": 113.05,
    "isbn": "0062976583",
    "publisher": "HARPER ONE",
    "publicationDate": "2023-08-01",
    "pages": 128,
    "postedBy": 34
  },
  {
    "id": 21,
  ...
]
```
---
**Author JSON**

```
[
    {
        "bookId": 20,
        "authorName": "Mackesy Charlie"
    },
    {
        "bookId": 21,
        "authorName": "Esquivel Eric M."
    },
    {
        "bookId": 21,
        "authorName": "Salas Ramon"
    },
    {
        "bookId": 21,
        "authorName": "Davis Darren G."
    },
    {
        "bookId": 22,
        "authorName": "Allo Deborah"
    },
    {
        "bookId": 22,
        "authorName": "Conner Dan"
    },
    {
        "bookId": 22,
        "authorName": "Scalia Roberto"
    }
]
```
