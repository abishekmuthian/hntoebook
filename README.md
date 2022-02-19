# HN to E-Book

![Shows screenshot of HN stories on Kindle](screenshot/hntoebook_light_3_thumbnail.gif)

### What?
This program converts the best HN stories to .mobi format to be read using an e-reader.

### Why?
[My Hacker News knowledge assimilation stack](https://abishekmuthian.com/my-hacker-news-knowledge-assimilation-stack/).

#### TL;DR
I read HN on my kindle as I find it useful for reading under sun(while walking in a controlled environment) and taking notes. I was running [hntokindle.com](https://web.archive.org/web/20220216140431/https://hntokindle.com/) to offer this as a service for others.

But Amazon has recently made sending bulk emails to the Kindle impossible by requiring 2FA over registered Amazon email address for each item sent to the Kindle.

Hence I've stripped down HN to Kindle code to enable local transfer to any e-book reader which supports .mobi format and made the project open-source.

### How
1. Retrieve HN stories using official API with a [Go wrapper](https://github.com/hoenn/go-hn/).
2. Filter best(Determined by HN) stories older than 9 hours but lesser than 24 hours with at least 20 comments and top comment older than 2 hours. For transferring specific HN item, Story/Comment is received through the item id.
3. Python classifier server is run in the background for classifying story titles against category keywords if category filter option is chosen.
4. Convert the HTML to .pdf after applying cosmetic changes using WKhtmlTopdf with a [Go wrapper](https://github.com/SebastiaanKlippert/go-wkhtmltopdf).
5. Convert the .pdf to .mobi using Calibre command line tool.
6. Place the .mobi file on the device.
7. Store the item id in the K,V database to prevent duplicates.

### Requirements
1. [WKhtmlTopdf](https://wkhtmltopdf.org/downloads.html)
2. [Calibre CLI](https://calibre-ebook.com/download)
3. [hntoebook](https://github.com/abishekmuthian/hntoebook/releases)
#### For Category Filter (Optional)
4. git lfs
5. pytorch
(Other python packages are installed through requirements.txt, Could update your existing packages.)

### Usage

#### Operating System
1. Linux amd64 (Tested)
2. Linux arm64 (Not tested, Reports are welcome)
3. darwin amd64 (Not tested, Reports are welcome)
4. darwin arm64 (Not tested, Reports are welcome)
5. Windows amd64 (Not tested, Reports are welcome)

#### Set the path to store the .mobi file on the e-book reader
./hntoebook -c 

#### Run hntoebook
./hntoebook

#### Send particular HN story or HN comment to the e-book reader
./hntoebook -i

#### Filter HN story categories(Whitelist)
./hntoebook -f

### Feature parity with HN To Kindle
#### Email
Local file transfer is used instead of Email.

#### Send HN item to Kindle
Individuals HN item (Story or Comment) can be sent to the e-book reader.

#### Category Filter
HN stories can be whitelisted by using category keywords and are filtered using a classifier.

### Troubleshooting

#### Errors with mobiPath
Make sure that the path for .mobi files on E-Reader ends with a trailing slash / and the folder where you want to place .mobi files exists prior to running this program.

#### Errors with Category Filter
Category filter requires specific requirements such as Python, PyTorch etc. Models downloaded during config.

#### Not functioning after an error
See if the process executed by the program e.g. uvicorn, calibre, wkhtmltopdf are still running, If so stop the process before executing the program again.

#### Database errors
##### SEGFAULT

Delete the db folder and start again. If you were using < v0.0.3 and upgraded to v0.0.3 then the db folder needs to be deleted regardless of any error as v0.0.3 uses new database.

### License

The MIT License


Copyright 2022 ABISHEK MUTHIAN

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.