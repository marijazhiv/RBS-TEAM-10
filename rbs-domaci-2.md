# Motivacija napadača

Razmatrajući potencijalne napadače na sistem *MegaTravel*, moguće je izdvojiti nekoliko grupa napadača, njihove nivoe veština, načine pristupa sistemu i krajnje ciljeve:

1. **Zlonamerni hakeri (*cyber* kriminalci)**  
   * **Motivacija:** Finansijska dobit. Žele da ukradu podatke o korisnicima, posebno finansijske informacije (podatke o kreditnim karticama, brojeve pasoša i druge lične podatke).  
   * **Nivo veštine:** Srednji do visoki. Mogu koristiti poznate ranjivosti, ali i napredne metode (*phishing* [\[1\]](#bookmark=id.9j7pihvp1afo)*, SQL injection* [\[2\]](#bookmark=id.hosyyvr87pv6)*, malware* [\[3\]](#bookmark=id.pp50vowtl9wi)*, ransomware* [\[4\]](#bookmark=id.6ucer5yi7okz)).  
   * **Pristup sistemu:** Eksterni. Napadi dolaze spolja, često putem interneta.  
   * **Ciljevi:**  
     * Krađa ličnih i finansijskih podataka.  
     * Prodaja ukradenih podataka na crnom tržištu.  
     * Ucenjivanje (npr. *ransomware* [\[4\]](#bookmark=id.6ucer5yi7okz))

2. **Insajderi (nezadovoljni zaposleni ili bivši zaposleni):**  
   * **Motivacija:** Osveta, nezadovoljstvo, lična korist.  
   * **Nivo veštine:** Nizak do srednji. Poznaju sistem iznutra, te postoji mogućnost zloupotreba internih privilegija.  
   * **Pristup sistemu:** Interni. Imaju direktan pristup sistemima, podacima ili mreži.  
   * **Ciljevi:**  
     * Sabotiranje poslovanja.  
     * Iznošenje poverljivih informacija.  
     * Krađa podataka o klijentima, partnerima ili kompaniji.

3. **Konkurentske kompanije (industrijska špijunaža):**  
   * **Motivacija:** Pribavljanje poslovnih tajni, strategija, ugovora, planova i korisničkih podataka.  
   * **Nivo veštine:** Srednji do visoki (mogu angažovati stručne hakere).  
   * **Pristup sistemu:** Eksterni, često prikriven (putem *phishing*\-a [\[1\]](#bookmark=id.9j7pihvp1afo), lažnih zaposlenih ili društvenog inženjeringa).  
   * **Ciljevi:**   
     * Pristup internim dokumentima i analizama.  
     * Sabotiranje reputacije *MegaTravel*\-a.  
     * Krađa inovativnih ideja i planova.

4. **Skripter kidiji (*script kiddies*** [\[5\]](#bookmark=id.c8qnmfohjgg3)**):**  
   * **Motivacija:** Dosada, dokazivanje, slava u hakerskoj zajednici.  
   * **Nivo veštine:** Nizak. Koriste gotove alate i skripte.  
   * **Pristup sistemu:** Eksterni.  
   * **Ciljevi:**   
     * Prekid rada sistema.  
     * *Defacement* sajta (promena izgleda)  
     * Ostavljanje tragova upada (tzv. *tagging* [\[6\]](#bookmark=id.3ab0es440kvz))

5. **Aktivisti (*Hacktivists*):**  
   * **Motivacija:** Politički ili ideološki motivisana akcija (npr. protiv korporacija, globalizacije, ekološki ciljevi)  
   * **Nivo veštine:** Srednji.  
   * **Pristup sistemu:** Eksterni.  
   * **Ciljevi:**  
     * Javno kompromitovanje kompanije (objavljivanje podataka).  
     * Sabotiranje usluga (npr. DDoS napadi [\[7\]](#bookmark=id.kf1b70yh5v57)).  
     * Skretanje pažnje na politička pitanja

6. **Državni akteri (državno sponzorisani napadači):**  
   * **Motivacija:** Geopolitički interesi, špijunaža, destabilizacija tržišta.  
   * **Nivo veštine:** Visoki. Koriste *zero-day vulnerabilities* [\[8\]](#bookmark=id.5m366b84tyyk) i napredne napade (*APT, Advanced Persistent Threat* [\[9\]](#bookmark=id.wzytwppv6tp3)).  
   * **Pristup sistemu:** Eksterni, ali sofisticiran.  
   * **Ciljevi:**   
     * Pristup informacijama o građanima, diplomatama ili strateškim podacima.  
     * Sabotiranje međunarodnog poslovanja.  
     * Uticaj na stabilnost država ili tržišta u kojima posluje *MegaTravel*.

*Tabela 1\. Tabela pregleda napadača*

| Klasa napadača | Nivo veštine | Pristup | Glavni ciljevi |
| :---: | :---: | :---: | :---: |
| *Cyber* kriminalci | Srednji do visok | Eksterni  | Finansijska dobit, krađa podataka |
| Insajderi  | Nizak do srednji | Interni  | Osveta, sabotaža |
| Konkurencija  | Srednji do visok | Eksterni | Špijunaža, kompromitacija |
| Skripter kidiji (*script kiddies*) | Nizak  | Eksterni  | Slava, zabava |
| Aktivisti (*hacktivists*) | Srednji  | Eksterni | Ideološki protesti |
| Državni akteri | Visok | Eksterni | Špijunaža, geopolitički ciljevi |

## Dodatna objašnjenja

* ***Phishing:*** Prevara kojom se korisnici navode da otkriju poverljive informacije, obično putem lažnih i-mejlova ili sajtova [\[1\]](#bookmark=id.9j7pihvp1afo).  
* ***SQL Injection:*** Napad kojim se zlonamerni SQL kod unosi u bazu podataka kroz ranjive forme [\[2\]](#bookmark=id.hosyyvr87pv6).  
* **Malver (*malware*):** Zlonamerni softver koji može izazvati štetu na računaru korisnika ili ukrasti podatke [\[3\]](#bookmark=id.pp50vowtl9wi).  
* ***Ransomware:*** Tip malvera koji zaključava fajlove korisnika i traži otkup [\[4\]](#bookmark=id.6ucer5yi7okz).  
* ***Script kiddies:*** Napadači bez većeg znanja koji koriste gotove alate [\[5\]](#bookmark=id.c8qnmfohjgg3).  
* ***Tagging (defacement):*** Promena izgleda sajta radi ostavljanja poruke ili potpisa [\[6\]](#bookmark=id.3ab0es440kvz).  
* ***DDoS (Distributed Denial of Service):*** Napad pri kome veliki broj kompromitovanih uređaja (botova) istovremeno šalje zahteve serveru ili mreži sa ciljem da preoptereti sistem i onemogući pristup legitimnim korisnicima [\[7\]](#bookmark=id.kf1b70yh5v57).  
* ***Zero-day vulnerabilities:*** Bezbednosni propusti koji nisu javno prepoznati ili zakrpljen [\[8\]](#bookmark=id.5m366b84tyyk).  
* ***Advanced Persistent Threats (APT):*** Dugotrajne, ciljno usmerene pretnje koje sprovode organizovani ili državno podržani napadači [\[9\]](#bookmark=id.wzytwppv6tp3).

# B. Osetljiva imovina MegaTravel-a 

U okviru procesa modelovanja pretnji (Threat Modeling), jedan od ključnih koraka je identifikacija i analiza osetljive imovine (assets). Ovaj korak je od suštinskog značaja za razumevanje potencijalnih rizika i prioritizaciju bezbednosnih mera. Kompanija MegaTravel, kao međunarodni lider u oblasti turističkih usluga, obrađuje veliki broj ličnih i poslovnih podataka korisnika, koristi kompleksnu infrastrukturu i sarađuje sa brojnim eksternim partnerima. Imajući u vidu različite vrste napadača, njihove ciljeve i motivaciju, kao i zakonske okvire kao što su GDPR i PCI DSS, u nastavku sledi detaljna analiza ključne imovine MegaTravel sistema. 

## 1\. Baza podataka korisnika (User Database) 

* **Inherentna izloženost:** Pristup ovoj bazi podataka imaju backend serveri aplikacije, sistem administratori i određeni članovi IT tima sa odgovarajućim privilegijama. Odnosno ovlašćeno IT osoblje.   
* **Bezbednosni ciljevi (CIA):**   
* *Poverljivost:* Baza sadrži lične podatke korisnika (ime, adresa, pasoš, podaci o plaćanju), čija zaštita je obavezna po zakonu.   
* *Integritet:* Neophodno je očuvanje tačnosti i konzistentnosti svih podataka.   
* *Dostupnost:* Sistem mora biti uvek dostupan bez prekida za potrebe rezervacija i korisničke podrške, odnosno omogućavati pristup korisničkim podacima.   
* **Uticaj narušavanja:**   
* Povreda poverljivosti može dovesti do narušavanja privatnosti korisnika, što može izazvati ozbiljne pravne posledice i dovesti do sudskih tužbi, kazni (npr. prema GDPR) kao i gubitka reputacije.   
* Povreda integriteta može uzrokovati greške u rezervacijama i stoga dovesti do pogrešnih rezervacija i nezadovoljstva korisnika i gubitak poverenja.   
* Nedostupnost baze blokira celokupnu funkcionalnost sistema. 

## 2\. Rezervacioni sistem (Reservation System) 

* **Inherentna izloženost:** Rezervacioni sistem je dostupan korisnicima putem interneta, kao i putem API-ja za partnere (interne servise).   
* **Bezbednosni ciljevi (CIA):**   
* *Poverljivost:* Čuva informacije i detalje o putovanjima korisnika i partnera.   
* *Integritet:* Mora osigurati tačnost svih rezervacija (datuma, destinacija, cena i uslova) i termina.   
* *Dostupnost:* Neophodno je neprekidno funkcionisanje kako bi korisnici mogli da pristupe uslugama. Dakle, mora biti dostupan 24/7.   
* **Uticaj narušavanja:**   
* Napad koji menja rezervacije i dovodi do pogrešnih može izazvati logistički haos.   
* Nedostupnost direktno utiče na prihode kompanije i korisničko zadovoljstvo. 

## 3\. Sistem plaćanja (Payment Gateway) 

* **Inherentna izloženost:** Pristupaju mu krajnji korisnici, finansijski provajderi (finansijske institucije) i interni procesni sistemi (integrisani partneri).   
* **Bezbednosni ciljevi (CIA):**   
* *Poverljivost:* Mora štititi podatke o kreditnim  karticama, brojevima računa i drugim platnim metodama.   
* *Integritet:* Svaka transakcija mora biti verifikovana i precizna. Iznosi i transakcije ne smeju biti izmenjeni.   
* *Dostupnost:* Plaćanje mora biti omogućeno u svakom trenutku.   
* **Uticaj narušavanja:**   
* Krađa finansijskih podataka može izazvati pravne posledice i gubitak poverenja odnosno reputacionu štetu.   
* Neuspeh transakcija i nemogućnost plaćanja onemogućava ostvarivanje prihoda (gubitak profita). 

## 4\. Poslovni i logistički podaci 

* **Inherentna izloženost:** Dostupni isključivo višem menadžmentu i odabranim zaposlenima.   
* **Bezbednosni ciljevi (CIA):**   
* *Poverljivost:* Obuhvata strateške dokumente, ugovore sa partnerima, strategije i interni know-how.   
* *Integritet:* Tačnost ovih podataka je presudna za donošenje odluka.   
* *Dostupnost:* Potrebni su u realnom vremenu za poslovne operacije.  
* **Uticaj narušavanja:**   
* Povreda može dovesti do industrijske špijunaže ili curenja podataka i konkurentske prednosti za druge firme. Time biva ugrožena tržišna pozicija kompanije. 

## 5\. E-mail sistem i komunikacija zaposlenih 

* **Inherentna izloženost:** Pristupaju mu svi zaposleni unutar kompanije. Svako svom e-mail nalogu.   
* **Bezbednosni ciljevi (CIA):**   
* *Poverljivost:* Obuhvata internu komunikaciju i razmenu sa klijentima.   
* *Integritet:* Potrebno je osigurati verodostojnost (autentičnost) poruka.   
* *Dostupnost:* Komunikacija mora teći bez prekida.   
* **Uticaj narušavanja:**   
* Mogućnost phishing napada, curenja poverljivih informacija i zlonamerne manipulacije sadržajem.   
* Prekid u komunikaciji remeti poslovne procese. 

## 6\. Web server i korisnička aplikacija (Front-end) 

* **Inherentna izloženost:** Direktno izloženi internetu i dostupni krajnjim korisnicima širom sveta.   
* **Bezbednosni ciljevi (CIA):**   
* *Poverljivost:* Zaštita korisničkih sesija i autentifikacionih tokena.   
* *Integritet:* Web stranice i aplikacija moraju biti autentične i izmenama otporne. Odnosi se na kôd aplikacije i interfejse.   
* *Dostupnost:* Ključno za pružanje usluga 24/7.   
* **Uticaj narušavanja:**   
* Defacement sajta, DDoS napadi ili zlonamerni kod mogu narušiti korisničko iskustvo i poverenje što utiče na ugled i funkcionalnost sistema. 

![dfd1][images/dfd1.png]

### E. Pretnje i mitigacije \- STRIDE metodologija (za svako sredstvo imovine)

---

| Slovo | Pretnja | Objašnjenje |
| ----- | ----- | ----- |
| **S** | **Spoofing** | Lažno predstavljanje korisnika ili sistema radi pristupa resursima. |
| **T** | **Tampering** | Neovlašćena izmena podataka ili koda. |
| **R** | **Repudiation** | Poricanje izvršene akcije bez dokaza (npr. korisnik poriče plaćanje). |
| **I** | **Information Disclosure** | Neovlašeno otkrivanje poverljivih informacija. |
| **D** | **Denial of Service (DoS)** | Sprečavanje legitimnog korišćenja sistema ili resursa. |
| **E** | **Elevation of Privilege** | Neovlašeno sticanje višeg nivoa pristupa nego što je dozvoljeno. |

---

### 1\. Baza korisnika

* **Pretnje (STRIDE):**

  * *Spoofing:* Lažno predstavljanje korisnika radi pristupa podacima.

  * *Information Disclosure:* Curenje ličnih podataka.

  * *Tampering:* Neautorizovana izmena korisničkih podataka.

* **Mitigacije:**

  * **Višefaktorska autentifikacija.**

Umesto da se korisnik identifikuje samo pomoću lozinke (nešto što zna-  lozinka, PIN), višefaktorska autentifikacija dodaje još **jedan ili više faktora** (nešto što imaš \- mobilni telefon, token, smart kartica; nešto što jesi \- otisak prsta, prepoznavanje lica).

Čak i ako napadač sazna korisničku lozinku, neće moći da pristupi nalogu bez drugog faktora (npr. koda sa telefona), čime se **drastično smanjuje rizik od spoofinga**.

* **Enkripcija u mirovanju i prenosu (Encryption at rest and in transit)**

**Enkripcija u mirovanju (at rest):**  
Podaci u bazi su šifrovani čak i kada nisu u upotrebi. Ako neko fizički dođe do fajla baze podataka, ne može da ga pročita bez ključa.  
 ➤ Primer: AES-256 enkripcija baze.

**Enkripcija u prenosu (in transit):**  
Kada podaci putuju između klijenta i servera (npr. login podaci), oni su šifrovani korišćenjem TLS/SSL (HTTPS), pa niko "između" ne može da ih presretne.  
 ➤ Primer: HTTPS sa TLS 1.2/1.3.

* **Kontrola pristupa (RBAC)**

RBAC je mehanizam koji definiše **ko može da radi šta u sistemu** na osnovu svoje uloge (npr. **korisnik** može menjati samo svoje podatke, **admin** može upravljati korisnicima, ali ne vidi njihove lozinke, **baza** se može menjati samo putem verifikovanog API-ja). Ako korisnik **nema prava da izmeni određene podatke**, ni slučajno ni zlonamerno ne može to da uradi. RBAC sprečava greške, sabotaže i neautorizovani pristup.  
Npr: U MegaTravel aplikaciji g**uest** ima samo pravo pregleda sadržaja, **registered user** može praviti rezervacije, plaćati i ostavljati recenzije, **partner (npr. vlasnik smeštaja)** upravlja sopstvenim ponudama i vidi rezervacije koje se njega tiču, **support agent** pomaže korisnicima i rešava probleme, ali nema pristup plaćanjima, **administrator** ima potpuni pristup svim podacima i sistemskim podešavanjima, **system (servisi)** automatski obrađuje rezervacije i šalje notifikacije, uz ograničen tehnički pristup.  
Svaka uloga ima jasno definisane dozvole prema principu "najmanjih privilegija" radi bezbednosti.

### 2\. Rezervacioni sistem

* **Pretnje:**

  * *Tampering:* Neautorizovana promena rezervacija.

  * *Denial of Service (DoS):* Sistem postaje nedostupan.

* **Mitigacije:**

  * **Validacija ulaza (input validation)**

Ako sistem **ne proverava šta korisnik unosi**, onda **zlonamerni korisnik može manipulisati podacima** (npr. u browseru, kroz "Inspect Element", Postman, ili skriptama) i poslati **namerne, nevažeće ili opasne vrednosti**.

- Negativna cena, nepravilan datum: **Business logic attack**.  
- ' OR 1=1 – : **SQL Injection**. Može da izvrši SQL kod koji vrati sve rezervacije, izbriše podatke, ili izmeni tuđe podatke.  
  * **Rate limiting, CAPTCHA zaštita**  
* **Rate limiting:** Ograničava broj zahteva koje korisnik može da pošalje u određenom vremenu (npr. najviše 5 rezervacija u minuti).  
* **CAPTCHA:** Pita korisnika da dokaže da je čovek (npr. klikni na semafore) pre nego što nastavi.  
* Štite od **DoS napada** i **automatskih botova** koji pokušavaju da preopterete server ili da masovno prave rezervacije.

(Ako neko pokuša da pošalje 1000 zahteva za rezervaciju u sekundi \- sistem automatski blokira ili usporava te zahteve).

* **Redudantni serveri (failover)**

Sistem ima više servera spremnih da preuzmu posao ako jedan padne (npr. backup serveri, klasteri, load balancing). Štite od **DoS napada i hardverskih kvarova**, jer omogućavaju da sistem nastavi da radi čak i ako jedan deo infrastrukture zakaže.

### 3\. Sistem za plaćanje

* **Pretnje:**

  * *Information Disclosure:* Krađa finansijskih podataka.

  * *Repudiation:* Korisnik poriče da je izvršio plaćanje.

  * *Tampering:* Promena iznosa u toku transakcije.

* **Mitigacije:**

  * **Tokenizacija i PCI-DSS usklađenost**

Štiti od Information Disclosure (štite korisničke kartice i finansijske informacije od curenja ili krađe).

* **Tokenizacija** zamenjuje osetljive podatke (npr. broj kartice) sa **nasumičnim tokenima** koji **nemaju vrednost van sistema** (`4111-1111-1111-1111` ➝ `tok_9f82hf2u).` Ako napadač ukrade token **ne može ga iskoristiti** bez konteksta sistema.

* **PCI-DSS (Payment Card Industry Data Security Standard)** je skup obaveznih standarda za zaštitu kartičnih podataka (Enkripcija, Ograničen pristup, Redovno testiranje ranjivosti, Mrežna segmentacija).

  * **Digitalni potpisi transakcija**

**Digitalni potpis** je kriptografski mehanizam koji **dokazuje autentičnost i integritet poruke ili transakcije** (privatni/javni ključ).

* **Logovanje i auditing sistema**

Štiti od Repudiation i Tampering. Sistem vodi **detaljne zapise (logove)** o svakoj transakciji, vremenu, IP adresi, korisniku, svakoj izmeni, pokušaju pristupa ili grešci.  
**Audit trail** (trag provere) omogućava da se svaki događaj **rekonstruiše kasnije**, detektuje zloupotreba,  pravno dokaže ko je i kada šta uradio.  
Sprečava korisnike da **poriču radnje**, i omogućava detekciju **zlonamernih izmena** u sistemu.

### 4\. Interna poslovna dokumentacija

* **Pretnje:**

  * *Information Disclosure:* Industrijska špijunaža.

  * *Elevation of Privilege:* Neovlašćeni pristup internim fajlovima.

* **Mitigacije:**

  * **Slojeviti pristup (least privilege)**

Štiti od Elevation of Privilege.

* Pristup dokumentima se zasniva na **potrebi za rad** \- korisnik **dobija pristup samo onome što mu je neophodno**.  
  Na primer **zaposleni u korisničkoj podršci** ne vidi finansijske izveštaje, **partner** ne vidi strategiju razvoja tržišta.  
  Time se sprečava da neko sa manje privilegija „zaluta” u delove sistema koje **ne bi smeo ni da vidi**, a kamoli da menja.


  * **Šifrovana interna komunikacija.**

Štiti od Information Disclosure.

Svi fajlovi i komunikacija (mejlovi, interne poruke, prenos fajlova) **moraju biti šifrovani**, čak i unutar organizacije.

Koriste se standardi kao što su:

* **TLS** za mrežnu komunikaciju,

* **AES-256** za enkripciju fajlova u mirovanju,

* VPN za pristup internoj mreži spolja.

  * **Monitoring pristupa i audit trail**

Štiti od Information Disclosure, Elevation of Privilege.

Sistem beleži:

* Ko je pristupio kom dokumentu,

  * Kada i odakle (IP adresa),

  * Da li je samo čitao, menjao ili brisao dokument.

Svi osetljivi pristupi se **automatski evidentiraju**, a neuobičajeno ponašanje (npr. previše downloada) može pokrenuti **alarm ili istragu**. Omogućava **pravovremenu reakciju** i **dokazivanje** ako dođe do incidenta.

### 5\. E-mail sistem

**Pretnje:**

* *Spoofing:* Lažni pošiljaoci mejlova.

  * *Phishing:* Prevara zaposlenih da odaju informacije.

**Mitigacije:**

* **SPF, DKIM zaštita**

**SPF (Sender Policy Framework)**

* SPF je mehanizam koji omogućava serverima da provere da li je IP adresa sa koje je stigao mejl ovlašćena da šalje mejlove u ime određenog domena.  
  Ako neko pokuša da pošalje lažni mejl (spoofing) sa tvog domena, primaoci mogu da otkriju da mejl nije sa legitimnog servera i da ga odbace ili označe kao sumnjiv.

**DKIM (DomainKeys Identified Mail)**

* DKIM dodaje digitalni potpis u zaglavlje mejla, koji omogućava primaocu da proveri da li mejl stvarno potiče sa domena pošiljaoca i da li nije izmenjen tokom prenosa. Pomaže u sprečavanju lažnog predstavljanja (spoofinga) i osigurava integritet sadržaja mejla.

  * **Antivirus i anti-phishing alati**

**Antivirus softver** skenira dolazne mejlove i priloge u potrazi za malicioznim kodom, virusima, trojancima ili drugim štetnim programima koji mogu ugroziti računar ili mrežu.

**Anti-phishing alati** automatski detektuju sumnjive ili lažne mejlove koji pokušavaju da prevarom (phishing) od korisnika izvuku poverljive informacije poput lozinki, brojeva kartica i sl.

Ovi alati mogu blokirati takve mejlove, upozoriti korisnika ili ih automatski preusmeriti u posebnu fasciklu za sumnjive mejlove.

* **Obuka zaposlenih o bezbednosti**

Čak i sa svim tehničkim merama, ljudski faktor je često najslabija karika. Zaposleni treba da budu edukovani kako da prepoznaju sumnjive mejlove, phishing pokušaje, lažne linkove i priloge.

### Trust Boundaries (Granice poverenja)

* Granica između korisničkog browsera i Web aplikacije (nepoveren korisnik).

* Granica između aplikacionog sloja i baze podataka.

* Granica između Web aplikacije i platnog provajdera (eksterni entitet).

* Granica između interne mreže i eksternih sistema (DMZ).


[image1]: <data:image/png;base64,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>


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

1. [NCSC \- UK: Phishing](https://www.ncsc.gov.uk/guidance/phishing)  
2. [OWASP: SQL Injection](https://owasp.org/www-community/attacks/SQL_Injection)  
3. [IBM: What is malware?](https://www.ibm.com/think/topics/malware)  
4. [CISA: Ransomware guide](https://www.cisa.gov/stopransomware/ransomware-guide)  
5. [Malwarebytes: What is a script kiddie?](https://www.malwarebytes.com/cybersecurity/basics/script-kiddie)  
6. [Wikipedia: Website defacement](https://en.wikipedia.org/wiki/Website_defacement)  
7. [Cloudflare: What is a DDoS attack?](https://www.cloudflare.com/learning/ddos/what-is-a-ddos-attack/)  
8. [Wikipedia: Zero-day vulnerability](https://en.wikipedia.org/wiki/Zero-day_vulnerability)  
9. [Kaspersky: What is APT (Advanced Persistent Threat)?](https://www.kaspersky.com/resource-center/definitions/advanced-persistent-threats)  
10. [Microsoft. “Threat Modeling: STRIDE Approach.”](https://learn.microsoft.com/en-us/security/compass/stride-threat-modeling)  
11. [OWASP Foundation. “OWASP Top Ten Security Risks.”](https://owasp.org/www-project-top-ten/)   
12. [Google reCAPTCHA](%20https://www.google.com/recaptcha/about/)

OWASP Top Ten: https://owasp.org/www-project-top-ten/

OWASP API Security Top 10: https://owasp.org/www-project-api-security/

OAuth 2.0 Security BCP: https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics

OWASP Cheat Sheet Series: https://cheatsheetseries.owasp.org/

Stripe Webhook Security: https://stripe.com/docs/webhooks

NIST SP 800-61 Incident Handling Guide: https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-61r2.pdf

CI/CD Security Risks: https://www.techtarget.com/searchitoperations/tip/9-ways-to-infuse-security-in-your-CI-CD-pipeline
