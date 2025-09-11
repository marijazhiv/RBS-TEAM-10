# Analiza OS-a, kernela, paketa i konfiguracije logovanja

Ovaj zadatak ima za cilj izvršavanje bezbednosne i konfiguracione
analiye Linux sistema koristeći automatizovanu python skriptu i
ispisujući rezultate u log fajl. Cilj je identifikovati bezbednosne
propuste, potencijalne ranjivosti i oblasti koje zahtevaju unapređenje
kako bi sistem bio sigurniji i u skladu sa standardima dobre prakse za
serverska okruženja.

Na osnovu izvršene analize, izdvojeni su ključni nalazi i preporuke za
unapređenje bezbednosti, optmizaciju konfiguracije i održavanje
minimalnog softverskog profila.

## Informacije o operativnom sistemu

Verzija operativnog sistema je Ubuntu 22.04.5 LTS, a verzija kernela je
6.8.0-79-generic. Ove verzije su moderne i podržane, što smanjuje rizik
od ranjivosti koje su bile prisutne u starijim verzijama operativnog
sistema i kernela. Korišćenje LTS verzija operativnog sistema u
serverskim okruženjima ima više prednosti, a neke od njih su: dugoročna
podrška, stabilnost, kompatibilnost, bezbednost i predvidivo održavanje.
Ipak, i dalje je potrebno redovno ažurirati OS i kernel koristeći *apt
update && apt-upgrade*. Takođe, bilo bi korisno postaviti automatizovane
bezbednosne nadogradnje, instaliranjem paketa *unattended-upgrades*.

Pregledom uptime-a sistema vidimo da je sistem bio aktivan 4 sata, što
znači da nije bio dugo aktivan bez resetovanja, što ukazuje na redovno
praćenje i ažuriranje.

Zaključujemo da je osnovna infrastruktura stabilna, ali dugoročno
praćenje i planirane nadogradnje su i dalje od velike važnosti za
bezbednost i pouzdanost sistema.

## Vremenska konfiguracija

Uvidom u trenutno stanje sistema, zaključujemo da je trenutno vremenska
zona postavljena na Europe/Belgrade, pto je lokalna vremenska zona i to
može dovesti do problema prilikom analize logova i sinhronizacije u
distribuiranim sistemima. Takođe, primećujemo da nije pronađen NTP
daemon (npr. Ntpd, chrony). Zatim primećujemo da je sinhronizacija
vremena trenutno neaktivna, što može otežati praćenje događaja i
autentifikaciju.

Kada se radi o serverskim okruženjima, preporučljivo je promeniti
vremensku zonu na UTC (koja je standardna za servere) komandom
*timedatectl set-timezone UTC* kako bismo izbegli greške pri različitim
vremenskim zonama. Takođe, neophodno je da instaliramo i konfigurišemo
NTP servis (ntpd ili chrony) komandom *apt install chrony*, a zatim
omogućiti servis komandom *systemctl enable \-- now chrony* kako bi
sistem bio sinhronizovan sa vremenskim izvorom.

## Instalirani paketi

Pregledom broja instaliranih paketa, zaključujemo da imamo veći broj
instaliranih paketa, tačnije, imamo ukupno 1848 paketa. Među njima se
nalaze i paketi koji nisu potrebni za serversko okruženje, uključujući
GUI komponente, office pakete i multimedijalni softver.

Bilo bi preporučljivo da uklonimo sve nepotrebne pakete koji povećavaju
površinu napada i troše resurse. Kako bi server bio što lakši za
održavanje i sigurniji, preporučljivo je održavati minimalnu
instalaciju.

Kako bismo to uradili, potrebno je prvo pregledati listu instaliranih
paketa upotrebom komande *apt list --installed*, a zatim ukloniti
nepotrebne pakete komandom *apt remove* ili *apt purge.* Preporučljivo
je instalirati samo neophodan minimum softvera radi manjeg vektora
napada i boljih performansi.

## Konfiguracija logovanja

Pregledom konfiguracije logovanja zaključeno je da je korišćen
systemd-journal, a log fajlovi se čuvaju u /var/log. Log rotation je
aktivan, ali nije postavljeno centralizovano ili udaljeno logovanje.

U narednim koracima bismo trebali da implementiramo remote logging (npr.
Rsyslog ili syslog-ng) kako bi svi logovi bili poslati na siguran
centralni server, čime bi se olakšao nadzor i analiza događaja.
Preporučljivo bi bilo i dodati enkripciju log prenosa (TLS) i kontrolu
pristupa log serveru.

## Sudo pristup i korisnici

Pregledom podešavanja sudo pristupa, utvrđeno je da pun sudo pristup
imaju korisnici root i Jelena. Takođe, utvrđeno je da nema drugih
korisnika sa UID 0 osim root-a, niti korisnika bez lozinke.

Preporučljivo je da ograničimo sudo pristup samo na one korisnike kojima
je zaista potreban i da razmotrimo korišćenje grupa (npr. admin) umesto
korišćenja direktne dodele sudo prava, radi bolje kontrole i evidencije.
Takođe, bilo bi dobro implementirati sudo audit logove radi bolje
evidencije. Dodatno, bilo bi poželjno razmotriti MFA (*multi-factor
authentication*) za privilegovane naloge.

## SSH konfiguracija

U trenutnim uslovima, konfiguracija ne pokazuje nebezbedne postavke, ali
postoje preporuke za dodatno jačanje sigurnosti.

Preporučljivo onemogućiti direktan root login podešavanjem
*PermitRootLogin no* u /etc/ssh/sshd_config, kao i onemogućiti login
preko lozinke i koristiti SSH ključeve, što se postavlja podešavanjem
*PasswordAuthentication no* u /etc/ssh/sshd_config. Opcionalno, možemo
koristiti port knocking ili promeniti default port (sa 22 na drugi) radi
smanjenja brute-force pokušaja.

## Zaključak

Sistem je u osnovi modern i stabilan, ali postoji nekoliko kritičnih
oblasti koje zahtevaju pažnju kako bi ovaj sistem mogao da se koristi
kao bezbedno okruženje za deployment:

- Ukloniti nepotrebne pakete radi smanjivanja površine napada.

- Prebaciti vremensku zonu na UTC i podesiti NTP servis za
  sinhronizaciju.

- Uvesti udaljeno logovanje radi centralne kontrole i audita.

- Ograničiti sudo privilegije i uvesti dodatne bezbednosne mehanizme
  poput MFA.

- Očvrsnuti SSH zabranom root i password autentifikacije.

Sprovođenjem ovih mera sistem postaje značajno sigurniji, pouzdaniji i
usklađeniji sa najboljim praksama za serverska okruženja. Kontinuirano
praćenje i redovno ažuriranje ostaju ključni za dugoročnu bezbednost.
