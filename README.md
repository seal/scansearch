# Public
Customer decided paying was optional, this project is *two* years old (2022) and has some questionable code in it.

Customise .env-example then run docker.

Apparently I didn't know how to use chmod, hence permissions issues on the certificates.....

##  ScanSearch 

![image](https://user-images.githubusercontent.com/25641834/201541535-9e328b36-18a3-4f02-8b9e-2c2b02318b75.png)


> Backend of ScanSearch created 11/13/2022
--- 
### Stuff to add:

- Send out emails for price drops
- Work on api URLs // Future versions

---

### Table of Contents

- [Description](#description)
- [How To Use](#how-to-use)
- [License](#license)

---

## Description

Backend api that essentially acts as a reverse proxy to SerpAPI and mushes with Javascript front-end ( To be coded)

#### Features

- Price comparison for clothing
- Mailchimp notifications on price drops
- TBS customisation on second request 



---

## How To Use
docker build . -t {dockerusername}/scansearch 

docker run -d -v ~/cert:/etc/letsencrypt/live/{yourdomain}/ -p 443:443 -p 80:80 {dockerusername/repo}

Use letsencrypt to get https signatures, then move to ~/cert ( Will remove auto-renew feature, but sorts out docker not having permission to read certificates)
Alternatively you could add docker to the permission group that can read /etc/letsencrypt etc. But I have not added this yet

#### Installation

Add .env file, add following lines:
serpapi={apikey}

#### API Reference

```go
    // Add options here
```


---

## License

Do not care


---


