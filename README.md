#goNews
# Introduction
News scrapes news articles from Hindustan times, which is an Indian english newspaper.

I am trying to increase my typing speed. Most of the online typing platform have predefined set's of text. After practicing for months, I felt bored to read the same text again and again. I was too lazy to copy text from websites. So I create goNews, which lets you read and copy article. Now I can copy article with just one click. And the I also avoid seeing ads.

**Currently goNews only works on linux systems.


# Getting Started
## Prerequistes
You need to have the following packages intalled
- xclip
- go

## Installation
```bash
git clone git@github.com:yumyum-pi/go-news.git
cd go-news
go install

# check installed application
go-news
```

[![Product Name Screen Shot][product-screenshot]]( https://raw.githubusercontent.com/yumyum-pi/go-news/master/images/go-news-01.jpg ) 
## Usage
```bash
go-news
```

go-news uses vi keybinding for movement along side the arrow keys. Use the 'c' to copy article. ( currently only work for linux with clip )

<!-- CONTRIBUTING -->
## Contributing
Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License
[MIT](https://choosealicense.com/licenses/mit/)


<!-- MARKDOWN LINKS & IMAGES -->
[product-screenshot]: images/go-news.jpg
