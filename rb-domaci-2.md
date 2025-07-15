# Tipovi korisnika i sistema koji komuniciraju sa MegaTravel sistemom

Na osnovu opisa sistema i u skladu sa OWASP principima bezbednosnog dizajna, izdvajamo su sledeće korisnike i eksterni akteri:

## Ljudski korisnici

* Krajnji korisnici (putnici) – pristupaju sistemu preko web/mobilne aplikacije da bi pretraživali, rezervisali i planirali odmor.

* Zaposleni u MegaTravel-u – koriste administrativne interfejse i alate za rad sa rezervacijama, korisnickim nalozima i konfiguracijom sistema.

* Partneri (hoteli, aviokompanije, agencije) – imaju interfejse za sinhronizaciju ponuda i dostupnosti kapaciteta.

## Sistemi i servisi

* Platni procesori – omogućavaju procesiranje online plaćanja.

* OAuth i SSO servisi – koriste se za autentifikaciju korisnika preko trećih strana (Google, Facebook itd).

* Marketing i analitika alati – prate ponašanje korisnika u cilju optimizacije usluga i personalizacije sadržaja.

* CI/CD i DevOps alati – koriste se za isporuku softverskih komponenti i konfiguraciju sistema.

# Površina napada – Identifikacija ulaznih tačaka

Na osnovu interakcije gore navedenih korisnika i sistema sa MegaTravel infrastrukturom, pronašli smo sledeće ulazne tačke koje čine površinu napada:

## 1. Web aplikacija i REST API

 ### Opis: Omogućava korisnicima pretragu, rezervaciju i upravljanje nalozima.

 ### Ulazne tačke:

* HTTPS zahtevi

* Forme za unos podataka

* API endpointi za rezervacije, plaćanja, pretrage

 ### Rizici:

* SQL Injection

* XSS, CSRF

* Neautorizovani pristup API funkcionalnostima

### Zaključak: Kao najčesće korisćena komponenta, predstavlja najveću površinu napada. Potrebno je implementirati strong input validation, rate limiting i authentication middleware.

## 2. Admin portali

### Opis: Koriste ih zaposleni i partneri za upravljanje podacima i rezervacijama.

### Ulazne tačke:

* Web login

* Backend dashboard funkcionalnosti

### Rizici:

* Privilege escalation

* Brute-force napadi

* Konfiguracione greške

### Zaključak: Potrebno je uvesti dvofaktorsku autentifikaciju i audit logove za sve administrativne aktivnosti.

## 3. Webhook endpointi za plaćanja

### Opis: Platni procesori šalju notifikacije o statusima transakcija.

### Ulazne tačke:

* Webhook endpointi

### Rizici:

* Lažni zahtevi (spoofing)

* Replay attacks

### Zaključak: Implementacija verifikacije potpisa i provera IP adresa procesora je ključna.

## 4. Partner API integracije

### Opis: Uvoze podatke o dostupnosti kapaciteta u realnom vremenu.

### Ulazne tačke:

* API pozivi

### Rizici:

* Nevalidni XML/JSON payloadi

* Poverenje u neproverene partnere

### Zaključak: Obezbediti sandbox okruženje i validaciju ulaza pre dalje obrade.

## 5. OAuth SSO endpointi

### Opis: Korisnici se autentifikuju preko trećih strana.

### Ulazne tačke:

* OAuth redirect/callback

### Rizici:

* Token substitution

* Open redirect napadi

### Zaključak: Obavezna validacija state parametra i ograničavanje dozvoljenih redirect URI-jeva.

## 6. Email tokeni i linkovi

### Opis: Linkovi za potvrdu rezervacija i resetovanje lozinki.

### Ulazne tačke:

* GET zahtevi sa tokenima

### Rizici:

* Phishing

* Token leakage

### Zaključak: Tokeni moraju biti vremenski ograničeni i jednokratni.

## 7. Agent desktop aplikacije (call centar)

### Opis: Interfejsi za zaposlene operatere.

### Ulazne tačke:

* Interni dashboardi

### Rizici:

* Insider threat

* Greške operatera

### Zaključak: Uvesti granularnu kontrolu pristupa (least privilege) i nadzor.

## 8. CI/CD i repozitorijumi

### Opis: Automatizovano postavljanje aplikacija i konfiguracija.

### Ulazne tačke:

* Git webhookovi, build serveri

### Rizici:

* Secrets leakage

* Neautorizovani push koda

### Zaključak: Svi tajni podaci moraju biti van koda (npr. vault), uz obaveznu proveru PR-ova i CI logova.

# Zaključak

* Web aplikacija i REST API su najveća površina napada i zahtevaju najveću paznju.

* Interfejsi za integraciju sa trećim stranama (plaćanja, partneri) otvaraju mogućnosti za sofisticirane napade i zahtevaju digitalne potpise, autentifikaciju i sandbox validaciju.

* Administrativni portali i agent alati zahtevaju strogu kontrolu pristupa, auditovanje i obaveznu autentifikaciju u više koraka.

* DevOps alati i repozitorijumi predstavljaju supply chain attack rizik i moraju biti adekvatno izolovani i nadgledani.

* Pokrivanje ovih povrsina kroz principe sigurnog dizajna, primenu OWASP praksi i automatizovane provere predstavlja minimalni prag bezbednosti za MegaTravel sistem.

# Korisćeni izvori:

OWASP Top Ten: https://owasp.org/www-project-top-ten/

OWASP API Security Top 10: https://owasp.org/www-project-api-security/

OAuth 2.0 Security BCP: https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics

OWASP Cheat Sheet Series: https://cheatsheetseries.owasp.org/

Stripe Webhook Security: https://stripe.com/docs/webhooks

NIST SP 800-61 Incident Handling Guide: https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-61r2.pdf

CI/CD Security Risks: https://www.techtarget.com/searchitoperations/tip/9-ways-to-infuse-security-in-your-CI-CD-pipeline