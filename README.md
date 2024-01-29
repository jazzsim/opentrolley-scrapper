REQUEST
---
POST REQUEST to http://localhost:8080/ with JSON body
```
{
  "url": "https://opentrolley.com.my/Search.aspx?category=young-adult-fiction&page=1&pgsz=10&sorttype=2"
}
```

Change ```pgsz``` query to determine the amount of books returned

 ---
RESPONSE
---
```
[
  {
    "image_url": "https://otimages.com/Bookcover/6573/9781339016573.jpg",
    "title": "The Ballad of Songbirds and Snakes (a Hunger Games Novel)",
    "authors": [
      "Collins Suzanne"
    ],
    "binding_description": "Paperback",
    "language": "English",
    "description": "Ambition will fuel him.Competition will drive him.But power has its price.It is the morning of the reaping that will kick off the tenth annual Hunger Games. In the Capitol, eighteen-year-old Coriolanus Snow is preparing for his one shot at glory as a mentor in the Games. The once-mighty house of Snow has fallen on hard times, its fate hanging on the slender chance that Coriolanus will be able to outcharm, outwit, and outmaneuver his fellow students to mentor the winning tribute.The odds are against him. He's been given the humiliating assignment of mentoring the female tribute from District 12, the lowest of the low. Their fates are now completely intertwined -- every choice Coriolanus makes could lead to favor or failure, triumph or ruin. Inside the arena, it will be a fight to the death. Outside the arena, Coriolanus starts to feel for his doomed tribute... and must weigh his need to follow the rules against his desire to survive no matter what it takes.",
    "price": 100,
    "discount_price": 89.91,
    "isbn": "1339016575",
    "publisher": "SCHOLASTIC",
    "publication_date": "01/08/2023",
    "pages": 528
  },
]
```
