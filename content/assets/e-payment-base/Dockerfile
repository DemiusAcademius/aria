FROM mcr.microsoft.com/dotnet/core/aspnet:2.2

RUN apt-get update -y \
    && apt-get install -y gnupg \
    && rm -rf /var/lib/apt/lists/*

ADD pubrings.tar.gz ./pubrings

RUN gpg -v --batch --import pubrings/ApaCanal/ApaCanal_pubring.pkr
RUN gpg -v --batch --import pubrings/ApaCanal/ApaCanal_secring.pkr

RUN gpg -v --batch --import pubrings/Bcr/bcr_pub.key
RUN gpg -v --batch --import pubrings/Comertbank/Procopan.asc
RUN gpg -v --batch --import pubrings/Energbank/office.asc
RUN gpg -v --batch --import pubrings/Eximbank/Eximbank.asc
RUN gpg -v --batch --import 'pubrings/Fincombank/Bat Aurel.asc'
RUN gpg -v --batch --import 'pubrings/MobiasBanca/Sergiu Virtan.asc'
RUN gpg -v --batch --import pubrings/Moldindconbank/MICB_Retail_public.asc
RUN gpg -v --batch --import pubrings/MoldovaAgroindbank/all-my-keys
RUN gpg -v --batch --import pubrings/PostaMoldovei/ANDREI.asc
RUN gpg -v --batch --import "pubrings/Victoria Bank/VictoriaBank_pubring.pkr"

RUN gpg -v --batch --import "pubrings/Demius/demius.md.public-key.asc"